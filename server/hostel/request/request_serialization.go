package request

import (
	"strconv"
	"strings"
)

// makeUserRequest analyzes incoming message to create request to hostel manager.
// Clients request must contain the exact amount of arguments.
func makeUserRequest(req string) (bool, HostelRequestable) {

	trimReq := strings.TrimSpace(req)
	splits := strings.Split(trimReq, " ")
	
	var hostelRequestable HostelRequestable

	switch splits[0] {
	case "LOGIN":
		if len(splits) < 2 {
			break
		}
		var req loginRequest
		req.clientName = splits[1]
		if len(splits) > 2 { req.setUsername(splits[2]) }

		hostelRequestable = &req

	case "LOGOUT":
		var req logoutRequest
		if len(splits) > 1 { req.setUsername(splits[1]) }
		hostelRequestable = &req

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
		if len(splits) > 4 { req.setUsername(splits[4]) }

		hostelRequestable = &req

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
		if len(splits) > 2 { req.setUsername(splits[2]) }
		hostelRequestable = &req

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

		if len(splits) > 3 { req.setUsername(splits[3]) }

		hostelRequestable = &req

	default:
		return false, nil
	}

	return true, hostelRequestable
}

func (r *loginRequest) serialize() string{
	return "biop bop"
}

func (r *logoutRequest) serialize() string{
	return "biop bop"
}

func (r *bookRequest) serialize() string{
	return "biop bop"
}

func (r *roomStateRequest) serialize() string{
	return "biop bop"
}

func (r *disponibilityRequest) serialize() string{
	return "biop bop"
}
