// Package servers implements some custom messages in order to do mutex negotiation.
package servers

import (
	"strconv"
	"strings"
)

// type of message to negotiate mutex.
const (
	ACK string = "ack"
	REQ        = "req"
	REL        = "rel"
)

// message is the struct passed between servers and network/mutex processes
type message struct {
	MessageType string
	Timestamp uint
	SrcServer uint
}

// serialize function that serializes a message into a string.
func serialize(m message) string {
	strTimestamp  := strconv.FormatUint(uint64(m.Timestamp), 10)
	strSrcServer  := strconv.FormatUint(uint64(m.SrcServer), 10)
	return m.MessageType + " " + strTimestamp + " " + strSrcServer
}

// deserialize function that deserializes a string into a message.
func deserialize(content string) message {
	splits := strings.Split(content, " ")

	var receivedMessageType string

	switch splits[0] {
	case ACK: receivedMessageType = ACK
	case REQ: receivedMessageType = REQ
	case REL: receivedMessageType = REL
	}

	receivedTimestamp, _ := strconv.ParseUint(splits[1], 10, 0)
	receivedSrcServer, _ := strconv.ParseUint(splits[2], 10, 0)

	return message{
		MessageType: receivedMessageType,
		Timestamp:   uint(receivedTimestamp),
		SrcServer:   uint(receivedSrcServer),
	}
}