package request

import (
	"server/hostel"
	"strings"
)

// Execute Calls hostel login function with provided clientName in request.
func (r *loginRequest) Execute(h *hostel.Hostel) (bool, string, string) {

	if r.username != "" {
		return false, "", "Must log out before log in."
	}

	response := h.Login(r.clientName)
	return strings.Split(response, " ")[0] == "RESULT_LOGIN", r.clientName, response
}

// Execute Calls hostel book function with provided period in request.
func (r *bookRequest) Execute(h *hostel.Hostel) (bool, string, string){

	response := h.Book(r.username, r.roomNumber, r.nightStart, r.duration)
	return strings.Split(response, " ")[0] == "RESULT_BOOK", r.username, response
}

// Execute Calls hostel getRoomState function with provided arguments in request.
func (r *roomStateRequest) Execute(h *hostel.Hostel) (bool, string, string) {

	response := h.GetRoomsState(r.username, r.nightNumber)
	return strings.Split(response, " ")[0] == "RESULT_ROOMLIST", r.username, response
}

// Execute Call hostel searchDisponibility function with provided arguments in request.
func (r *disponibilityRequest) Execute(h *hostel.Hostel) (bool, string, string) {

	response := h.SearchDisponibility(r.username, r.nightStart, r.duration)
	return strings.Split(response, " ")[0] == "RESULT_FREEROOM", r.username, response
}

// Execute Calls hostel logout function with request username in request.
func (r *logoutRequest) Execute(h *hostel.Hostel) (bool, string, string) {

	response := h.Logout(r.username)
	return strings.Split(response, " ")[0] == "RESULT_LOGOUT", "", response
}
