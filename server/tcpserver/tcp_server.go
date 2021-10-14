// Package tcpserver implements the Server TCP logic.
// The tcp_server.go file is responsible for TCP "worker" and
// the request.go file is responsible for requests that can
// be submitted to the job logic layer.
package tcpserver

import (
	"Server/hostel"
	"bufio"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
)

// StartServer Launch TCP Server, starts client handler goroutine and hostel logic goroutine.
func StartServer(nbRooms, nbNights uint) {

	// Opening TCP Server.
	listener, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}

	// Starting concurrent hostel manager.
	go hostelManager(uint(nbRooms), uint(nbNights))

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

// client to server message channel.
type client chan<- string

var (
	// leaving channel to notify clients that leave. Usually when clients shut down
	leaving = make(chan client)

	// requests channel to submit requests to job layer.
	requests = make(chan hostelRequestable)
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

		if DebugMode {
			debugModeLog("--------- Enter shared zone ---------")
		}

		select {
		case request := <-requests:

			request.execute(hostelManager, clients)

		case cli := <-leaving:
			delete(clients, cli)
			close(cli)
		}

		if DebugMode {
			debugModeLog("--------- Leave shared zone ---------")
		}
	}
}

// handleConnection handles client communication. Listen and write to clients or hostel Manager
func handleConnection(conn net.Conn) {
	ch := make(chan string) // bidirectional

	// Starting client writer
	go func() {
		for msg := range ch { // client writer <- hostelManager / handleConnection

			if DebugMode {
				debugModeLog("To " + conn.RemoteAddr().String() + " : " + msg)
			}

			_, _ = fmt.Fprintln(conn, msg) // TCP Client <- client writer
		}
	}()

	ch <- "WELCOME Welcome in the FH Hotel !" +
		"- LOGIN userName" +
		"- LOGOUT" +
		"- BOOK roomNumber arrivalNight nbNights" +
		"- ROOMLIST night" +
		"- FREEROOM arrivalNight nbNights"

	// Scanning incoming client message.
	input := bufio.NewScanner(conn)
	for input.Scan() {

		goodRequest, req := makeUserRequest(input.Text(), ch)

		if DebugMode {
			debugModeLog("From " + conn.RemoteAddr().String() + " : " + input.Text())
		}

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
func makeUserRequest(req string, ch client) (bool, hostelRequestable) {

	trimReq := strings.TrimSpace(req)
	splits := strings.Split(trimReq, " ")

	switch splits[0] {
	case "LOGIN":
		if len(splits) != 2 {
			break
		}
		var req loginRequest
		req.chanToHandler = ch
		req.clientName = splits[1]
		return true, &req

	case "LOGOUT":
		var req logoutRequest
		req.chanToHandler = ch
		return true, &req

	case "BOOK":
		if len(splits) != 4 {
			break
		}
		var req bookRequest
		req.chanToHandler = ch

		roomNumber, err1 := strconv.ParseUint(splits[1], 10, 0)
		arrivalNight, err2 := strconv.ParseUint(splits[2], 10, 0)
		nbNights  , err3 := strconv.ParseUint(splits[3], 10, 0)
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


