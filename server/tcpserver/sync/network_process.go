// Package network Sync network and listen/send to other servers
package sync

import (
	"bufio"
	"fmt"
	"log"
	"net"
	configReader "prr.configuration/reader"
	"strconv"
	"strings"
)

var connections []net.Conn

// ServersSync synchronize server with the other servers. Blocks until all servers are online.
// The server try to reach servers with a smaller number and wait that server with bigger number reach it.
func ServersSync(listener net.Listener) {

	connections = make([]net.Conn, len(configReader.GetServers()))


	// Connect to servers (smaller numbers)
	localServerNumber := configReader.GetLocalServerNumber()
	for otherServ := uint(0) ; otherServ < localServerNumber ; otherServ++ {

		var conn net.Conn
		var err error

		for notConnected := true; notConnected; notConnected = err != nil {
			strPort := strconv.FormatUint(uint64(configReader.GetServerById(otherServ).Port), 10)
			conn, err = net.Dial("tcp", configReader.GetServerById(otherServ).Ip + ":" + strPort)
		}

		fmt.Fprintln(conn, strconv.FormatUint(uint64(localServerNumber), 10))

		fmt.Println("server " + strconv.FormatUint(uint64(otherServ), 10)  + " connected " )
		connections[otherServ] = conn
	}

	// Wait for others to connect ( bigger numbers)
	for otherServ := localServerNumber + 1 ; otherServ < uint(len(configReader.GetServers())) ; otherServ++ {
		var conn net.Conn
		var err error
		for isServer := false; !isServer; {
			conn, err = listener.Accept()
			if err != nil {
				log.Fatal(err)
			}
			isServer = configReader.IsServerIP(conn.RemoteAddr().String())
			if !isServer {
				fmt.Fprintln(conn, "Sync running, try again later")
				conn.Close()
			}
		}

		data, _ := bufio.NewReader(conn).ReadString('\n')
		num, _ := strconv.ParseInt(strings.Trim(data, "\n"), 10, 0)


		fmt.Println("server " + strconv.FormatUint(uint64(num), 10)  + " connected" )

		connections[num] = conn
	}

	// Start worker to listen for other servers messages
	/*for _, conn := range connections {
		go handleServer(conn)
	}*/
}

// handleServer listens to the connection for incoming messages from other servers
func handleServer(conn net.Conn) {
	input := bufio.NewScanner(conn)

	for input.Scan() {
		HandleMessage(Deserialize(input.Text()))
	}
}

// SendToOne sends a message to a dstServer with a given message type.
func SendToOne(dstServer uint, message string) {
	fmt.Fprintln(connections[dstServer], message)
}

// SendToAll sends a Massage to others servers
func SendToAll(message string) {
	for i, _ := range connections {
		if i != int(configReader.GetLocalServerNumber()) {
			SendToOne(uint(i), message)
		}
	}
}

