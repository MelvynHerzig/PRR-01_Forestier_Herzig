// Package servers implements some custom messages in order to do mutex negotiation.
package servers

import (
	"strconv"
	"strings"
)

// type of message to negotiate mutex.
const (
	req string = "req"
	token      = "token"
)

// message is the struct passed between servers and network/mutex processes
type message struct {
	MessageType string
	SrcServer uint
}

// serialize function that serializes a message into a string.
func serialize(m message) string {
	strSrcServer  := strconv.FormatUint(uint64(m.SrcServer), 10)
	return m.MessageType + " " + strSrcServer
}

// deserialize function that deserializes a string into a message.
func deserialize(content string) message {
	splits := strings.Split(content, " ")

	var receivedMessageType string

	switch splits[0] {
	case token: receivedMessageType = token
	case req  : receivedMessageType = req
	}

	receivedSrcServer, _ := strconv.ParseUint(splits[1], 10, 0)

	return message{
		MessageType: receivedMessageType,
		SrcServer:   uint(receivedSrcServer),
	}
}