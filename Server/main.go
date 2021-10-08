package main

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

func main() {

	// Getting program args (#rooms and #day of hostel).
	if len(os.Args) < 3 {
		log.Fatal("Paramètres manquants")
	}

	NbRooms, errRooms := strconv.ParseUint(os.Args[1], 10, 0)
	NbDays, errDays   := strconv.ParseUint(os.Args[2], 10, 0)

	if errRooms != nil || errDays != nil {
		log.Fatal("Paramètres invalides")
	}

	// Opening TCP server.
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

type client chan<- string // Server to client message channel.

var (
	 leaving = make(chan client)
	entering = make(chan client)
	requests = make(chan hostelRequestable) // all incoming client request
)

func hostelManager(nbRooms, nbDays uint) {
	// All connected client with their unique name
	clients := make(map[client]string)

	hostelManager, hostelError := hostel.NewHostel(nbRooms, nbDays)
	if hostelError != nil {
		log.Fatal(hostelError)
	}

	for {
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
}

func handleConnection(conn net.Conn) {
	ch := make(chan string) // bidirectional

	// Starting client writer
	go func() {
		for msg := range ch { // client writer <- hostelManager / handleConnection
			_, _ = fmt.Fprintln(conn, msg) // TCP Client <- client writer
		}
	}()

	ch <- "Welcome in the FH Hotel ! \n" +
		  "First log in with: LOGIN userName " +
	      "Then, here the commands you can use once legged in :\n" +
		  "- BOOK roomNumber arrivalDay nbNights\n" +
		  "- ROOMLIST day\n" +
		  "- FREEROOM arrivalDay nbNights \n" +
		  "- LOGOUT"

	input := bufio.NewScanner(conn)

	for input.Scan() {

		goodRequest, req := makeUserRequest(input.Text(), ch)

		if goodRequest {
			requests <- req
		}
	}

	leaving <- ch
	conn.Close()
}

func makeUserRequest(req string, ch client) (bool, hostelRequestable) {

	trimmReq := strings.TrimSpace(req)
	splits := strings.Split(trimmReq, " ")

	switch splits[0] {
	case "LOGIN":
		if len(splits) != 2 {
			break
		}
		var req registerRequest
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
		var req getRoomStateRequest
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
		var req searchDisponibilityRequest
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

// Interface of request that can be made to hostel.
type hostelRequestable interface {
	execute(h *hostel.Hostel, clients map[client]string)
}

type hostelRequest struct {
	chanToHandler client
}

type registerRequest struct {
	hostelRequest
	clientName string
}

type logoutRequest struct {
	hostelRequest
}

type bookRequest struct {
	hostelRequest
	roomNumber uint
	dayStart uint
	duration uint
}

type getRoomStateRequest struct {
	hostelRequest
	dayNumber uint
}

type searchDisponibilityRequest struct {
	hostelRequest
	dayStart uint
	duration uint
}

func (r *registerRequest) execute(h *hostel.Hostel, clients map[client]string) {

	if clients[r.chanToHandler] != "" {
		r.chanToHandler <- "You are already connected as " + clients[r.chanToHandler] + "."
		return
	}

	clients[r.chanToHandler] = r.clientName
	h.TryRegister(r.clientName)

	r.chanToHandler <- "Connection success."
}

func (r *logoutRequest) execute(h *hostel.Hostel, clients map[client]string) {

	if clients[r.chanToHandler] == "" {
		r.chanToHandler <- "You are not logged in."
		return
	}

	clients[r.chanToHandler] = ""
	r.chanToHandler <- "Logged out successfully."
}

func (r *bookRequest) execute(h *hostel.Hostel, clients map[client]string) {

	if clients[r.chanToHandler] == "" {
		r.chanToHandler <- "You are not logged in."
		return
	}

	if err := h.Book(clients[r.chanToHandler], r.roomNumber, r.dayStart, r.duration); err != nil {
		r.chanToHandler <- err.Error()
		return
	}

	r.chanToHandler <- "You succcessfully booked room " + strconv.FormatUint(uint64(r.roomNumber), 10) + "."
}

func (r *getRoomStateRequest) execute(h *hostel.Hostel, clients map[client]string) {

	if clients[r.chanToHandler] == "" {
		r.chanToHandler <- "You are not logged in."
		return
	}

	states, err := h.GetRoomsState(clients[r.chanToHandler], r.dayNumber)
	if err != nil {
		r.chanToHandler <- err.Error()
		return
	}

	var res string

	for _, state := range states {
		res += fmt.Sprintf("%s ", hostel.RoomsStateSignification[state])
	}

	r.chanToHandler <- res
}

func (r *searchDisponibilityRequest) execute(h *hostel.Hostel, clients map[client]string) {

	if clients[r.chanToHandler] == "" {
		r.chanToHandler <- "You are not logged in."
		return
	}

	room, err := h.SearchDisponibility(r.dayStart, r.duration)
	if err != nil {
		r.chanToHandler <- err.Error()
		return
	}

	r.chanToHandler <- "Room " + strconv.FormatUint(uint64(room), 10) + " is free."
}