package hostel

import "log"

// Channels used with the principle of "Communicating sequential processes". They are user to transfer message between
// Manager and clientHandler/serverHandler (that uses SubmitRequest function).
var (
	// Requests is a channel used to pass hostel requests from SubmitRequest.
	requests = make(chan Request)

	// Responses is a channel used to pass hostel responses to SubmitRequest.
	responses = make (chan Response)
)

// Manager is a function that creates a hostel and then listens to SubmitRequest for incoming Request.
func Manager(roomCount, nightCount uint) {

	hostelManager, hostelError := newHostel(roomCount, nightCount)
	if hostelError != nil {
		log.Fatal(hostelError)
	}

	for {
		select {
			case request := <-requests:
				responses <- request.execute(hostelManager)
		}
	}
}

// SubmitRequest is a function used by clientHandler/serverHandler to submit hostelRequestables.
func SubmitRequest(request Request) Response {
	requests<-request
	return <-responses
}
