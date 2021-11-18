// Package main starts the hostel server manager
package main

import (
	"log"
	"net"
	"os"
	config "prr.configuration/reader"
	"server/tcpserver/clients"
	"server/tcpserver/servers"
	"strconv"
)

var ServerNumber uint = 0

// main Gets programs arguments, configuration file and start server.
func main() {

	// Getting program args (server number).
	argsLen := len(os.Args)
	if argsLen < 2 {
		log.Fatal("Un paramètre attendu: <no de server>")
	}

	noServ, errNoServ   := strconv.ParseUint(os.Args[1], 10, 0)
	if  errNoServ != nil {
		log.Fatal("Paramètre invalide. Devrait être <no serveur>")
	}

	// Init configuration
	config.Init("../config.json", uint(noServ))

	if noServ < 0 || noServ >= uint64(len(config.GetServers())) {
		log.Fatal("Le numero de serveur doit être entre [0, nb serveurs [")
	}

	// Starting TCP Server.
	localPort := config.GetServerById(config.GetLocalServerNumber()).Port
	listener, err := net.Listen("tcp", "localhost:" + strconv.FormatUint(uint64(localPort), 10))
	if err != nil {
		log.Fatal(err)
	}

	servers.ServersSync(listener)
	clients.HandleClients(listener)
}

