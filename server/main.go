// Package main gets program arguments, starts the client handler, starts the mutex core, sync the servers and
// finally starts listener for client connexions.
package main

import (
	"log"
	"net"
	"os"
	"prr.configuration/config"
	"server/hostel"
	"server/tcpserver/clients"
	"server/tcpserver/servers"
	"strconv"
)

// main Gets programs arguments, configuration file and start server.
func main() {

	// Getting program args (server number).
	argsLen := len(os.Args)
	if argsLen < 2 {
		log.Fatal("Usage: <server number>")
	}

	noServ, errNoServ   := strconv.ParseUint(os.Args[1], 10, 0)
	if  errNoServ != nil {
		log.Fatal("Invalid parameter. Must be <no serveur>")
	}

	// Init configuration
	config.Init("../config.json", uint(noServ))

	if noServ < 0 || noServ >= uint64(len(config.GetServers())) {
		log.Fatal("Server number is an integer between [0, servers count [")
	}

	// Starting TCP Server.
	localPort := config.GetServerById(config.GetLocalServerNumber()).Port
	listener, err := net.Listen("tcp", "localhost:" + strconv.FormatUint(uint64(localPort), 10))
	if err != nil {
		log.Fatal(err)
	}

	// Hostel manager to handle sequentially incoming request.
	go hostel.Manager(config.GetRoomsCount(), config.GetNightsCount())

	// Starting mutex core that client processes used to access hostel
	go servers.MutexCore()

	// Waiting for servers to come online and opening connexions.
	servers.Sync(listener)

	// Waiting for client.
	clients.HandleClients(listener)
}

