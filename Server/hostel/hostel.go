package hostel

import "errors"

type Hostel struct {
	rooms     [][]uint
	clients   map[string]uint
	nbClients uint
	nbRooms   uint
	nbDays    uint
}

// Room states
const (
	Free = iota
	SelfReserved
	Occupied
)

var RoomsStateSignification = map[uint]string{
	Free:         "free",
	SelfReserved: "Self reserved",
	Occupied:     "Occupied",
}

const clientStartId = 1

// Creates a new hosteé
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

// true: register successfully, false: already exists
func (h *Hostel) RegisterClient(name string) bool {

	if err := h.checkClient(name); err == nil {
		return false
	}

	h.clients[name] = clientStartId + h.nbClients
	h.nbClients++

	return true
}

// true, nil if booked successfully else false, err if book failed
func (h *Hostel) Book(name string, noRoom, dayStart, duration uint) (bool, error) {

	// Checks
	if err := h.checkClient(name); err != nil {
		return false, err
	}

	if err := h.checkRoomNumber(noRoom); err != nil {
		return false, err
	}

	if err := h.checkPeriod(dayStart, duration); err != nil {
		return false, err
	}

	// Room free during booking time ?
	for day := dayStart; day < dayStart+duration; day++ {
		if h.rooms[noRoom-1][day-1] != Free {
			return false, errors.New("room already booked")
		}
	}

	// Booking
	clientId := h.clients[name]
	for day := dayStart; day < dayStart+duration; day++ {
		h.rooms[noRoom-1][day-1] = clientId
	}

	return true, nil
}

// Returns an error if args are invalid else nil with FREE, OCCUPIED or SELF_RESERVED slice
func (h *Hostel) GetRoomsState(name string, noDay uint) ([]uint, error) {

	roomsState := make([]uint, h.nbRooms)

	// Checks
	if err := h.checkClient(name); err != nil {
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
			roomsState[room] = Free
		case clientId:
			roomsState[room] = SelfReserved
		default:
			roomsState[room] = Occupied
		}
	}

	return roomsState, nil
}

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

// Return nil if client is registered else not nil
func (h *Hostel) checkClient(name string) error {

	if _, ok := h.clients[name]; ok != true {
		return errors.New("unknown client")
	}

	return nil
}

// Return nil if room number is valid [1, h.nbRooms] else not nil
func (h *Hostel) checkRoomNumber(noRoom uint) error {

	if noRoom == 0 || noRoom > h.nbRooms {
		return errors.New("invalid room number")
	}

	return nil
}

// Return nil if day number is valid [1, h.nbDays] else not nil
func (h *Hostel) checkDayNumber(noDay uint) error {

	if noDay == 0 || noDay > h.nbDays {
		return errors.New("invalid day number")
	}

	return nil
}

// Return nil if period is valid [1, h.nbDays] else not nil
func (h *Hostel) checkPeriod(noDay, duration uint) error {

	if err := h.checkDayNumber(noDay); err != nil {
		return err
	}

	if duration == 0 || noDay+duration-1 > h.nbDays {
		return errors.New("invalid duration")
	}

	return nil
}
