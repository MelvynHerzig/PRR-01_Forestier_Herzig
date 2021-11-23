package hostel

// Request interface of communication that can be submitted to Manager with SubmitRequest.
type Request interface {
	execute(h *hostel) Response
	ShouldReplicate() bool
	ToString() string
	SetUsername(name string)
	Serialize() string
}

// hostelRequest base content of all kind of requestables.
type hostelRequest struct {
	username string // Logged username to use, can be null in the case of a loginRequest
}

// SetUsername sets the username used to execute the request.
func (r *hostelRequest) SetUsername(name string) {
	r.username = name
}

// loginRequest request that logs clients in, in order to book, inspect and search free room.
// Client can't login with the associated clientName if the username is already logged in
type loginRequest struct {
	hostelRequest
	clientName string
}

// logoutRequest request that logs client out, in order to change "account".
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

// disponibilityRequest request used to find a free room for a given night and duration.
type disponibilityRequest struct {
	hostelRequest
	nightStart uint
	duration   uint
}
