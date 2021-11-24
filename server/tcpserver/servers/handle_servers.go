// Package servers handles servers to servers connexions, servers synchronization and servers to servers communication.
package servers

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"prr.configuration/config"
	"server/hostel"
	"strconv"
	"strings"
)

// replicationDone channel used to signal that a replication communication confirmation has been received
// Used between Replicate function and serverHandler
var replicationDone = make(chan struct{})

// confirmReplication is the message to confirm that replication has been done.
const confirmReplication = "OK"

var connections []net.Conn

// Sync synchronize server with the other servers. It blocks until all servers are online and a connexion could
// be established. The server try to reach servers with a smaller number and wait that server with bigger
// number reach it.
func Sync(listener net.Listener) {

	connections = make([]net.Conn, len(config.GetServers()))

	// Connect to servers (smaller numbers)
	localServerNumber := config.GetLocalServerNumber()
	for otherServ := uint(0); otherServ < localServerNumber; otherServ++ {

		var conn net.Conn
		var err error

		for notConnected := true; notConnected; notConnected = err != nil {
			strPort := strconv.FormatUint(uint64(config.GetServerById(otherServ).Port), 10)
			conn, err = net.Dial("tcp", config.GetServerById(otherServ).Ip+":"+strPort)
		}

		fmt.Fprintln(conn, strconv.FormatUint(uint64(localServerNumber), 10))

		fmt.Println("server " + strconv.FormatUint(uint64(otherServ), 10) + " connected ")
		connections[otherServ] = conn
	}

	// Wait for others to connect ( bigger numbers)
	for otherServ := localServerNumber + 1; otherServ < uint(len(config.GetServers())); otherServ++ {
		var conn net.Conn
		var err error
		for isServer := false; !isServer; {
			conn, err = listener.Accept()
			if err != nil {
				log.Fatal(err)
			}
			isServer = config.IsServerIP(conn.RemoteAddr().String())
			if !isServer {
				fmt.Fprintln(conn, "Sync running, try again later")
				conn.Close()
			}
		}

		data, _ := bufio.NewReader(conn).ReadString('\n')
		num, _ := strconv.ParseInt(strings.Trim(data, "\n"), 10, 0)

		fmt.Println("server " + strconv.FormatUint(uint64(num), 10) + " connected")

		connections[num] = conn
	}

	// Start worker to listen for other servers messages
	for conn := 0; conn < len(connections); conn++ {
		if uint(conn) != localServerNumber {
			go serverHandler(connections[conn])
		}
	}
}

// Replicate sends the given requests to other server and waits that they answer back that they have applied it.
// Blocks until other all server confirmed.
func Replicate(req hostel.Request){

	sendToAll(req.Serialize())

	// Waiting for other servers to confirm that the replication has been processed
	for nbResponseReceived := 0; nbResponseReceived < len(connections) - 1; nbResponseReceived++ {
		<-replicationDone
	}
}

// serverHandler listens the connection for incoming messages from other servers
func serverHandler(conn net.Conn) {
	input := bufio.NewScanner(conn)

	for input.Scan() {
		text := input.Text()

		if text == confirmReplication { // If replication confirmation

			replicationDone <- struct{}{}

		} else if good, request := hostel.MakeRequest(text, true); good {  // If replication requested

			hostel.SubmitRequest(request)
			fmt.Fprintln(conn, confirmReplication)

		} else { // or mutex access management

			handleMessage(deserialize(input.Text()))
		}
	}
}

// sendToOne sends a message to a distant server of number dstServer with a given message type.
func sendToOne(dstServer uint, message string) {
	fmt.Fprintln(connections[dstServer], message)
}

// sendToAll sends a message to others servers
func sendToAll(message string) {
	for i := range connections {
		if i != int(config.GetLocalServerNumber()) {
			sendToOne(uint(i), message)
		}
	}
}

