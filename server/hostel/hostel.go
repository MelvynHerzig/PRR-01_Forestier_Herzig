// Package hostel is the job logic layer. It implements the hostel room managment.
package hostel

import "errors"

// Hostel This struct if the base struct defining a hostel.
type Hostel struct {
	rooms     [][]uint        // 0 freeRoom else client id
	clients   map[string]uint // name -> id from [0, max uint]
	nbClients uint
	nbRooms   uint
	nbNights  uint
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

	hostel.clients = make(map[string]uint)
	hostel.nbRooms = nbRooms
	hostel.nbNights = nbNights

	return hostel, nil
}

// TryRegister try to register a client and gives him an id. If the client already exists for the given hostel
// this is effectless because clients name are supposed unique.
func (h *Hostel) TryRegister(name string) {

	if _, ok := h.clients[name]; ok == false {
		h.clients[name] = 1 + h.nbClients // 1 + because client's id start from 1
		h.nbClients++
	}
}

// Book try to book a room for a given night and duration. Client name must be registered.
// Rooms are going from 1 to h.nbRooms. Nights are going from 1 to h.nbRooms. Duration cannot be 0.
func (h *Hostel) Book(name string, noRoom, nightStart, duration uint) error {

	// Checks
	if err := h.checkClientRegistered(name); err != nil {
		return err
	}

	if err := h.checkRoomNumber(noRoom); err != nil {
		return err
	}

	if err := h.checkPeriod(nightStart, duration); err != nil {
		return err
	}

	// Room free during booking time ?
	for night := nightStart; night < nightStart+duration; night++ {
		if h.rooms[noRoom-1][night-1] != freeRoom {
			return errors.New("room already booked")
		}
	}

	// Booking
	clientId := h.clients[name]
	for night := nightStart; night < nightStart+duration; night++ {
		h.rooms[noRoom-1][night-1] = clientId
	}

	return nil
}

// GetRoomsState returns state for each rooms: "free", "self reserved" or "occupied". Client must be registered.
// Nights are going from 1 to h.nbRooms.
func (h *Hostel) GetRoomsState(name string, noNight uint) ([]string, error) {

	roomsState := make([]string, h.nbRooms)

	// Checks
	if err := h.checkClientRegistered(name); err != nil {
		return roomsState, err
	}

	if err := h.checkNightNumber(noNight); err != nil {
		return roomsState, err
	}

	clientId := h.clients[name]

	// Filling room state slice
	for room := uint(0); room < h.nbRooms; room++ {

		switch h.rooms[room][noNight- 1] {
		case freeRoom:
			roomsState[room] = "Free"
		case clientId:
			roomsState[room] = "Self reserved"
		default:
			roomsState[room] = "Occupied"
		}
	}

	return roomsState, nil
}

// SearchDisponibility looks for a free room starting from a given night during a given duration.
// Nights are going from 1 to h.nbRooms. Duration cannot be 0.
func (h *Hostel) SearchDisponibility(nightStart, duration uint) (uint, error) {

	if err := h.checkPeriod(nightStart, duration); err != nil {
		return 0, err
	}

	for room := uint(0); room < h.nbRooms; room++ {

		free := true

		for night := nightStart; night < nightStart+duration; night++ {
			if h.rooms[room][night-1] != freeRoom {
				free = false
			}
		}

		if free == true {
			return room + 1, nil
		}
	}

	return 0, nil
}

// checkClientRegistered checks if the client is registered in hostel clients map. Returns nil if this is the case.
func (h *Hostel) checkClientRegistered(name string) error {

	if _, ok := h.clients[name]; ok != true {
		return errors.New("unknown client name. Please register first")
	}

	return nil
}

// checkRoomNumber checks if the room number is between 1 and hostel number of rooms. Returns nil if this is the case.
func (h *Hostel) checkRoomNumber(noRoom uint) error {

	if noRoom == 0 || noRoom > h.nbRooms {
		return errors.New("invalid room number")
	}

	return nil
}

// checkNightNumber checks if the night number is between 1 and hostel max night plan. Returns nil if this is the case.
func (h *Hostel) checkNightNumber(noNight uint) error {

	if noNight == 0 || noNight > h.nbNights {
		return errors.New("invalid night number")
	}

	return nil
}

// checkPeriod checks if the duration starting from the given night is not going further than the hostel max night plan.
func (h *Hostel) checkPeriod(noNight, duration uint) error {

	if err := h.checkNightNumber(noNight); err != nil {
		return err
	}

	if duration == 0 || noNight+duration-1 > h.nbNights {
		return errors.New("invalid duration")
	}

	return nil
}
