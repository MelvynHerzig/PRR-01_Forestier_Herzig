// Package tcpserver manages the tcp server
package tcpserver

import (
	"log"
	"net"
	"server/config"
	"server/tcpserver/clients"
	"server/tcpserver/sync"
)

func StartServer(serverNumber uint) {

	// TODO add to config parser method to check serverNumber


	// Starting TCP Server.
	listener, err := net.Listen("tcp", "localhost:" + config.Servers[serverNumber].Port)
	if err != nil {
		log.Fatal(err)
	}

	sync.ServersSync(listener, serverNumber)
	clients.HandleClients(listener, config.GetRoomsCount(), config.GetDaysCount())
}
