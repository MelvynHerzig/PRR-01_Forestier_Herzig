// Package hostel implements Request function to serialize/deserialize .
package hostel

import (
	"fmt"
	"strconv"
	"strings"
)

// MakeRequest analyzes incoming message to create a hostelRequestable.
// Returns false, nil, if req was not sufficient to form a hostelRequestable.
func MakeRequest(str string, mayContainUsername bool) (bool, Request) {

	trimReq := strings.TrimSpace(str)
	splits := strings.Split(trimReq, " ")

	// Checking if there is good amount of arguments
	minArgs, maxArgs := 0, 0
	switch splits[0] {
	case "LOGIN": 	 minArgs, maxArgs = 2, 2
	case "LOGOUT": 	 minArgs, maxArgs = 1, 1
	case "BOOK":     minArgs, maxArgs = 4, 4
	case "ROOMLIST": minArgs, maxArgs = 2, 2
	case "FREEROOM": minArgs, maxArgs = 3, 3
	}

	if mayContainUsername { maxArgs++ }
	if len(splits) < minArgs || len(splits) > maxArgs {
		return false, nil
	}

	// Forming request
	var responseReq Request

	switch splits[0] {
	case "LOGIN":
		var req loginRequest
		req.clientName = splits[1]
		responseReq = &req

	case "LOGOUT":
		var req logoutRequest
		if len(splits) > 1 { req.SetUsername(splits[1]) }
		responseReq = &req

	case "BOOK":
		var req bookRequest

		roomNumber,   err1 := strconv.ParseUint(splits[1], 10, 0)
		arrivalNight, err2 := strconv.ParseUint(splits[2], 10, 0)
		nbNights,     err3 := strconv.ParseUint(splits[3], 10, 0)
		if err1 != nil || err2 != nil || err3 != nil {
			return false, nil
		}

		req.roomNumber = uint(roomNumber)
		req.nightStart = uint(arrivalNight)
		req.duration   = uint(nbNights)
		responseReq = &req

	case "ROOMLIST":
		var req roomStateRequest

		night, err := strconv.ParseUint(splits[1], 10, 0)
		if err != nil {
			return false, nil
		}

		req.nightNumber = uint(night)
		responseReq = &req

	case "FREEROOM":
		var req disponibilityRequest

		arrivalNight, err1 := strconv.ParseUint(splits[1], 10, 0)
		nbNights  , err2 := strconv.ParseUint(splits[2], 10, 0)
		if err1 != nil || err2 != nil  {
			return false, nil
		}

		req.nightStart = uint(arrivalNight)
		req.duration   = uint(nbNights)

		responseReq = &req
	}

	if mayContainUsername && len(splits) == maxArgs { responseReq.SetUsername(splits[maxArgs - 1]) }

	return true, responseReq
}

// Serialize transform a loginRequest into a string for communication protocol.
func (r *loginRequest) Serialize() string{
	return "LOGIN " + r.clientName
}

// Serialize transform a logoutRequest into a string for communication protocol.
func (r *logoutRequest) Serialize() string{
	return "LOGOUT " + r.username
}

// Serialize transform a bookRequest into a string for communication protocol.
func (r *bookRequest) Serialize() string{
	return fmt.Sprintf("BOOK %d %d %d %s", r.roomNumber, r.nightStart, r.duration, r.username)
}

// Serialize transform a roomStateRequest into a string for communication protocol.
func (r *roomStateRequest) Serialize() string{
	return fmt.Sprintf("ROOMLIST %d %s", r.nightNumber, r.username)
}

// Serialize transform a disponibilityRequest into a string for communication protocol.
func (r *disponibilityRequest) Serialize() string{
	return fmt.Sprintf("FREEROOM %d %d %s", r.nightStart, r.duration, r.username)
}
