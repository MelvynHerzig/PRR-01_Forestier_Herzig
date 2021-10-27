package sync

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"server/config"
	"strconv"
	"strings"
)

var connections []net.Conn

// ServersSync synchronize server with the other servers. Blocks until all servers are online.
// The server try to reach servers with a smaller number and wait that server with bigger number reach it.
func ServersSync(listener net.Listener, noServer uint) {

	connections = make([]net.Conn, len(config.Servers))

	// Connect to servers (smaller numbers)
	for otherServ := uint(0) ; otherServ < noServer ; otherServ++ {

		var conn net.Conn
		var err error

		for notConnected := true; notConnected; notConnected = err != nil {
			conn, err = net.Dial("tcp", config.Servers[otherServ].Ip+ ":" + config.Servers[otherServ].Port)
		}

		fmt.Fprintln(conn, strconv.FormatUint(uint64(noServer), 10))

		fmt.Println("server " + strconv.FormatUint(uint64(otherServ), 10)  + " connected " )
		connections[otherServ] = conn
	}

	// Wait for others to connect ( bigger numbers)
	for otherServ := noServer + 1 ; otherServ < uint(len(config.Servers)) ; otherServ++ {
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


		fmt.Println("server " + strconv.FormatUint(uint64(num), 10)  + " connected" )

		connections[num] = conn
	}
}


