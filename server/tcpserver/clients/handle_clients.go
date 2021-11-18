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
	"strconv"
	"strings"
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

var (
	// leaving channel to notify clients that leave. Usually when clients shut down
	leaving = make(chan client)

	// requests channel to submit requests to job layer.
	requests = make(chan HostelRequestable)
)

// hostelManager concurrency safe function that handles hostel rooms management.
func hostelManager(nbRooms, nbNights uint) {

	// All connected client with their unique name, on entering string is "".
	clients := make(map[client]string)

	hostelManager, hostelError := hostel.NewHostel(nbRooms, nbNights)
	if hostelError != nil {
		log.Fatal(hostelError)
	}

	// Handling clients.
	for {

		select {
		case request := <-requests:

			// TODO call demande (Processus client -> processus mutex)

			// TODO call attente (Processus client -> processus mutex)

			request.execute(hostelManager, clients)

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

		goodRequest, req := makeUserRequest(conn.RemoteAddr().String(), input.Text(), ch)

		if goodRequest {
			requests <- req
		} else {
			sendError(ch, "Unknown request")
		}
	}

	leaving <- ch
	conn.Close()
}

// makeUserRequest analyzes incoming message to create request to hostel manager.
// Clients request must contain the exact amount of arguments.
func makeUserRequest(clientAddr, req string, ch client) (bool, HostelRequestable) {

	trimReq := strings.TrimSpace(req)
	splits := strings.Split(trimReq, " ")

	switch splits[0] {
	case "LOGIN":
		if len(splits) != 2 {
			break
		}
		var req loginRequest
		req.chanToHandler = ch
		req.clientAddr = clientAddr
		req.clientName = splits[1]
		return true, &req

	case "LOGOUT":
		var req logoutRequest
		req.chanToHandler = ch
		req.clientAddr = clientAddr
		return true, &req

	case "BOOK":
		if len(splits) != 4 {
			break
		}
		var req bookRequest
		req.chanToHandler = ch
		req.clientAddr = clientAddr

		roomNumber,   err1 := strconv.ParseUint(splits[1], 10, 0)
		arrivalNight, err2 := strconv.ParseUint(splits[2], 10, 0)
		nbNights,     err3 := strconv.ParseUint(splits[3], 10, 0)
		if err1 != nil || err2 != nil || err3 != nil {
			break
		}

		req.roomNumber = uint(roomNumber)
		req.nightStart = uint(arrivalNight)
		req.duration   = uint(nbNights)

		return true, &req

	case "ROOMLIST":
		if len(splits) != 2 {
			break
		}
		var req roomStateRequest
		req.chanToHandler = ch
		req.clientAddr = clientAddr

		night, err := strconv.ParseUint(splits[1], 10, 0)
		if err != nil {
			break
		}

		req.nightNumber = uint(night)

		return true, &req

	case "FREEROOM":
		if len(splits) != 3 {
			break
		}
		var req disponibilityRequest
		req.chanToHandler = ch
		req.clientAddr = clientAddr

		arrivalNight, err1 := strconv.ParseUint(splits[1], 10, 0)
		nbNights  , err2 := strconv.ParseUint(splits[2], 10, 0)
		if err1 != nil || err2 != nil  {
			break
		}

		req.nightStart = uint(arrivalNight)
		req.duration   = uint(nbNights)

		return true, &req
	}


	return false, nil
}

// sendError sends errors to client with error prefix ("ERROR").
func sendError(ch client, err string) {
	ch <- "ERROR " + err
}


