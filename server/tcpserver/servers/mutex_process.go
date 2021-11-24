// Package servers handles access to critical section between servers
package servers

import (
	"prr.configuration/config"
	"server/tcpserver"
	"server/tcpserver/servers/clock"
)

// Channels used for communication between client process and mutex process to manage critical section access
var (
	demand = make(chan struct{})
	leave  = make(chan struct{})
	allow  = make(chan struct{}) // to awake a client process waiting for critical section
)

// Channel used to get received message from network process.
var onMessage = make(chan message)

// Stores the messages of all servers
var messages []message

// The server isInSC
var isInSC bool

// MutexCore method executed in a goroutine. It is responsible to be the mutex "engine".
func MutexCore() {

	// Stores the messages of all servers
	messages = make([]message, len(config.GetServers()))

	// Init the messages
	for i := 0; i < len(messages); i++ {
		messages[i] = message{
			MessageType: REL,
			Timestamp:   0,
			SrcServer:   uint(i),
		}
	}

	for {
		select {
		case _ = <-demand:
			doDemand()
		case _ = <-leave:
			doLeave()
		case msg := <-onMessage:
			doHandleMessage(msg)
		}
	}
}

// AccessMutex function called by client process to signal that it wants critical section access. May block.
func AccessMutex() {
	tcpserver.LogMutex("--------- Asking ---------")
	demand <- struct{}{}
	tcpserver.LogMutex("--------- Waiting ---------")
	<-allow
	tcpserver.LogMutex("--------- Entering ---------")

}

// LeaveMutex function called by client process to signal that it finished critical section
func LeaveMutex() {
	leave <- struct{}{}
	tcpserver.LogMutex("--------- Leaving ---------")
}

// HandleMessage function called by network process to signal that a message arrived.
func handleMessage(message message) {
	onMessage <- message
}

// doDemand function called when client process ask for mutex
func doDemand() {
	clock.IncTimestamp()
	message := message{MessageType: REQ,
		Timestamp: clock.GetTimestamp(),
		SrcServer: config.GetLocalServerNumber()}
	messages[config.GetLocalServerNumber()] = message
	sendToAll(serialize(message))

}

// doLeave function called when client process leaves the mutex
func doLeave() {
	clock.IncTimestamp()
	message := message{MessageType: REL,
		Timestamp: clock.GetTimestamp(),
		SrcServer: config.GetLocalServerNumber()}
	messages[config.GetLocalServerNumber()] = message
	sendToAll(serialize(message))
	isInSC = false
}

// doHandleMessage handles a mutex message incoming from distant server
func doHandleMessage(msg message) {
	clock.SyncTimestamp(msg.Timestamp)

	// Because we have implemented optimised version.
	// All incoming message can be stored
	messages[msg.SrcServer] = msg

	// Optimisation
	if msg.MessageType == REQ && messages[config.GetLocalServerNumber()].MessageType != REQ {
		response := message{MessageType: ACK,
			Timestamp: clock.GetTimestamp(),
			SrcServer: config.GetLocalServerNumber()}
		sendToOne(msg.SrcServer, serialize(response))
	}

	checkCriticalSection()
}

// checkCriticalSection checks if server has access to mutex
func checkCriticalSection() {
	if messages[config.GetLocalServerNumber()].MessageType != REQ || isInSC {
		return
	}
	myNumber := config.GetLocalServerNumber()
	for i := 0; i < len(messages); i++ {
		if uint(i) == myNumber {
			continue
		}

		// if our timestamp is not the oldest
		if messages[myNumber].Timestamp > messages[i].Timestamp {
			return
		} else if messages[myNumber].Timestamp == messages[i].Timestamp && myNumber > uint(i) {
			return
		}
	}

	// Signal client that he can enter
	allow <- struct{}{}
	isInSC = true
}
