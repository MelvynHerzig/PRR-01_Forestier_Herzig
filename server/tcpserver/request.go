// Package main implements the Server TCP logic.
// The main.go file is responsible for TCP "worker" and
// the request.go file is responsible for request that can
// be submitted to the job logic layer.
package tcpserver

import (
	"Server/hostel"
	"strconv"
)


// hostelRequestable Interface of request that can be made to hostel.
type hostelRequestable interface {
	execute(h *hostel.Hostel, clients map[client]string)
}

// hostelRequest Base content of all kind of request.
type hostelRequest struct {
	chanToHandler client
}

// loginRequest Request that log clients in, in order to book, inspect and search free room.
// Clients names are supposed to be unique. Two clients with the same name would be considered (from Server pov) as
// the same person. So their actions would impact each other.
type loginRequest struct {
	hostelRequest
	clientName string
}

// logoutRequest Request that log client out, in order to change "account".
type logoutRequest struct {
	hostelRequest
}

// bookRequest Request used to book a room.
type bookRequest struct {
	hostelRequest
	roomNumber uint
	dayStart uint
	duration uint
}

// roomStateRequest Request used to get the state of rooms for a given day.
type roomStateRequest struct {
	hostelRequest
	dayNumber uint
}

// disponiblityRequest Request used to find a free room for a given day and duration.
type disponibilityRequest struct {
	hostelRequest
	dayStart uint
	duration uint
}

// execute If user is not already logged in, log him in, otherwise explains error. All is notified to clients channel.
func (r *loginRequest) execute(h *hostel.Hostel, clients map[client]string) {

	if clients[r.chanToHandler] != "" {
		sendError(r.chanToHandler, "You are already connected as " + clients[r.chanToHandler] + ".")
		return
	}

	clients[r.chanToHandler] = r.clientName
	h.TryRegister(r.clientName)

	r.chanToHandler <- "RESULT_LOGIN"
}

// execute If user is logged in, log him out, otherwise explains error. All is notified to clients channel.
func (r *logoutRequest) execute(h *hostel.Hostel, clients map[client]string) {

	if clients[r.chanToHandler] == "" {
		sendError(r.chanToHandler, "You are not logged in.")
		return
	}

	clients[r.chanToHandler] = ""
	r.chanToHandler <- "RESULT_LOGOUT"
}

// execute Try booking a room. Client must be logged in. All is notified to clients channel.
func (r *bookRequest) execute(h *hostel.Hostel, clients map[client]string) {

	if clients[r.chanToHandler] == "" {
		sendError(r.chanToHandler, "You are not logged in.")
		return
	}

	if err := h.Book(clients[r.chanToHandler], r.roomNumber, r.dayStart, r.duration); err != nil {
		sendError(r.chanToHandler,  err.Error())
		return
	}

	strRoom     := strconv.FormatUint(uint64(r.roomNumber), 10)
	strDay      := strconv.FormatUint(uint64(r.dayStart), 10)
	strDuration := strconv.FormatUint(uint64(r.duration), 10)

	r.chanToHandler <- "RESULT_BOOK " + strRoom + " " + strDay + " " + strDuration
}

// execute Try getting room state for a given day. Client must be logged in. All is notified to clients channel.
func (r *roomStateRequest) execute(h *hostel.Hostel, clients map[client]string) {

	if clients[r.chanToHandler] == "" {
		sendError(r.chanToHandler, "You are not logged in.")
		return
	}

	states, err := h.GetRoomsState(clients[r.chanToHandler], r.dayNumber)
	if err != nil {
		sendError(r.chanToHandler, err.Error())
		return
	}

	var res string

	for i, state := range states {
		if i != 0 {
			res += ","
		}
		res += state
	}

	r.chanToHandler <- "RESULT_ROOMLIST " + res
}

// execute Try finding a free room. Client must be logged in. All is notified to clients channel.
func (r *disponibilityRequest) execute(h *hostel.Hostel, clients map[client]string) {

	if clients[r.chanToHandler] == "" {
		sendError(r.chanToHandler, "You are not logged in.")
		return
	}

	room, err := h.SearchDisponibility(r.dayStart, r.duration)
	if err != nil {
		sendError(r.chanToHandler, err.Error())
		return
	}

	strRoom     := strconv.FormatUint(uint64(room), 10)
	strDay      := strconv.FormatUint(uint64(r.dayStart), 10)
	strDuration := strconv.FormatUint(uint64(r.duration), 10)

	r.chanToHandler <- "RESULT_FREEROOM " + strRoom + " " + strDay + " " + strDuration
}
