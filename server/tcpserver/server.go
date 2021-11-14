// Package tcpserver manages the tcp server start
package tcpserver

import (
	"log"
	"net"
	configReader "prr.configuration/reader"
	"server/tcpserver/clients"
	"server/tcpserver/sync/network"
	"strconv"
)

var serverNumber uint = 0

// StartServer start the server
func StartServer(noServer uint) {

	serverNumber = noServer

	// Starting TCP Server.
	listener, err := net.Listen("tcp", "localhost:" + strconv.FormatUint(uint64(configReader.GetServerById(serverNumber).Port), 10))
	if err != nil {
		log.Fatal(err)
	}

	network.ServersSync(listener)
	clients.HandleClients(listener, configReader.GetRoomsCount(), configReader.GetNightsCount())
}

// GetServerNumber gets the current server number
func GetServerNumber() uint {
	return serverNumber
}