// Package hostel is the job logic. It implements the hostel room managment.
package hostel

import "errors"

// Hostel This struct if the base struct defining a hostel.
type Hostel struct {
	rooms     [][]uint
	clients   map[string]uint
	nbClients uint
	nbRooms   uint
	nbDays    uint
}

// free If a room value at a given day is 0, the room is not booked.
const free = 0

// clientStartId Client id start at value 1 since 0 means room is not booked.
const clientStartId = 1

// NewHostel Creates a new hostel for a given amount of rooms and days.
func NewHostel(nbRooms, nbDays uint) (*Hostel, error) {
	if nbRooms == 0 || nbDays == 0 {
		return nil, errors.New("number of rooms or number of days cannot be 0")
	}

	hostel := new(Hostel)

	hostel.rooms = make([][]uint, nbRooms)
	for room := range hostel.rooms {
		hostel.rooms[room] = make([]uint, nbDays)
	}

	hostel.clients = make(map[string]uint)
	hostel.nbRooms = nbRooms
	hostel.nbDays = nbDays

	return hostel, nil
}

// TryRegister Try registering a client and gives him an id. If the client already exists for the given hostel
// this is effectless because clients name are supposed unique.
func (h *Hostel) TryRegister(name string) {

	if _, ok := h.clients[name]; ok == false {
		h.clients[name] = clientStartId + h.nbClients
		h.nbClients++
	}
}

// Book Try booking a room for a given day and duration. Client name must be registered.
// Rooms are going from 1 to h.nbRooms. Days are going from 1 to h.nbRooms. Duration cannot be 0.
func (h *Hostel) Book(name string, noRoom, dayStart, duration uint) error {

	// Checks
	if err := h.checkClientRegistered(name); err != nil {
		return err
	}

	if err := h.checkRoomNumber(noRoom); err != nil {
		return err
	}

	if err := h.checkPeriod(dayStart, duration); err != nil {
		return err
	}

	// Room free during booking time ?
	for day := dayStart; day < dayStart+duration; day++ {
		if h.rooms[noRoom-1][day-1] != free {
			return errors.New("room already booked")
		}
	}

	// Booking
	clientId := h.clients[name]
	for day := dayStart; day < dayStart+duration; day++ {
		h.rooms[noRoom-1][day-1] = clientId
	}

	return nil
}

// GetRoomsState Returns state for each rooms: free, self reserved or occupied. Client must be registered.
// Days are going from 1 to h.nbRooms.
func (h *Hostel) GetRoomsState(name string, noDay uint) ([]string, error) {

	roomsState := make([]string, h.nbRooms)

	// Checks
	if err := h.checkClientRegistered(name); err != nil {
		return roomsState, err
	}

	if err := h.checkDayNumber(noDay); err != nil {
		return roomsState, err
	}

	clientId := h.clients[name]

	// Filling room state slice
	for room := uint(0); room < h.nbRooms; room++ {

		switch h.rooms[room][noDay - 1] {
		case 0:
			roomsState[room] = "Free"
		case clientId:
			roomsState[room] = "Self reserved"
		default:
			roomsState[room] = "Occupied"
		}
	}

	return roomsState, nil
}

// SearchDisponibility Looks for a free room starting from a given day during a given duration.
// Days are going from 1 to h.nbRooms. Duration cannot be 0.
func (h *Hostel) SearchDisponibility(dayStart, duration uint) (uint, error) {

	if err := h.checkPeriod(dayStart, duration); err != nil {
		return 0, err
	}

	for room := uint(0); room < h.nbRooms; room++ {

		free := true

		for day := dayStart; day < dayStart+duration; day++ {
			if h.rooms[room][day-1] != 0 {
				free = false
			}
		}

		if free == true {
			return room + 1, nil
		}
	}

	return 0, errors.New("no free room for this period")
}

// checkClientRegistered Checks if the client is registered in hostel clients map. Return nil if this is the case.
func (h *Hostel) checkClientRegistered(name string) error {

	if _, ok := h.clients[name]; ok != true {
		return errors.New("unknown client name. Please register first")
	}

	return nil
}

// checkRoomNumber Checks if the room number is between 1 and hostel number of rooms. Returns nil if this is the case.
func (h *Hostel) checkRoomNumber(noRoom uint) error {

	if noRoom == 0 || noRoom > h.nbRooms {
		return errors.New("invalid room number")
	}

	return nil
}

// checkDayNumber Checks if the day number is between 1 and hostel max day plan. Returns nil if this is the case.
func (h *Hostel) checkDayNumber(noDay uint) error {

	if noDay == 0 || noDay > h.nbDays {
		return errors.New("invalid day number")
	}

	return nil
}

// checkPeriod Checks if the duration starting from the given day is not going further the hostel max day plan.
func (h *Hostel) checkPeriod(noDay, duration uint) error {

	if err := h.checkDayNumber(noDay); err != nil {
		return err
	}

	if duration == 0 || noDay+duration-1 > h.nbDays {
		return errors.New("invalid duration")
	}

	return nil
}
