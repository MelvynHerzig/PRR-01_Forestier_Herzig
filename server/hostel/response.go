package hostel

// Response is a struct that holds the result from the execution of a Hostelrequestable.
// Success is true if request could be executed without error else false.
// Username is that username used to perform the request. In case of logout with success, it is an empty string
// Message returned by the hostel to be sent to the request provider.
type Response struct {
	Success bool
	Username string
	Message string
}
