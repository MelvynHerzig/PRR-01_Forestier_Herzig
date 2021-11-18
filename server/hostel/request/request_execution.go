package request

import (
	"server/hostel"
)

// execute Calls hostel login function with provided clientName in request.
func (r *loginRequest) execute(h *hostel.Hostel) string {

	return h.Login(r.clientName)
}

// execute Calls hostel book function with provided period in request.
func (r *bookRequest) execute(h *hostel.Hostel) string{

	return h.Book(r.username, r.roomNumber, r.nightStart, r.duration)
}

// execute Calls hostel getRoomState function with provided arguments in request.
func (r *roomStateRequest) execute(h *hostel.Hostel) string {

	return h.GetRoomsState(r.username, r.nightNumber)
}

// execute Call hostel searchDisponibility function with provided arguments in request.
func (r *disponibilityRequest) execute(h *hostel.Hostel) string {

	return h.SearchDisponibility(r.username, r.nightStart, r.duration)
}

// execute Calls hostel logout function with request username in request.
func (r *logoutRequest) execute(h *hostel.Hostel) string {

	return h.Logout(r.username)
}
