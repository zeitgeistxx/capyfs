package p2p

import "net"

// Peer represents the remote node.
type Peer interface {
	net.Conn
	Send([]byte) error
}

// Transport handles communication between nodes in the network.
// This can be of form (TCP, UDP, Websockets).
type Transport interface {
	Dial(string) error
	ListenAndAccept() error
	Consume() <-chan RPC
	Close() error
}
