// Package mutex handles access to critical section between servers
package mutex

import (
	configReader "prr.configuration/reader"
	"server/tcpserver"
	"server/tcpserver/sync"
	"server/tcpserver/sync/clock"
	"server/tcpserver/sync/network"
)

// Channels used for communication between client process and mutex process to manage critical section access
var (
	demand  chan struct{}
	leave   chan struct{}
	allow   chan struct{} // to awake a client process waiting for critical section
)

// Channel used to get received mesage from network process.
var onMessage chan sync.Message

// Stores the messages of all servers
var messages []sync.Message

// mutexCore method executed in a go routine. It is responsible to be the mutex "engine".
// Other thread can use the channels to pass
func mutexCore() {

	// Stores the messages of all servers
    messages = make([]sync.Message, len(configReader.GetServers()))

	// Init the messages
	for i := 0; i < len(messages); i++ {
		messages[i] = sync.Message{
			MessageType: sync.REL,
			Timestamp:   0,
			SrcServer:   uint(i),
		}
	}

	for {

		select {
			case <-demand:
			case <-leave:
			case <-onMessage:
		}
	}
}

// Demand function called by client process to signal that it wants critical section access.
func Demand() {
	demand <- struct{}{}
}

// Waiting function called by client process to signal that it waits for critical section access
func Waiting() {
	<-allow
}

// Leave function called by client process to signal that it finished critical section
func Leave() {
	leave <- struct{}{}
}

// HandleMessage function called by network process to signal that a message arrived.
func HandleMessage(message sync.Message) {
	onMessage <- message
}

// doDemand function called when client process ask for mutex
func doDemand() {
	clock.IncTimestamp()
	message := sync.Message{ MessageType: sync.REQ,
		                     Timestamp: clock.GetTimestamp(),
		                     SrcServer: tcpserver.GetServerNumber()}
	messages[tcpserver.GetServerNumber()] = message
	network.SendToAll(sync.Serialize(message))
}

// doLeave function called when client process leaves the mutex
func doLeave() {
	clock.IncTimestamp()
	message := sync.Message{ MessageType: sync.REL,
							 Timestamp: clock.GetTimestamp(),
							 SrcServer: tcpserver.GetServerNumber()}
	messages[tcpserver.GetServerNumber()] = message
	network.SendToAll(sync.Serialize(message))
}

// doHandleMessage handles a mutex message incoming from distant server
func doHandleMessage (msg sync.Message) {
	clock.SyncTimestamp(msg.Timestamp)

	// Because we have implemented optimised version.
	// All incoming message can be stored
	messages[msg.SrcServer] = msg

	// Optimisation
	if msg.MessageType == sync.REQ && messages[tcpserver.GetServerNumber()].MessageType != sync.REQ {
		response := sync.Message{ MessageType: sync.ACK,
								  Timestamp: clock.GetTimestamp(),
								  SrcServer: tcpserver.GetServerNumber()}
		network.SendToOne(msg.SrcServer, sync.Serialize(response))
	}

	checkCriticalSection()
}

// checkCriticalSection checks if server has access to mutex
func checkCriticalSection() {
	if messages[tcpserver.GetServerNumber()].MessageType != sync.REQ {
		return
	}

	myNumber := tcpserver.GetServerNumber()

	for i := 0; i < len(messages); i++ {
		if uint(i) == myNumber {
			continue
		}

		// if our timestamp is not the oldest
		if messages[myNumber].Timestamp > messages[i].Timestamp {
			return
		} else if  messages[myNumber].Timestamp == messages[i].Timestamp && myNumber > uint(i) {
			return
		}
	}

	// Signal client that he can enter
	allow<- struct{}{}
}