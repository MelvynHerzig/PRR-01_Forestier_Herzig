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

type Message struct {
	messageType string
	timestamp uint
	srcServer uint
}

func (m *Message) serialize() string {
	strTimestamp  := strconv.FormatUint(uint64(m.timestamp), 10)
	strSrcServer  := strconv.FormatUint(uint64(m.srcServer), 10)
	return m.messageType + " " + strTimestamp + " " + strSrcServer
}

func deserialize(content string) Message {
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
		messageType: receivedMessageType,
		timestamp:   uint(receivedTimestamp),
		srcServer:   uint(receivedSrcServer),
	}
}