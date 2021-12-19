// Package servers handles servers to servers connexions: broadcast, replication, mutex negotiation.
package servers

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"prr.configuration/config"
	"server/hostel"
	"server/tcpserver"
	"strconv"
	"strings"
)

// replicationDone channel used to signal that a replication communication confirmation has been received
// Used between Replicate function and serverHandler
var replicationDone = make(chan struct{})

// confirmReplication is the message to confirm that replication has been done.
const confirmReplication = "OK"

type treeNode struct {
	id         uint
	connection net.Conn
}

var children = make(map[uint]treeNode)

var parent treeNode

// Sync synchronize server with the other servers. It blocks until all servers are online and a connexion could
// be established. The server try to reach servers with a smaller number and wait that server with bigger
// number reach it.
func Sync(listener net.Listener) {

	localServerNumber := config.GetLocalServerNumber()

	// Waiting for children to connect
	for range config.GetInitialChildrenIds(localServerNumber) {
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
		num, _ := strconv.ParseUint(strings.Trim(data, "\n"), 10, 0)

		fmt.Println("server " + strconv.FormatUint(uint64(num), 10) + " connected")

		children[uint(num)] = treeNode{id: uint(num), connection: conn}
	}

	parentId := config.GetServerById(localServerNumber).Parent
	// If we have a parent
	if parentId != localServerNumber {
		// Connect to parent
		var conn net.Conn
		var err error

		parentId := config.GetServerById(localServerNumber).Parent
		strPort := strconv.FormatUint(uint64(config.GetServerById(parentId).Port), 10)

		for notConnected := true; notConnected; notConnected = err != nil {
			conn, err = net.Dial("tcp", config.GetServerById(parentId).Ip+":"+strPort)
		}

		parent = treeNode{id: parentId, connection: conn}
		fmt.Println("Parent " + strconv.FormatUint(uint64(parentId), 10) + " connected ")

		// Sending our local number
		sendToParent(strconv.FormatUint(uint64(localServerNumber), 10))

		// Waiting for ready signal ( it doesn't matter what it is, we just want a message back)
		bufio.NewReader(conn).ReadString('\n')
	} else {
		parent = treeNode{id: localServerNumber, connection: nil}
	}

	sendToChildren("GO")

	// Start worker to listen for other servers messages
	for _, node := range children {
		go serverHandler(node.connection)
	}
	if parentId != localServerNumber {
		go serverHandler(parent.connection)
	}

	fmt.Printf("Server %v ready to handle clients\n", localServerNumber)
}

// Replicate sends the given requests to other server and waits that they answer back that they have applied it.
// Blocks until other all server confirmed.
func Replicate(req hostel.Request) {

	sendToChildren(req.Serialize())

	// Waiting for other servers to confirm that the replication has been processed
	for nbResponseReceived := 0; nbResponseReceived < len(children); nbResponseReceived++ {
		<-replicationDone
	}

	// if we have a parent we notify that our children and us are done
	sendToParent(confirmReplication)
}

// serverHandler listens the connection for incoming messages from other servers
func serverHandler(conn net.Conn) {
	input := bufio.NewScanner(conn)

	for input.Scan() {
		text := input.Text()
		tcpserver.LogServer(fmt.Sprintf("%s received", text))
		if text == confirmReplication { // If replication confirmation

			replicationDone <- struct{}{}

		} else if good, request := hostel.MakeRequest(text, true); good { // If replication requested
			hostel.SubmitRequest(request)
			tcpserver.LogRequestReplicating(request)
			Replicate(request)

		} else { // or mutex access management
			handleMessage(deserialize(input.Text()))
		}
	}
}

// sendToParent sends a message to the parent server of number.
func sendToParent(message string) {
	if parent.connection != nil {
		fmt.Fprintln(parent.connection, message)

		tcpserver.LogServer(fmt.Sprintf("%s to server %d", message, parent.id))
	}
}

// sendToChildren sends a message to children servers
func sendToChildren(message string) {
	childrenStr := "["
	for _, tn := range children {
		fmt.Fprintln(tn.connection, message)
		childrenStr += fmt.Sprintf(" %d", tn.id)
	}
	childrenStr += " ]"

	tcpserver.LogServer(fmt.Sprintf("%s to servers %s", message, childrenStr))
}
