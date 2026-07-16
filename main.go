package main

import (
	"bytes"
	"log"
	"time"

	"github.com/zeitgeistxx/capyfs/p2p"
)

func makeServer(listenAddr string, nodes ...string) *FileServer {
	tcpTransportOpts := p2p.TCPTransportOpts{
		ListenAddr:    listenAddr,
		HandshakeFunc: p2p.NOPHandshakeFunc,
		Decoder:       p2p.DefaultDecoder{},
	}
	tcpTransport := p2p.NewTCPTransport(tcpTransportOpts)

	fileServerOpts := FileServerOpts{
		StorageRoot:       listenAddr + "_network",
		PathTransformFunc: CASPathTransformFunc,
		Transport:         tcpTransport,
		BootstrapNodes:    nodes,
	}

	s := NewFileServer(fileServerOpts)

	tcpTransport.OnPeer = s.OnPeer

	return s
}

func main() {
	s1 := makeServer(":8000", "")
	s2 := makeServer(":8001", ":8000")

	go func() {
		log.Fatal(s1.Start())
	}()

	time.Sleep(time.Second * 2)

	go s2.Start()
	time.Sleep(time.Second * 2)

	data := bytes.NewReader([]byte("my data file!"))
	s2.StoreData("verysecret", data)

	select {}
}
