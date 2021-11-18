package request

import "server/hostel"

// HostelRequestable interface of request that can be made to hostel.
type HostelRequestable interface {
	execute(h *hostel.Hostel) string
	ToString() string
	setUsername(name string)
	serialize() string
}

// hostelRequest base content of all kind of request.
type hostelRequest struct {
	username string
}

// setUsername sets the username used to execute the request.
func (r *hostelRequest) setUsername(name string) {
	r.username = name
}

// loginRequest request that logs clients in, in order to book, inspect and search free room.
// Client cant login if the username is already logged in
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
