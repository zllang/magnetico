package metadata

import (
	"net"
	"reflect"
	"testing"
	"time"
)

func TestSink_NewSink(t *testing.T) {
	t.Parallel()

	sink := NewSink(time.Second, 10)
	if sink == nil ||
		len(sink.PeerID) != 20 ||
		sink.deadline != time.Second ||
		sink.maxNLeeches != 10 ||
		sink.drain == nil ||
		sink.incomingInfoHashes == nil ||
		sink.termination == nil {
		t.Error("One or more fields of Sink were not initialized correctly")
	}

}

type TestResult struct {
	infoHash  [20]byte
	peerAddrs []net.TCPAddr
}

func (tr *TestResult) InfoHash() [20]byte {
	return tr.infoHash
}

func (tr *TestResult) PeerAddrs() []net.TCPAddr {
	return tr.peerAddrs
}

func TestSink_Sink(t *testing.T) {
	t.Parallel()

	sink := NewSink(time.Minute, 2)
	if len(sink.incomingInfoHashes) != 0 {
		t.Error("incomingInfoHashes field of Sink has not been initialized correctly")
	}
	testResult := &TestResult{
		infoHash:  [20]byte{255},
		peerAddrs: []net.TCPAddr{{IP: net.ParseIP("127.0.0.1"), Port: 443, Zone: ""}},
	}

	sink.Sink(testResult)
	if len(sink.incomingInfoHashes) != 1 {
		t.Error("incomingInfoHashes field of Sink has not been filled in correctly")
	}

	sink.Sink(testResult)
	if len(sink.incomingInfoHashes) != 1 {
		t.Error("the same InfoHash should not be processed multiple times")
	}
}

func TestSink_Terminate(t *testing.T) {
	t.Parallel()

	sink := NewSink(time.Minute, 1)
	sink.Terminate()

	if !sink.terminated {
		t.Error("terminated field of Sink has not been set to true")
	}
}

func TestSink_Drain(t *testing.T) {
	t.Parallel()

	defer func() {
		if r := recover(); r == nil {
			t.Error("No panic while Draining an already closed Sink!")
		}
	}()

	sink := NewSink(time.Minute, 1)
	sink.Terminate()
	sink.Drain()
}

func TestFlush(t *testing.T) {
	t.Parallel()

	sink := NewSink(time.Minute, 1)
	testMetadata := Metadata{
		InfoHash: []byte{1, 2, 3, 4, 5, 6},
	}

	go func() {
		select {
		case result := <-sink.drain:
			if !reflect.DeepEqual(result.InfoHash, testMetadata.InfoHash) {
				t.Errorf("Expected flushed InfoHash to be %v, but got %v", testMetadata.InfoHash, result.InfoHash)
			}

		case <-time.After(1 * time.Second):
			t.Error("Timeout waiting for flush result")
		}
	}()

	sink.flush(testMetadata)

	time.Sleep(500 * time.Millisecond)

	var infoHash [20]byte
	copy(infoHash[:], testMetadata.InfoHash)
	sink.incomingInfoHashesMx.RLock()
	_, exists := sink.incomingInfoHashes[infoHash]
	sink.incomingInfoHashesMx.RUnlock()
	if exists {
		t.Error("InfoHash was not deleted after flush")
	}
}
