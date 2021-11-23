// Package clients implements logic to communicate with clients.
package clients

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"prr.configuration/config"
	"server/hostel"
	"server/tcpserver/servers"
)

// HandleClients starts a goroutine to handle concurrently the clients and accepts client connexions.
func HandleClients (listener net.Listener) {

	go clientsManager()

	// Listening for TCP connections.
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}

		// Starting handler for client connection.
		go clientHandler(conn)
	}
}

// client to server message channel.
type client chan<- string

// Messages passed from clientHandler to clientsManager
type clientRequest struct {
	ch      client         // channel to post response message
	request hostel.Request // Request received
}

// Channels used between clientHandler and clientsManager
var (
	// leaving channel to notify clients that leave. Usually when clients shut down
	leaving = make(chan client)

	// clientDemands channel used by client handlers to transmit Request. Need mutex access.
	clientDemands = make(chan clientRequest)
)

// clientsManager is a function to start in a goroutine. It uses the Communicating sequential process in order to handle
// client concurrently. When a request is received, it asks for mutex access before submitting it to Manager.
func clientsManager() {

	// All connected client sessions with their unique name, on entering string is "".
	clients := make(map[client]string)

	// Handling clients.
	for {
		select {

		case clientDemand := <-clientDemands:

			servers.AccessMutex() // May block

			// Getting session client name and assigning it to request. If client is not logged, will be "".
			if _, exists := clients[clientDemand.ch]; exists {
				clientDemand.request.SetUsername(clients[clientDemand.ch])
			}

			response := hostel.SubmitRequest(clientDemand.request)

			// If request succeed, we replicate the state in other servers, and we set the name used in the request
			// and in the session.
			if response.Success {

				if clientDemand.request.ShouldReplicate() {
					servers.Replicate(clientDemand.request)
				}

				clients[clientDemand.ch] = response.Username
				clientDemand.ch <- response.Message
			} else {
				sendError(clientDemand.ch, response.Message)
			}

			servers.LeaveMutex()

		case cli := <-leaving:

			// We remove client session
			delete(clients, cli)
			close(cli)

		}
	}
}

// clientHandler handles client communication. Listen and write to clients or clientsManager
func clientHandler(conn net.Conn) {
	ch := make(chan string) // bidirectional

	// Starting client writer
	go func() {
		for msg := range ch { // client writer <- clientsManager / clientHandler
			_, _ = fmt.Fprintln(conn, msg) // TCP Client <- client writer
		}
	}()

	ch <- fmt.Sprintf("WELCOME Welcome in the FH Hostel ! Nb rooms: %v, nb nights: %v" +
					"- LOGIN <userName>" +
					"- LOGOUT" +
					"- BOOK <roomNumber> <arrivalNight> <nbNights>" +
					"- ROOMLIST <night>" +
					"- FREEROOM <arrivalNight> <nbNights>", config.GetRoomsCount(), config.GetNightsCount())

	// Scanning incoming client message.
	input := bufio.NewScanner(conn)
	for input.Scan() {

		goodRequest, req := hostel.MakeRequest(input.Text())

		if goodRequest {
			clientDemands <- clientRequest{ch, req}
		} else {
			sendError(ch, "Unknown communication")
		}
	}

	leaving <- ch
	_ = conn.Close()
}

// sendError sends errors to client with error prefix ("ERROR").
func sendError(ch client, err string) {
	ch <- "ERROR " + err
}


