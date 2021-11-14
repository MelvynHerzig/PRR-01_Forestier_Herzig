package sync

import (
	"strconv"
	"strings"
)

// type of message
const (
	ACK string = "ack"
	REQ        = "req"
	REL        = "rel"
)

// Message is the struct passed between servers and network/mutex processes
type Message struct {
	MessageType string
	Timestamp uint
	SrcServer uint
}

// Serialize function that serializes a message into a string.
func Serialize(m Message) string {
	strTimestamp  := strconv.FormatUint(uint64(m.Timestamp), 10)
	strSrcServer  := strconv.FormatUint(uint64(m.SrcServer), 10)
	return m.MessageType + " " + strTimestamp + " " + strSrcServer
}

// Deserialize function that deserializes a string into a message.
func Deserialize(content string) Message {
	splits := strings.Split(content, " ")

	var receivedMessageType string

	switch splits[0] {
	case ACK: receivedMessageType = ACK
	case REQ: receivedMessageType = REQ
	case REL: receivedMessageType = REL
	}

	receivedTimestamp, _ := strconv.ParseUint(splits[1], 10, 0)
	receivedSrcServer, _ := strconv.ParseUint(splits[2], 10, 0)

	return Message{
		MessageType: receivedMessageType,
		Timestamp:   uint(receivedTimestamp),
		SrcServer:   uint(receivedSrcServer),
	}
}