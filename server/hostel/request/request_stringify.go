package request

import "strconv"

// ToString converts loginRequest to string.
func (r *loginRequest) ToString() string {
	return "LOGIN with username: " + r.clientName
}

// ToString converts bookRequest to string.
func (r *bookRequest) ToString() string {

	strRoom     := strconv.FormatUint(uint64(r.roomNumber), 10)
	strNight    := strconv.FormatUint(uint64(r.nightStart), 10)
	strDuration := strconv.FormatUint(uint64(r.duration), 10)

	return "BOOK room " + strRoom + " from night" + strNight + " for " + strDuration + " night(s)"
}

// ToString converts roomStateRequest to string.
func (r *roomStateRequest) ToString() string {
	return "ROOMLIST for night" + strconv.FormatUint(uint64(r.nightNumber), 10)
}

// ToString converts disponibilityRequest to string.
func (r *disponibilityRequest) ToString() string {

	strNight    := strconv.FormatUint(uint64(r.nightStart), 10)
	strDuration := strconv.FormatUint(uint64(r.duration), 10)
	return "FREEROOM from night" + strNight + " for " + strDuration + " night(s)"
}

// ToString converts logoutRequest to string.
func (r *logoutRequest) ToString() string {
	return "LOGOUT from: " + r.username
}