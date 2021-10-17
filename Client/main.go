package main

import (
	"Client/tcpclient"
	"log"
	"os"
)

// main starts hostel tcp client.
func main() {

	// Getting program args (server address).
	if len(os.Args) < 2 {
		log.Fatal("Paramètres manquants")
	}

	tcpclient.StartClient(os.Args[1])
}
