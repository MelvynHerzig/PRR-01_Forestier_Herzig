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
	"os"
	"strconv"
	"strings"
)

// StartServer Launch TCP Server, starts client handler goroutine and hostel logic goroutine.
func StartServer() {

	// Getting program args (#rooms and #day of hostel).
	if len(os.Args) < 3 {
		log.Fatal("Paramètres manquants")
	}

	NbRooms, errRooms := strconv.ParseUint(os.Args[1], 10, 0)
	NbDays, errDays   := strconv.ParseUint(os.Args[2], 10, 0)

	if errRooms != nil || errDays != nil {
		log.Fatal("Paramètres invalides")
	}

	// Opening TCP Server.
	listener, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}

	// Starting concurrent hostel manager.
	go hostelManager(uint(NbRooms), uint(NbDays))

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

// client Server to client message channel.
type client chan<- string

var (
	// leaving Channel to notify clients that leave. Usually when clients shut down
	leaving = make(chan client)

	// entering Channel to notify clients that enter/connect. Usually when clients launch.
	entering = make(chan client)

	// requests Channel to submit requests to job layer.
	requests = make(chan hostelRequestable)
)

// hostelManager Concurrency safe function that handle hostel rooms management.
func hostelManager(nbRooms, nbDays uint) {

	// All connected client with their unique name, on entering string is "".
	clients := make(map[client]string)

	hostelManager, hostelError := hostel.NewHostel(nbRooms, nbDays)
	if hostelError != nil {
		log.Fatal(hostelError)
	}

	// Handling clients.
	for {
		select {
		case request := <-requests:
			request.execute(hostelManager, clients)

		case cli := <-entering:
			clients[cli] = ""

		case cli := <-leaving:
			delete(clients, cli)
			close(cli)
		}
	}
}

// handleConnection Client handler function. Listen and write to clients or hostel Manager
func handleConnection(conn net.Conn) {
	ch := make(chan string) // bidirectional

	// Starting client writer
	go func() {
		for msg := range ch { // client writer <- hostelManager / handleConnection
			_, _ = fmt.Fprintln(conn, msg) // TCP Client <- client writer
		}
	}()

	ch <- "WELCOME Welcome in the FH Hotel !" +
		"- LOGIN userName" +
		"- LOGOUT" +
		"- BOOK roomNumber arrivalDay nbNights" +
		"- ROOMLIST day" +
		"- FREEROOM arrivalDay nbNights"

	// Scanning incoming client message.
	input := bufio.NewScanner(conn)
	for input.Scan() {

		goodRequest, req := makeUserRequest(input.Text(), ch)

		if goodRequest {
			requests <- req
		} else {
			sendError(ch, "Unknown request")
		}
	}

	leaving <- ch
	conn.Close()
}

// makeUserRequest Analyze incoming message to create request to hostel manager.
// Clients request must contain the exact amount of argument.
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
		arrivalDay, err2 := strconv.ParseUint(splits[2], 10, 0)
		nbNights  , err3 := strconv.ParseUint(splits[3], 10, 0)
		if err1 != nil || err2 != nil || err3 != nil {
			break
		}

		req.roomNumber = uint(roomNumber)
		req.dayStart   = uint(arrivalDay)
		req.duration   = uint(nbNights)

		return true, &req

	case "ROOMLIST":
		if len(splits) != 2 {
			break
		}
		var req roomStateRequest
		req.chanToHandler = ch

		day, err := strconv.ParseUint(splits[1], 10, 0)
		if err != nil {
			break
		}

		req.dayNumber = uint(day)

		return true, &req

	case "FREEROOM":
		if len(splits) != 3 {
			break
		}
		var req disponibilityRequest
		req.chanToHandler = ch

		arrivalDay, err1 := strconv.ParseUint(splits[1], 10, 0)
		nbNights  , err2 := strconv.ParseUint(splits[2], 10, 0)
		if err1 != nil || err2 != nil  {
			break
		}

		req.dayStart   = uint(arrivalDay)
		req.duration   = uint(nbNights)

		return true, &req
	}


	return false, nil
}

// sendError Sends to client error with error prefix ("ERROR").
func sendError(ch client, err string) {
	ch <- "ERROR " + err
}


