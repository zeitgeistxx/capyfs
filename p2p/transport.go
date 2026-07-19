package p2p

import "net"

// Peer represents the remote node.
type Peer interface {
	net.Conn
	Send([]byte) error
	CloseStream()
}

// Transport handles communication between nodes in the network.
// This can be of form (TCP, UDP, Websockets).
type Transport interface {
	Addr() string
	Dial(string) error
	ListenAndAccept() error
	Consume() <-chan RPC
	Close() error
}
