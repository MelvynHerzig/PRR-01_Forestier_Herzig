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

	if len(os.Args) < 3 {
		log.Fatal("Paramètres manquants")
		return
	}

	NB_ROOMS, errRooms := strconv.ParseUint(os.Args[1], 10, 0)
	NB_DAYS, errDays := strconv.ParseUint(os.Args[2], 10, 0)

	if errRooms != nil || errDays != nil {
		log.Fatal("Paramètres invalides")
		return
	}

	listener, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}

	go roomManager(uint(NB_ROOMS), uint(NB_DAYS))
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go handleConnection(conn)
	}
}

type client chan<- string // an outgoing message channel

var (
	login    = make(chan *UserRequest)
	quit     = make(chan *User)
	requests = make(chan *UserRequest)
)

type User struct {
	Username   string
	Channel    client
	Connection net.Conn
}

type UserRequest struct {
	User    *User
	Request Request
}

type Request struct {
	Command string
	Params  []string
}

func roomManager(nbRooms, nbDays uint) {
	clients := make(map[*User]bool)

	// 2d slice for reservations
	reservations := make([][]string, nbRooms)
	for room := range reservations {
		reservations[room] = make([]string, nbDays)
	}

	hostel, hostelError := hostel.NewHostel(nbRooms, nbDays)

	if hostelError != nil {
		log.Fatal(hostelError)
	}

	for {
		select {
		case req := <-requests:

			// switch sur req, et faire en fonction de la requête
			switch req.Request.Command {

			case "LOGIN":
				// TODO: Check if username is already connected and return error if true
				if hostel.RegisterClient(req.User.Username) {

					req.User.Username = req.Request.Params[0]
					req.User.Channel <- "Hello " + req.User.Username + "\n" +
						"You can book our rooms (1-3) for days (1-5)\n" +
						"Here is the list of commands you can use :\n" +
						"- BOOK roomNumber, arrivalDay, nbNights\n" +
						"- ROOMLIST day\n" +
						"- FREEROOM arrivalDay nbNights\n" +
						"- QUIT"
				}
			// Book a room (roomNumber, arrivalDay, nbNights)
			case "BOOK":
				roomNumber, _ := strconv.ParseUint(req.Request.Params[0], 10, 0)
				arrivalDay, _ := strconv.ParseUint(req.Request.Params[1], 10, 0)
				nbNights, _ := strconv.ParseUint(req.Request.Params[2], 10, 0)

				hostel.Book(req.User.Username, uint(roomNumber), uint(arrivalDay), uint(nbNights))

				// Get occupationList (day)
			case "ROOMLIST":
				day, _ := strconv.ParseUint(req.Request.Params[0], 10, 0)

				rooms, _ := hostel.GetRoomsState(req.User.Username, uint(day))

				var result string

				for index, _ := range rooms {
					result += fmt.Sprintf("%s ", hostel.RoomsStateSignification[index])
				}

				req.User.Channel <- result

				// Get free room (arrivalDay, nbNights)
			case "FREEROOM":
				arrivalDay, _ := strconv.ParseUint(req.Request.Params[0], 10, 0)
				nbNights, _ := strconv.ParseUint(req.Request.Params[1], 10, 0)
				val, error := hostel.SearchDisponibility(uint(arrivalDay), uint(nbNights))

				if error != nil {
					log.Print(error)
				} else {
					req.User.Channel <- strconv.FormatUint(uint64(val), 10)
				}
			}

		case cli := <-login:
			clients[cli.User] = hostel.RegisterClient(cli.User.Username)

		case cli := <-quit:
			delete(clients, cli)
			close(cli.Channel)
		}
	}
}

func handleConnection(conn net.Conn) {
	ch := make(chan string) // channel 'client' mais utilisé ici dans les 2 sens
	go func() {             // clientwriter
		for msg := range ch { // clientwriter <- broadcaster, handleConn
			fmt.Fprintln(conn, msg) // netcat Client <- clientwriter
		}
	}()

	var user User
	user.Channel = ch
	user.Connection = conn
	user.Channel <- "Welcome in the FH Hotel !\n"

	input := bufio.NewScanner(conn)

	doLogin(&user, input)

	for input.Scan() {

		req := makeUserRequest(input.Text(), &user)

		requests <- &req
	}

	quit <- &user
	conn.Close()
}

func makeUserRequest(req string, user *User) UserRequest {
	split := strings.Split(req, " ")
	var uReq UserRequest
	uReq.Request.Command = split[0]
	uReq.Request.Params = split[1:]
	uReq.User = user

	return uReq
}

func doLogin(user *User, input *bufio.Scanner) {
	user.Channel <- "Please, identify yourself. Use command \"LOGIN {your name}\""
	for input.Scan() {
		uReq := makeUserRequest(input.Text(), user)

		if uReq.Request.Command != "LOGIN" {
			user.Channel <- "Please, identify yourself. Use command \"LOGIN {your name}\""
			continue
		}

		requests <- &uReq

		fmt.Println(uReq.User.Username)
		fmt.Println(user.Username)
		fmt.Println("-------------")

		if user.Username != "" {
			return
		}
	}
}
