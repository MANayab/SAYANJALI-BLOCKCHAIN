package p2p

import (
	"io"
	"testing"
)

func TestHelloFinishRoundTrip(t *testing.T) {
	v := HelloFinish{NodeID: "node-1", PublicKey: make([]byte, 64), EchoChallenge: make([]byte, 32), Signature: make([]byte, 64)}
	b, err := EncodeHelloFinish(v, 7)
	if err != nil {
		t.Fatal(err)
	}
	f, err := ReadFrame(bytesReader(b))
	if err != nil {
		t.Fatal(err)
	}
	got, err := DecodeHelloFinish(f)
	if err != nil {
		t.Fatal(err)
	}
	if got.NodeID != v.NodeID {
		t.Fatal("node id changed")
	}
	if string(got.PublicKey) != string(v.PublicKey) || string(got.EchoChallenge) != string(v.EchoChallenge) || string(got.Signature) != string(v.Signature) {
		t.Fatal("HELLO_FINISH payload changed during round trip")
	}
}

// bytesReader keeps this test independent of network I/O.
func bytesReader(b []byte) *sliceReader { return &sliceReader{b: b} }

type sliceReader struct {
	b []byte
	i int
}

func (r *sliceReader) Read(p []byte) (int, error) {
	if r.i >= len(r.b) {
		return 0, io.EOF
	}
	n := copy(p, r.b[r.i:])
	r.i += n
	return n, nil
}
