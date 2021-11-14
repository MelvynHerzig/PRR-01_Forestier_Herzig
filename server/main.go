// Package main starts the hostel server manager
package main

import (
	"log"
	"os"
	configReader "prr.configuration/reader"
	"server/tcpserver"
	"strconv"
)

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
	configReader.Init("../config.json")

	if noServ < 0 || noServ >= uint64(len(configReader.GetServers())) {
		log.Fatal("Le numero de serveur doit être entre [0, nb serveurs [")
	}

	tcpserver.StartServer(uint(noServ))
}

