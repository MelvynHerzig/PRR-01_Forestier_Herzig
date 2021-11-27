// Package hostel a hostel with custom function to manipulate the hostel.
package hostel

import (
	"errors"
	"fmt"
)

// hostel is the base struct that defines a hostel.
type hostel struct {
	rooms     [][]uint               // 0 freeRoom else client id
	clients   map[string] *client    // map of client name to their id and online status
	nbClients uint					 // count of hostel registered clients, in other words len of clients
	nbRooms   uint					 // number of rooms
	nbNights  uint                   // number of nights supported by the planning
}

// client is the struct that is mapped to a client name in hostel.clients. It consists of a unique identifier and an
// online status.
type client struct {
	id uint        // id from [0, max uint]
	logged bool
}

// freeRoom is the constant value stored in rooms to declare that the room is free.
const freeRoom = 0

// newHostel creates a new hostel for a given amount of rooms and nights.
func newHostel(nbRooms, nbNights uint) (*hostel, error) {
	if nbRooms == 0 || nbNights == 0 {
		return nil, errors.New("number of rooms or number of nights cannot be 0")
	}

	hostel := new(hostel)

	hostel.rooms = make([][]uint, nbRooms)
	for room := range hostel.rooms {
		hostel.rooms[room] = make([]uint, nbNights)
	}

	hostel.clients = make(map[string]*client)
	hostel.nbRooms = nbRooms
	hostel.nbNights = nbNights

	return hostel, nil
}

// Login logs a client in. If he doesn't already exist, the client is registered. Fails if client already logged.
func (h* hostel) login(name string) string {

	// Try to register client if he doesn't already exist
	if _, ok := h.clients[name]; name != "" && ok == false {
		h.clients[name] = &client{id: 1 + h.nbClients, logged: false}
		h.nbClients++
	}

	if c := h.clients[name]; c.logged == false {
		h.clients[name].logged = true

		return "RESULT_LOGIN"
	} else {
		return "Client already logged in."
	}
}

// book try to book a room for a given night and duration. Client name must be registered.
// Rooms are going from 1 to h.nbRooms. Nights are going from 1 to h.nbRooms. Duration cannot be 0.
// Fails if username is not logged, noRoom is an invalid number or period (nightStart and duration) not valid.
func (h *hostel) book(username string, noRoom, nightStart, duration uint) string {

	// Checks
	if ok, msg := h.checkClientLogged(username); ok == false {
		return msg
	}

	if ok, msg := h.checkRoomNumber(noRoom); ok == false {
		return msg
	}

	if ok, msg := h.checkPeriod(nightStart, duration); ok == false {
		return msg
	}

	// Room free during booking time ?
	for night := nightStart; night < nightStart+duration; night++ {
		if h.rooms[noRoom-1][night-1] != freeRoom {
			return "Room already booked"
		}
	}

	// Booking
	clientId := h.clients[username].id
	for night := nightStart; night < nightStart+duration; night++ {
		h.rooms[noRoom-1][night-1] = clientId
	}

	return fmt.Sprintf("RESULT_BOOK %d %d %d", noRoom, nightStart, duration)
}

// getRoomsState returns state for each room: "free", "self reserved" or "occupied". Client must be registered.
// Nights are going from 1 to h.nbRooms.
// Fails if username is not logged in or noNight in bounds of hostel planning.
func (h *hostel) getRoomsState(username string, noNight uint) string {
	// Checks
	if ok, msg := h.checkClientLogged(username); ok == false {
		return msg
	}

	if ok, msg := h.checkNightNumber(noNight); ok == false {
		return msg
	}

	// Filling room state response
	res := ""
	for room := uint(0); room < h.nbRooms; room++ {

		if room != 0 {
			res += ","
		}

		switch h.rooms[room][noNight- 1] {
		case freeRoom:
			res += "Free"
		case h.clients[username].id:
			res += "Self reserved"
		default:
			res += "Occupied"
		}
	}

	return "RESULT_ROOMLIST " + res
}

// searchDisponibility looks for a free room starting from a given night during a given duration.
// Nights are going from 1 to h.nbRooms. Duration cannot be 0.
// Fails if username is not logged or if the period (nightStart + duration) is not in hostel planning bounds.
func (h *hostel) searchDisponibility(username string, nightStart, duration uint) string {

	// Checks
	if ok, msg := h.checkClientLogged(username); ok == false {
		return msg
	}

	if ok, msg := h.checkPeriod(nightStart, duration); ok == false {
		return msg
	}

	for room := uint(0); room < h.nbRooms; room++ {

		free := true

		for night := nightStart; night < nightStart + duration; night++ {
			if h.rooms[room][night-1] != freeRoom {
				free = false
			}
		}

		if free == true {
			return fmt.Sprintf("RESULT_FREEROOM %d %d %d", room + 1, nightStart, duration)
		}
	}

	return "RESULT_FREEROOM 0"
}

// logout logs a client out.
// Fails if client is not logged in
func (h* hostel) logout(name string) string {

	// Checks
	if ok, msg := h.checkClientLogged(name); ok == false {
		return msg
	}

	h.clients[name].logged = false

	return "RESULT_LOGOUT"
}

// checkClientLogged checks if the client is registered in hostel clients map. Returns nil if this is the case.
func (h *hostel) checkClientLogged(name string) (bool, string) {

	if client, registered := h.clients[name]; registered == false || client.logged == false{
		return false, "Client not logged"
	}

	return true, ""
}

// checkRoomNumber checks if the room number is between 1 and hostel number of rooms. Returns nil if this is the case.
func (h *hostel) checkRoomNumber(noRoom uint) (bool, string) {

	if noRoom == 0 || noRoom > h.nbRooms {
		return false, "Invalid room number"
	}

	return true, ""
}

// checkNightNumber checks if the night number is between 1 and hostel max night plan. Returns nil if this is the case.
func (h *hostel) checkNightNumber(noNight uint) (bool, string) {

	if noNight == 0 || noNight > h.nbNights {
		return false, "Invalid night number"
	}

	return true, ""
}

// checkPeriod checks if the duration starting from the given night is not going further than the hostel max night plan.
func (h *hostel) checkPeriod(noNight, duration uint) (bool, string) {

	if result, message := h.checkNightNumber(noNight); result == false {
		return result, message
	}

	if duration == 0 || noNight+duration-1 > h.nbNights {
		return false, "Invalid duration"
	}

	return true, ""
}
