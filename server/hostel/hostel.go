// Package hostel is the job logic layer. It implements the hostel room managment.
package hostel

import (
	"errors"
	"strconv"
)

// Hostel This struct if the base struct defining a hostel.
type Hostel struct {
	rooms     [][]uint        // 0 freeRoom else client id
	clients   map[string] *client
	nbClients uint
	nbRooms   uint
	nbNights  uint
}

// Properties of a client, name is already stores in clients map.
type client struct {
	id uint        // id from [0, max uint]
	logged bool
}

// freeRoom is the constant value stored in rooms to declare that the room is free.
const freeRoom = 0

// NewHostel creates a new hostel for a given amount of rooms and nights.
func NewHostel(nbRooms, nbNights uint) (*Hostel, error) {
	if nbRooms == 0 || nbNights == 0 {
		return nil, errors.New("number of rooms or number of nights cannot be 0")
	}

	hostel := new(Hostel)

	hostel.rooms = make([][]uint, nbRooms)
	for room := range hostel.rooms {
		hostel.rooms[room] = make([]uint, nbNights)
	}

	hostel.clients = make(map[string]*client)
	hostel.nbRooms = nbRooms
	hostel.nbNights = nbNights

	return hostel, nil
}

// Login logs a client in. If he doesn't already exist, the client is registered.
func (h* Hostel) Login(name string) string {

	// Try to register client if he doesn't already exist
	if _, ok := h.clients[name]; ok == false {
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

// Book try to book a room for a given night and duration. Client name must be registered.
// Rooms are going from 1 to h.nbRooms. Nights are going from 1 to h.nbRooms. Duration cannot be 0.
func (h *Hostel) Book(username string, noRoom, nightStart, duration uint) string {

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

	strRoom     := strconv.FormatUint(uint64(noRoom), 10)
	strNight    := strconv.FormatUint(uint64(nightStart), 10)
	strDuration := strconv.FormatUint(uint64(duration ), 10)

	return "RESULT_BOOK " + strRoom + " " + strNight + " " + strDuration
}

// GetRoomsState returns state for each rooms: "free", "self reserved" or "occupied". Client must be registered.
// Nights are going from 1 to h.nbRooms.
func (h *Hostel) GetRoomsState(username string, noNight uint) string {

	roomsState := make([]string, h.nbRooms)

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
			roomsState[room] = "Free"
		case h.clients[username].id:
			roomsState[room] = "Self reserved"
		default:
			roomsState[room] = "Occupied"
		}
	}

	return "RESULT_ROOMLIST " + res
}

// SearchDisponibility looks for a free room starting from a given night during a given duration.
// Nights are going from 1 to h.nbRooms. Duration cannot be 0.
func (h *Hostel) SearchDisponibility(username string, nightStart, duration uint) string {

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
			strRoom     := strconv.FormatUint(uint64(room + 1), 10)
			strNight    := strconv.FormatUint(uint64(nightStart), 10)
			strDuration := strconv.FormatUint(uint64(duration ), 10)
			return "RESULT_FREEROOM " + strRoom + " " + strNight + " " + strDuration
		}
	}

	return "No free room found for this period"
}

// Logout logs a client out.
func (h* Hostel) Logout(name string) string {

	// Checks
	if ok, msg := h.checkClientLogged(name); ok == false {
		return msg
	}

	return "RESULT_LOGOUT"
}

// checkClientLogged checks if the client is registered in hostel clients map. Returns nil if this is the case.
func (h *Hostel) checkClientLogged(name string) (bool, string) {

	if client, registered := h.clients[name]; registered == false || client.logged == false{
		return false, "Client not logged"
	}

	return true, ""
}

// checkRoomNumber checks if the room number is between 1 and hostel number of rooms. Returns nil if this is the case.
func (h *Hostel) checkRoomNumber(noRoom uint) (bool, string) {

	if noRoom == 0 || noRoom > h.nbRooms {
		return false, "Invalid room number"
	}

	return true, ""
}

// checkNightNumber checks if the night number is between 1 and hostel max night plan. Returns nil if this is the case.
func (h *Hostel) checkNightNumber(noNight uint) (bool, string) {

	if noNight == 0 || noNight > h.nbNights {
		return false, "Invalid night number"
	}

	return true, ""
}

// checkPeriod checks if the duration starting from the given night is not going further than the hostel max night plan.
func (h *Hostel) checkPeriod(noNight, duration uint) (bool, string) {

	if result, message := h.checkNightNumber(noNight); result == false {
		return result, message
	}

	if duration == 0 || noNight+duration-1 > h.nbNights {
		return false, "Invalid duration"
	}

	return true, ""
}
