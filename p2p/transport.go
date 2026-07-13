package p2p

// Peer represents the remote node.
type Peer interface {
	Close() error
}

// Transport handles communication between nodes in the network.
// This can be of form (TCP, UDP, Websockets).
type Transport interface {
	ListenAndAccept() error
	Consume() <-chan RPC
}
