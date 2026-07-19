package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/zeitgeistxx/capyfs/p2p"
)

func makeServer(listenAddr string, nodes ...string) *FileServer {
	tcptransportOpts := p2p.TCPTransportOpts{
		ListenAddr:    listenAddr,
		HandshakeFunc: p2p.NOPHandshakeFunc,
		Decoder:       p2p.DefaultDecoder{},
	}
	tcpTransport := p2p.NewTCPTransport(tcptransportOpts)

	fileServerOpts := FileServerOpts{
		EncKey:            newEncryptionKey(),
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
	s2 := makeServer(":8001", "")
	s3 := makeServer(":8002", ":8000", ":8001")

	go func() { log.Fatal(s1.Start()) }()
	time.Sleep(time.Millisecond * 500)
	go func() { log.Fatal(s2.Start()) }()

	time.Sleep(time.Second * 2)

	go s3.Start()
	time.Sleep(time.Second * 2)

	for i := range 20 {
		key := fmt.Sprintf("capybara_%d.png", i)
		data := bytes.NewReader([]byte("ehstoph!! capybara no ugly"))
		s3.Store(key, data)

		if err := s3.store.Delete(s3.ID, key); err != nil {
			log.Fatal(err)
		}

		r, err := s3.Get(key)
		if err != nil {
			log.Fatal(err)
		}

		b, err := io.ReadAll(r)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(string(b))
	}
}
