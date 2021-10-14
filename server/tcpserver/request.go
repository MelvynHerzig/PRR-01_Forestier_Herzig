// Package tcpserver implements the Server TCP logic.
// The main.go file is responsible for TCP "worker" and
// the request.go file is responsible for request that can
// be submitted to the job logic layer.
package tcpserver

import (
	"Server/hostel"
	"fmt"
	"strconv"
	"time"
)

// DebugMode enables debug mode if set to true.
// Debug mode consist of logging server actions and sleeping 20s when 2 clients logs in, in order to
// let enough time to create a race condition.
var DebugMode = false

// Number of logged client. It is race safe since incremented and decremented in loginRequest.execute
//and logoutRequest.execute. All this, because the execute methodes are handled in hostelManager function
// that respects Communicating Sequential Processes.
var nbLoggedClient = 0;

// hostelRequestable interface of request that can be made to hostel.
type hostelRequestable interface {
	execute(h *hostel.Hostel, clients map[client]string)
}

// hostelRequest base content of all kind of request.
type hostelRequest struct {
	chanToHandler client
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
func (r *loginRequest) execute(h *hostel.Hostel, clients map[client]string) {

	if clients[r.chanToHandler] != "" {
		sendError(r.chanToHandler, "You are already connected as " + clients[r.chanToHandler] + ".")
		return
	}

	clients[r.chanToHandler] = r.clientName
	h.TryRegister(r.clientName)

	if DebugMode {
		debugModeLog(r.clientName + " logged in.")
		nbLoggedClient++

		if nbLoggedClient == 2 {
			debugModeLog("Server request handler suspended. Resume in 20s.")
			time.Sleep(20 * time.Second)
			debugModeLog("Server request handler resumed.")
		}
	}

	r.chanToHandler <- "RESULT_LOGIN"
}

// execute if user is logged in, log him out, otherwise explains error. All is notified to clients channel.
func (r *logoutRequest) execute(h *hostel.Hostel, clients map[client]string) {

	if clients[r.chanToHandler] == "" {
		sendError(r.chanToHandler, "You are not logged in.")
		return
	}

	oldName := clients[r.chanToHandler]
	clients[r.chanToHandler] = ""

	if DebugMode {
		nbLoggedClient--
		debugModeLog(oldName + " logged out.")
	}

	r.chanToHandler <- "RESULT_LOGOUT"
}

// execute tries to book a room. Client must be logged in. All is notified to clients channel.
func (r *bookRequest) execute(h *hostel.Hostel, clients map[client]string) {

	if clients[r.chanToHandler] == "" {
		sendError(r.chanToHandler, "You are not logged in.")
		return
	}

	if err := h.Book(clients[r.chanToHandler], r.roomNumber, r.nightStart, r.duration); err != nil {
		sendError(r.chanToHandler,  err.Error())
		return
	}

	strRoom     := strconv.FormatUint(uint64(r.roomNumber), 10)
	strNight    := strconv.FormatUint(uint64(r.nightStart), 10)
	strDuration := strconv.FormatUint(uint64(r.duration), 10)

	if DebugMode {
		debugModeLog( clients[r.chanToHandler] + " booked room. Room " + strRoom + " from night " + strNight + " for " + strDuration + " night(s)")
	}

	r.chanToHandler <- "RESULT_BOOK " + strRoom + " " + strNight + " " + strDuration
}

// execute tries to get room state for a given night. Client must be logged in. All is notified to clients channel.
func (r *roomStateRequest) execute(h *hostel.Hostel, clients map[client]string) {

	if clients[r.chanToHandler] == "" {
		sendError(r.chanToHandler, "You are not logged in.")
		return
	}

	states, err := h.GetRoomsState(clients[r.chanToHandler], r.nightNumber)
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

	if DebugMode {
		debugModeLog( clients[r.chanToHandler] + " is consulting rooms for night " + strconv.FormatUint(uint64(r.nightNumber), 10))
	}

	r.chanToHandler <- "RESULT_ROOMLIST " + res
}

// execute tries to find a free room. Client must be logged in. All is notified to clients channel.
func (r *disponibilityRequest) execute(h *hostel.Hostel, clients map[client]string) {

	if clients[r.chanToHandler] == "" {
		sendError(r.chanToHandler, "You are not logged in.")
		return
	}

	room, err := h.SearchDisponibility(r.nightStart, r.duration)
	if err != nil {
		sendError(r.chanToHandler, err.Error())
		return
	}

	strRoom     := strconv.FormatUint(uint64(room), 10)
	strNight    := strconv.FormatUint(uint64(r.nightStart), 10)
	strDuration := strconv.FormatUint(uint64(r.duration), 10)

	if DebugMode {
		debugModeLog( clients[r.chanToHandler] + " found a free room. Room " + strRoom + " from night " + strNight + " for " + strDuration + " night(s)")
	}

	r.chanToHandler <- "RESULT_FREEROOM " + strRoom + " " + strNight + " " + strDuration
}

func debugModeLog(message string) {
	fmt.Println("DEBUG >> ", message)
}