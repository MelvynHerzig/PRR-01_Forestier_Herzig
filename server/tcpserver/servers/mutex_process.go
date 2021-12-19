// Package servers handles access to critical section between servers
package servers

import (
	"container/list"
	"prr.configuration/config"
	"server/tcpserver"
)

// Channels used for communication between client process and mutex process to manage critical section access
var (
	demand = make(chan struct{})
	leave  = make(chan struct{})
	allow  = make(chan struct{}) // to awake a client process waiting for critical section
)

// Channel used to get received message from network process.
var onMessage = make(chan message)

// Stores the processes' id that have transmitted a message
var queue = list.New()

// Defines the state of the mutex noDemand (= 0), demanding (= 1), inSection (= 2)
var state int

const noDemand = 0
const demanding = 1
const inSection = 2

// MutexCore method executed in a goroutine. It is responsible to be the mutex "engine".
func MutexCore() {

	for {
		select {
		case _ = <-demand:
			doDemand()
		case _ = <-leave:
			doLeave()
		case msg := <-onMessage:
			switch msg.MessageType {
			case token:
				doHandleToken(msg)
			case req:
				doHandleReq(msg)
			}
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
	// We are a children we ask for token.
	if parent.id != config.GetLocalServerNumber() {
		queue.PushBack(config.GetLocalServerNumber())

		if state == noDemand {
			state = demanding
			sendToParent(serialize(message{MessageType: req,
				SrcServer: config.GetLocalServerNumber()}))
		}
	} else { // we are root we going in.
		state = inSection
		allow <- struct{}{}
	}

}

// doLeave function called when client process leaves the mutex
func doLeave() {
	state = noDemand
	if queue.Len() != 0 {

		// Getting next req.
		front := queue.Front()
		queue.Remove(front)
		// Child become parent
		parent = children[front.Value.(uint)]
		delete(children, front.Value.(uint))

		sendToParent(serialize(message{MessageType: token,
			SrcServer: config.GetLocalServerNumber()}))

		if queue.Len() != 0 {
			state = demanding
			sendToParent(serialize(message{MessageType: req,
				SrcServer: config.GetLocalServerNumber()}))

		}
	}
}

// doHandleReq handles an incoming req message
func doHandleReq(msg message) {
	if parent.id == config.GetLocalServerNumber() && state == noDemand {
		// Child become parent
		parent = children[msg.SrcServer]
		delete(children, msg.SrcServer)

		sendToParent(serialize(message{MessageType: token,
			SrcServer: config.GetLocalServerNumber()}))
	} else {
		queue.PushBack(msg.SrcServer)
		if parent.id != config.GetLocalServerNumber() && state == noDemand {
			state = demanding
			sendToParent(serialize(message{MessageType: req,
				SrcServer: config.GetLocalServerNumber()}))
		}
	}
}

// doHandleToken handles an incoming token message
func doHandleToken(msg message) {
	// Getting next req.
	front := queue.Front()
	queue.Remove(front)

	if front.Value.(uint) == config.GetLocalServerNumber() {
		children[parent.id] = parent
		parent = treeNode{id: front.Value.(uint), connection: nil}
		state = inSection
		allow <- struct{}{}
	} else {
		// Child becomes parent and parent becomes child
		temp := parent
		parent = children[front.Value.(uint)]
		delete(children, front.Value.(uint))
		children[temp.id] = temp

		sendToParent(serialize(message{MessageType: token,
			SrcServer: config.GetLocalServerNumber()}))

		if queue.Len() != 0 {
			state = demanding
			sendToParent(serialize(message{MessageType: req,
				SrcServer: config.GetLocalServerNumber()}))
		} else {
			state = noDemand
		}
	}
}
