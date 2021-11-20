// Package main checks program arguments, get the server from config and then connects to it by starting client.
package main

import (
	"client/tcpclient"
	"log"
	"os"
	"prr.configuration/config"
	"strconv"
)

// main starts hostel tcp client.
func main() {

	config.InitSimple("../config.json")

	var server *config.Server

	// Getting program args (server number)
	if len(os.Args) < 2 {
		server = config.GetServerRandomly()
	} else {
		id, err := strconv.ParseUint(os.Args[1], 10, 0)
		if err != nil || int(id) >= len(config.GetServers()){
			log.Fatal("Server number must be an integer between [0, Servers count[.")
		}
		server = config.GetServerById(uint(id))
	}

	tcpclient.StartClient(server)
}
