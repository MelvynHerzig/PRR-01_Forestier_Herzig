// Package clients implements server side tcp logic.
// It handles each client on a specific goroutine and manages access to hostel rooms.
package clients

import (
	"bufio"
	"fmt"
	"log"
	"net"
	config "prr.configuration/reader"
	"server/hostel"
	"server/hostel/request"
	"server/tcpserver/servers"
	"strconv"
)

// HandleClients starts client handler goroutine and hostel logic goroutine.
func HandleClients (listener net.Listener) {

	// Starting concurrent hostel manager.
	go hostelManager(config.GetRoomsCount(), config.GetNightsCount())

	// Listening for TCP connections.
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}

		// Starting handler for client connection.
		go handleConnection(conn)
	}
}

// client to server messages channel.
type client chan<- string

// Messages passed from client handler to hostelManager
type clientDemand struct {
	ch     client
	demand request.HostelRequestable
}

var (
	// leaving channel to notify clients that leave. Usually when clients shut down
	leaving = make(chan client)

	// clientDemands channel used by client handlers to transmit HostelRequestable. Need mutex access.
	clientDemands = make(chan clientDemand)
)

// hostelManager concurrency safe function that handles hostel rooms management.
func hostelManager(nbRooms, nbNights uint) {

	// All connected client with their unique name, on entering string is "".
	clients := make(map[client]string)

	hostelManager, hostelError := hostel.NewHostel(nbRooms, nbNights)
	if hostelError != nil {
		log.Fatal(hostelError)
	}

	waitingMutex := false
	
	// Handling clients.
	for {

		select {
		case clientDemand := <-clientDemands:

			// TODO call demande (Processus client -> processus mutex)
			servers.Demand()
			waitingMutex = true

			// Getting session client name ans assigning it to request.
			if _, exists := clients[clientDemand.ch]; exists {
				clientDemand.demand.SetUsername(clients[clientDemand.ch])
			}

			// TODO call attente (Processus client -> processus mutex)

			success, username, message := clientDemand.demand.Execute(hostelManager)

			// If request succeed, we replicate, we set the name associated to client.
			if success {

				// TODO replicate

				clients[clientDemand.ch] = username
				clientDemand.ch <- message
			} else {
				sendError(clientDemand.ch, message)
			}

			// TODO call fin (Processus client -> processus mutex)

		case cli := <-leaving:
			delete(clients, cli)
			close(cli)
		}
	}
}

// handleConnection handles client communication. Listen and write to clients or hostel Manager
func handleConnection(conn net.Conn) {
	ch := make(chan string) // bidirectional

	// Starting client writer
	go func() {
		for msg := range ch { // client writer <- hostelManager / handleConnection
			_, _ = fmt.Fprintln(conn, msg) // TCP Client <- client writer
		}
	}()

	strRooms  := strconv.FormatUint(uint64(config.GetRoomsCount()) , 10)
	strNights := strconv.FormatUint(uint64(config.GetNightsCount()), 10)

	ch <- "WELCOME Welcome in the FH Hostel ! Nb rooms: "  + strRooms + ", nb nights: " + strNights +
		"- LOGIN <userName>" +
		"- LOGOUT" +
		"- BOOK <roomNumber> <arrivalNight> <nbNights>" +
		"- ROOMLIST <night>" +
		"- FREEROOM <arrivalNight> <nbNights>"

	// Scanning incoming client message.
	input := bufio.NewScanner(conn)
	for input.Scan() {

		goodRequest, req := request.MakeRequest(input.Text())

		if goodRequest {
			clientDemands <- clientDemand{ch, req}
		} else {
			sendError(ch, "Unknown request")
		}
	}

	leaving <- ch
	conn.Close()
}

// sendError sends errors to client with error prefix ("ERROR").
func sendError(ch client, err string) {
	ch <- "ERROR " + err
}


