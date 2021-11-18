package main

import (
	"client/tcpclient"
	"os"
	config "prr.configuration/reader"
	"strconv"
)

// main starts hostel tcp client.
func main() {

	config.InitSimple("../config.json")

	var server *config.Server
	// Getting program args (server address).
	if len(os.Args) < 2 {
		server = config.GetServerRandomly()
	} else {
		id, _ := strconv.ParseUint(os.Args[1], 10, 0)
		server = config.GetServerById(uint(id))
	}

	tcpclient.StartClient(server)
}
