package p2p

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTCPTransport(t *testing.T) {
	opts := TCPTransportOpts{
		ListenAddrress: ":8000",
		HandshakeFunc:  NOPHandshakeFunc,
		Decoder:        DefaultDecoder{},
	}
	tr := NewTCPTransport(opts)

	assert.Equal(t, tr.ListenAddrress, ":8000")
	assert.Nil(t, tr.ListenAndAccept())
}
