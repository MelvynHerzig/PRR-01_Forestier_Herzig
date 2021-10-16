// Package tcpserver implements the requests that can
// be submitted to the hostel logic layer in tcp_server.go.
package tcpserver

import (
	"Server/hostel"
	"strconv"
	"time"
)

// hostelRequestable interface of request that can be made to hostel.
type hostelRequestable interface {
	execute(h *hostel.Hostel, clients map[client]string) bool
	toString() string
}

// hostelRequest base content of all kind of request.
type hostelRequest struct {
	chanToHandler client
	clientAddr string
}

// loginRequest request that logs clients in, in order to book, inspect and search free room.
// Clients names are supposed to be unique. Two clients with the same name would be considered (from Server pov) as
// the same person. So their actions would impact each other.
type loginRequest struct {
	hostelRequest
	clientName string
}

// logoutRequest request that log client out, in order to change "account".
type logoutRequest struct {
	hostelRequest
}

// bookRequest request used to book a room.
type bookRequest struct {
	hostelRequest
	roomNumber uint
	nightStart uint
	duration uint
}

// roomStateRequest request used to get the state of rooms for a given night.
type roomStateRequest struct {
	hostelRequest
	nightNumber uint
}

// disponiblityRequest Request used to find a free room for a given night and duration.
type disponibilityRequest struct {
	hostelRequest
	nightStart uint
	duration   uint
}

// execute if user is not already logged in, log him in, otherwise explains error. All is notified to clients channel.
func (r *loginRequest) execute(h *hostel.Hostel, clients map[client]string) bool{

	if clients[r.chanToHandler] != "" {
		sendError(r.chanToHandler, "You are already connected as " + clients[r.chanToHandler] + ".")
		return false
	}

	clients[r.chanToHandler] = r.clientName
	h.TryRegister(r.clientName)

	nbLoggedClient++

	r.chanToHandler <- "RESULT_LOGIN"

	if nbLoggedClient == 2 {
		debugLogRisk("Server request handler suspended. Resume in 20s.")
		time.Sleep(20 * time.Second)
		debugLogRisk("Server request handler resumed.")
	}

	return true
}

// toString converts loginRequest to string.
func (r *loginRequest) toString() string {
	return "From " + r.clientAddr + " login as " + r.clientName
}

// execute if user is logged in, log him out, otherwise explains error. All is notified to clients channel.
func (r *logoutRequest) execute(h *hostel.Hostel, clients map[client]string) bool {

	if clients[r.chanToHandler] == "" {
		sendError(r.chanToHandler, "You are not logged in.")
		return false
	}

	clients[r.chanToHandler] = ""

	nbLoggedClient--

	r.chanToHandler <- "RESULT_LOGOUT"
	return true
}

// toString converts logoutRequest to string.
func (r *logoutRequest) toString() string {
	return "From " + r.clientAddr + " logout"
}

// execute tries to book a room. Client must be logged in. All is notified to clients channel.
func (r *bookRequest) execute(h *hostel.Hostel, clients map[client]string) bool{

	if clients[r.chanToHandler] == "" {
		sendError(r.chanToHandler, "You are not logged in.")
		return false
	}

	if err := h.Book(clients[r.chanToHandler], r.roomNumber, r.nightStart, r.duration); err != nil {
		sendError(r.chanToHandler,  err.Error())
		return false
	}

	strRoom     := strconv.FormatUint(uint64(r.roomNumber), 10)
	strNight    := strconv.FormatUint(uint64(r.nightStart), 10)
	strDuration := strconv.FormatUint(uint64(r.duration), 10)

	r.chanToHandler <- "RESULT_BOOK " + strRoom + " " + strNight + " " + strDuration
	return true
}

// toString converts bookRequest to string.
func (r *bookRequest) toString() string {

	strRoom     := strconv.FormatUint(uint64(r.roomNumber), 10)
	strNight    := strconv.FormatUint(uint64(r.nightStart), 10)
	strDuration := strconv.FormatUint(uint64(r.duration), 10)

	return "From " + r.clientAddr + " BOOK room " + strRoom + " from night" + strNight + " for " + strDuration + " night(s)"
}

// execute tries to get room state for a given night. Client must be logged in. All is notified to clients channel.
func (r *roomStateRequest) execute(h *hostel.Hostel, clients map[client]string) bool {

	if clients[r.chanToHandler] == "" {
		sendError(r.chanToHandler, "You are not logged in.")
		return false
	}

	states, err := h.GetRoomsState(clients[r.chanToHandler], r.nightNumber)
	if err != nil {
		sendError(r.chanToHandler, err.Error())
		return false
	}

	var res string

	for i, state := range states {
		if i != 0 {
			res += ","
		}
		res += state
	}

	r.chanToHandler <- "RESULT_ROOMLIST " + res
	return true
}

// toString converts roomStateRequest to string.
func (r *roomStateRequest) toString() string {
	return "From " + r.clientAddr + " ROOMLIST for night" + strconv.FormatUint(uint64(r.nightNumber), 10)
}

// execute tries to find a free room. Client must be logged in. All is notified to clients channel.
func (r *disponibilityRequest) execute(h *hostel.Hostel, clients map[client]string) bool {

	if clients[r.chanToHandler] == "" {
		sendError(r.chanToHandler, "You are not logged in.")
		return false
	}

	room, err := h.SearchDisponibility(r.nightStart, r.duration)
	if err != nil {
		sendError(r.chanToHandler, err.Error())
		return false
	}

	strRoom     := strconv.FormatUint(uint64(room), 10)
	strNight    := strconv.FormatUint(uint64(r.nightStart), 10)
	strDuration := strconv.FormatUint(uint64(r.duration), 10)

	r.chanToHandler <- "RESULT_FREEROOM " + strRoom + " " + strNight + " " + strDuration
	return true
}

// toString converts disponibilityRequest to string.
func (r *disponibilityRequest) toString() string {

	strNight    := strconv.FormatUint(uint64(r.nightStart), 10)
	strDuration := strconv.FormatUint(uint64(r.duration), 10)
	return "From " + r.clientAddr + " FREEROOM from night" + strNight + " for " + strDuration + " night(s)"
}