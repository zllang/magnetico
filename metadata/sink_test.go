package metadata

import (
	"net"
	"testing"
	"time"
)

func TestRandomDigit(t *testing.T) {
	t.Parallel()
	for i := 0; i < 100; i++ {
		digit := randomDigit()
		if digit < '0' || digit > '9' {
			t.Errorf("Expected a digit in range(0 - 9), got %c", digit)
		}
	}
}

func TestPeerId(t *testing.T) {
	t.Parallel()
	for i := 0; i < 100; i++ {
		peerID := randomID()
		lenPeerID := len(peerID)
		if lenPeerID > PeerIDLength {
			t.Errorf("peerId longer than 20 bytes: %s (%d)", peerID, lenPeerID)
		}
	}
}

func TestSink_NewSink(t *testing.T) {
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
	sink := NewSink(time.Minute, 1)
	sink.Terminate()

	if !sink.terminated {
		t.Error("terminated field of Sink has not been set to true")
	}
}
