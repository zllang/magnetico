package metadata

import (
	"log"
	"net"
	"sync"
	"time"

	"github.com/tgragnato/magnetico/dht"
	"github.com/tgragnato/magnetico/persistence"
)

const (
	// PeerIDLength The peer_id is exactly 20 bytes (characters) long.
	// https://wiki.theory.org/BitTorrentSpecification#peer_id
	PeerIDLength = 20
	// PeerPrefix Azureus-style
	PeerPrefix = "-UT3600-"
)

type Metadata struct {
	InfoHash []byte
	// Name should be thought of "Title" of the torrent. For single-file torrents, it is the name
	// of the file, and for multi-file torrents, it is the name of the root directory.
	Name         string
	TotalSize    uint64
	DiscoveredOn int64
	// Files must be populated for both single-file and multi-file torrents!
	Files []persistence.File
}

type Sink struct {
	PeerID      []byte
	deadline    time.Duration
	maxNLeeches int
	drain       chan Metadata

	incomingInfoHashes   map[[20]byte][]net.TCPAddr
	incomingInfoHashesMx sync.RWMutex

	terminated  bool
	termination chan interface{}
}

func NewSink(deadline time.Duration, maxNLeeches int) *Sink {
	ms := new(Sink)

	ms.PeerID = randomID()
	ms.deadline = deadline
	ms.maxNLeeches = maxNLeeches
	ms.drain = make(chan Metadata, 10)
	ms.incomingInfoHashes = make(map[[20]byte][]net.TCPAddr)
	ms.termination = make(chan interface{})

	return ms
}

func (ms *Sink) Sink(res dht.Result) {
	if ms.terminated {
		log.Panicln("Trying to Sink() an already closed Sink!")
	}

	// cap the max # of leeches
	ms.incomingInfoHashesMx.RLock()
	currentLeeches := len(ms.incomingInfoHashes)
	ms.incomingInfoHashesMx.RUnlock()
	if currentLeeches >= ms.maxNLeeches {
		return
	}

	infoHash := res.InfoHash()
	peerAddrs := res.PeerAddrs()

	ms.incomingInfoHashesMx.RLock()
	_, exists := ms.incomingInfoHashes[infoHash]
	ms.incomingInfoHashesMx.RUnlock()
	if exists || len(peerAddrs) <= 0 {
		return
	}

	peer := peerAddrs[0]
	ms.incomingInfoHashesMx.Lock()
	ms.incomingInfoHashes[infoHash] = peerAddrs[1:]
	ms.incomingInfoHashesMx.Unlock()

	go NewLeech(infoHash, &peer, ms.PeerID, LeechEventHandlers{
		OnSuccess: ms.flush,
		OnError:   ms.onLeechError,
	}).Do(time.Now().Add(ms.deadline))
}

func (ms *Sink) Drain() <-chan Metadata {
	if ms.terminated {
		log.Panicln("Trying to Drain() an already closed Sink!")
	}
	return ms.drain
}

func (ms *Sink) Terminate() {
	ms.terminated = true
	close(ms.termination)
	close(ms.drain)
}

func (ms *Sink) flush(result Metadata) {
	if ms.terminated {
		return
	}

	ms.drain <- result

	var infoHash [20]byte
	copy(infoHash[:], result.InfoHash)
	ms.delete(infoHash)
}

func (ms *Sink) onLeechError(infoHash [20]byte, err error) {
	ms.incomingInfoHashesMx.RLock()
	peers, exists := ms.incomingInfoHashes[infoHash]
	ms.incomingInfoHashesMx.RUnlock()
	if !exists || len(peers) == 0 {
		return
	}

	if len(peers) == 1 {
		ms.delete(infoHash)
	} else {
		ms.incomingInfoHashesMx.Lock()
		ms.incomingInfoHashes[infoHash] = peers[1:]
		ms.incomingInfoHashesMx.Unlock()
	}

	go NewLeech(infoHash, &peers[0], ms.PeerID, LeechEventHandlers{
		OnSuccess: ms.flush,
		OnError:   ms.onLeechError,
	}).Do(time.Now().Add(ms.deadline))
}

func (ms *Sink) delete(infoHash [20]byte) {
	ms.incomingInfoHashesMx.Lock()
	defer ms.incomingInfoHashesMx.Unlock()
	delete(ms.incomingInfoHashes, infoHash)
}
