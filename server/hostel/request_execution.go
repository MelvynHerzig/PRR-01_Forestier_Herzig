// Package hostel implements Request execution function.
package hostel

import "strings"

// execute calls hostel.login function with provided clientName.
// In case of a successful login, username return equals clientName.
func (r *loginRequest) execute(h *hostel) Response {

	// If not null, client already logged in.
	if r.username != "" {
		return Response{false, "", "Must log out before log in."}
	}

	// execution
	response := h.login(r.clientName)
	return Response{strings.Split(response, " ")[0] == "RESULT_LOGIN", r.clientName, response }
}

// execute calls hostel.book function with provided room and period.
func (r *bookRequest) execute(h *hostel) Response {

	// execution
	response := h.book(r.username, r.roomNumber, r.nightStart, r.duration)
	return Response{strings.Split(response, " ")[0] == "RESULT_BOOK", r.username, response}
}

// execute calls hostel.getRoomState function with provided night number.
func (r *roomStateRequest) execute(h *hostel) Response {

	// execution
	response := h.getRoomsState(r.username, r.nightNumber)
	return Response{strings.Split(response, " ")[0] == "RESULT_ROOMLIST", r.username, response}
}

// execute calls hostel.searchDisponibility function with provided period.
func (r *disponibilityRequest) execute(h *hostel) Response {

	//execution
	response := h.searchDisponibility(r.username, r.nightStart, r.duration)
	return Response{strings.Split(response, " ")[0] == "RESULT_FREEROOM", r.username, response}
}

// execute calls hostel logout function with request username provided.
// In case of successful logout, returned username is empty string.
func (r *logoutRequest) execute(h *hostel) Response {

	// execution
	response := h.logout(r.username)
	return Response{strings.Split(response, " ")[0] == "RESULT_LOGOUT", "", response}
}
