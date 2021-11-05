package sync

import (
	"log"
	"net"
	"server/config"
	"server/hostel"
)

// mutexCore this method is executed in a go routine. It is responsible to be the mutex "engine".
// Other thread can use the channels to pass
func mutexCore() {

	// Stores the messages of all servers
	var messages = make([]Message, len(config.Servers))


	// Init the messages
	for i := 0; i < len(messages); i++ {
		messages[i] = Message{
			messageType: REL,
			timestamp: 0,
			srcServer: uint(i),
		}
	}

	hostelManager, hostelError := hostel.NewHostel(nbRooms, nbNights)
	if hostelError != nil {
		log.Fatal(hostelError)
	}

	// Handling clients.
	for {

		if DebugMode {
			debugLogRisk("--------- Enter shared zone ---------")
		}

		select {
		case request := <-requests:

			// TODO call demande (Processus client -> processus mutex)

			// TODO call attente (Processus client -> processus mutex)

			if DebugMode {
				debugLogRequestHandling(request)
			}

			success :=  request.execute(hostelManager, clients)

			if DebugMode {
				debugLogRequestResult(request, success)
			}

			// TODO call fin (Processus client -> processus mutex)

		case cli := <-leaving:
			delete(clients, cli)
			close(cli)
		}

		if DebugMode {
			debugLogRisk("--------- Leave shared zone ---------")
		}
	}
}


func Ask() {

}

func Wait() {

}

func End() {

}
