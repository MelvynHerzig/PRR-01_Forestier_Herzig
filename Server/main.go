package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

func main() {

	if len(os.Args) < 3 {
		log.Fatal("Paramètres manquants")
		return
	}

	/*NB_ROOMS, errRooms := strconv.Atoi(os.Args[1])
	NB_DAYS, errDays := strconv.Atoi(os.Args[2])

	if errRooms != nil || errDays != nil {
		log.Fatal("Paramètres invalides")
		return
	}*/

	listener, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}

	//go roomManager(NB_ROOMS, NB_DAYS)
go broadcaster()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go handleConn(conn)
	}
}

type client chan<- string // an outgoing message channel

var (
	entering = make(chan client)
	leaving  = make(chan client)
	requests = make(chan string)
	messages = make(chan string)// all incoming client messages
)

type User struct {
	Username   string
	Connection net.Conn
}

func roomManager(nbRooms, nbDays int) {
	clients := make(map[client]bool)

	// 2d slice for reservations
	reservations := make([][]string, nbRooms)
	for room := range reservations {
		reservations[room] = make([]string, nbDays)
	}

	for {
		select {
		case req := <-requests:

			params := strings.Split(req, " ")
			// switch sur req, et faire en fonction de la requête
			switch params[0] {

			// Book a room (roomNumber, arrivalDay, nbNights)
			case "BOOK":
				fmt.Println("BOOK")

				// Get occupationList (day)
			case "ROOMLIST":
				fmt.Println("ROOMLIST")

				// Get free room (arrivalDay, nbNights)
			case "FREEROOM":
				fmt.Println("FREEROOM")
			}

		case cli := <-entering:
			clients[cli] = true

		case cli := <-leaving:
			delete(clients, cli)
			close(cli)
		}
	}
}

func handleConnection(conn net.Conn) {

}

func broadcaster() {
	clients := make(map[client]bool) // all connected clients
	for {
		select {
		case msg := <-messages: // broadcaster <- handleConn
			// Broadcast incoming message to all
			// clients' outgoing message channels.
			for cli := range clients {
				cli <- msg // clientwriter (handleConn) <- broadcaster
			}

		case cli := <-entering:
			clients[cli] = true

		case cli := <-leaving:
			delete(clients, cli)
			close(cli)
		}
	}
}

func handleConn(conn net.Conn) {
	ch := make(chan string) // channel 'client' mais utilisé ici dans les 2 sens
	go func() {             // clientwriter
		for msg := range ch { // clientwriter <- broadcaster, handleConn
			fmt.Fprintln(conn, msg) // netcat Client <- clientwriter
		}
	}()

	who := conn.RemoteAddr().String()
	ch <- "You are " + who // clientwriter <- handleConn
	entering <- ch
	messages <- who + " has arrived" // broadcaster <- handleConn

	input := bufio.NewScanner(conn)
	for input.Scan() { // handleConn <- netcat client
		messages <- who + ": " + input.Text() // broadcaster <- handleConn
	}

	leaving <- ch
	messages <- who + " has left" // broadcaster <- handleConn
	conn.Close()
}
