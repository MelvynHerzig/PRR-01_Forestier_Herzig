package request

import (
	"fmt"
	"strconv"
	"strings"
)

// MakeRequest analyzes incoming message to create request to hostel manager.
// Clients request must contain the exact amount of arguments.
func MakeRequest(req string) (bool, HostelRequestable) {

	trimReq := strings.TrimSpace(req)
	splits := strings.Split(trimReq, " ")

	switch splits[0] {
	case "LOGIN":
		if len(splits) < 2 {
			break
		}
		var req loginRequest
		req.clientName = splits[1]
		if len(splits) > 2 { req.SetUsername(splits[2]) }

		return true, &req

	case "LOGOUT":
		var req logoutRequest
		if len(splits) > 1 { req.SetUsername(splits[1]) }
		return true, &req

	case "BOOK":
		if len(splits) < 4 {
			break
		}
		var req bookRequest

		roomNumber,   err1 := strconv.ParseUint(splits[1], 10, 0)
		arrivalNight, err2 := strconv.ParseUint(splits[2], 10, 0)
		nbNights,     err3 := strconv.ParseUint(splits[3], 10, 0)
		if err1 != nil || err2 != nil || err3 != nil {
			break
		}

		req.roomNumber = uint(roomNumber)
		req.nightStart = uint(arrivalNight)
		req.duration   = uint(nbNights)
		if len(splits) > 4 { req.SetUsername(splits[4]) }

		return true, &req

	case "ROOMLIST":
		if len(splits) < 2 {
			break
		}
		var req roomStateRequest

		night, err := strconv.ParseUint(splits[1], 10, 0)
		if err != nil {
			break
		}

		req.nightNumber = uint(night)
		if len(splits) > 2 { req.SetUsername(splits[2]) }
		return true, &req

	case "FREEROOM":
		if len(splits) < 3 {
			break
		}
		var req disponibilityRequest

		arrivalNight, err1 := strconv.ParseUint(splits[1], 10, 0)
		nbNights  , err2 := strconv.ParseUint(splits[2], 10, 0)
		if err1 != nil || err2 != nil  {
			break
		}

		req.nightStart = uint(arrivalNight)
		req.duration   = uint(nbNights)

		if len(splits) > 3 { req.SetUsername(splits[3]) }

		return true, &req
	}

	return false, nil
}

func (r *loginRequest) Serialize() string{
	return "LOGIN " + r.clientName
}

func (r *logoutRequest) Serialize() string{
	return "LOGOUT " + r.username
}

func (r *bookRequest) Serialize() string{
	return fmt.Sprintf("BOOK %d %d %d %s", r.roomNumber, r.nightStart, r.duration, r.username)
}

func (r *roomStateRequest) Serialize() string{
	return fmt.Sprintf("ROOMLIST %d %s", r.nightNumber, r.username)
}

func (r *disponibilityRequest) Serialize() string{
	return fmt.Sprintf("FREEROOM %d %d %s", r.nightStart, r.duration, r.username)
}
