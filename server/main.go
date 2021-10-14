package main

import (
	"Server/tcpserver"
	"fmt"
	"log"
	"os"
	"strconv"
)

// main Starts hostel tcp server.
func main() {
	// Getting program args (#rooms and #day of hostel).
	argsLen := len(os.Args)
	if argsLen < 3 {
		log.Fatal("Paramètres manquants")
	}

	nbRooms, errRooms := strconv.ParseUint(os.Args[1], 10, 0)
	nbDays, errDays   := strconv.ParseUint(os.Args[2], 10, 0)

	if errRooms != nil || errDays != nil {
		log.Fatal("Paramètres invalides. Devrait être <nb chambres> <nb nuits>")
	}

	// Fetching last arguments
	if argsLen == 4 {
		for argIndex := 3 ; argIndex < argsLen; argIndex++ {

			switch os.Args[argIndex] {
			case "-debug":
				tcpserver.DebugMode = true
				fmt.Println("Starting with debug on")

			default: log.Fatal("Argument inconnu ", os.Args[argIndex])
			}
		}
	}

	tcpserver.StartServer(uint(nbRooms), uint(nbDays))
}

