// Package main starts the hostel server manager
package main

import (
	"fmt"
	"log"
	"os"
	"server/tcpserver"
	"server/tcpserver/clients"
	"strconv"
)

// main Gets programs arguments, configuration file and start server.
func main() {

	// Getting program args (server number).
	argsLen := len(os.Args)
	if argsLen < 2 {
		log.Fatal("Paramètres manquants")
	}

	noServ, errNoServ   := strconv.ParseUint(os.Args[1], 10, 0)

	if  errNoServ != nil {
		log.Fatal("Paramètre invalide. Devrait être <no serveur> <options>")
	}

	// Fetching options
	if argsLen == 3 {
		for argIndex := 2 ; argIndex < argsLen; argIndex++ {

			switch os.Args[argIndex] {
			case "-debug":
				clients.DebugMode = true
				fmt.Println("Starting with debug on")

			default: log.Fatal("Argument inconnu ", os.Args[argIndex])
			}
		}
	}

	// TODO parse config
	// something like config.Parse()

	tcpserver.StartServer(uint(noServ))
}

