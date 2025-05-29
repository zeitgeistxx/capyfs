package main

import (
	"log"

	"github.com/zeitgeistxx/capyfs/p2p"
)

func main() {
	tr := p2p.NewTCPTransport(":8000")
	
	if err := tr.ListenAndAccept(); err != nil {
		log.Fatal(err)
	}

	select {}
}
