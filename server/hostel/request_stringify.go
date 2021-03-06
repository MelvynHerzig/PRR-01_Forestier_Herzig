// Package hostel implements Request function to stringify Request (human-readable)
package hostel

import "fmt"

// ToString converts loginRequest to string for human reading.
func (r *loginRequest) ToString() string {
	return "LOGIN with username: " + r.clientName
}

// ToString converts bookRequest to string for human reading.
func (r *bookRequest) ToString() string {
	return fmt.Sprintf("BOOK room %d from night %d for %d night(s)", r.roomNumber, r.nightStart, r.duration)
}

// ToString converts roomStateRequest to string for human reading.
func (r *roomStateRequest) ToString() string {
	return fmt.Sprintf("ROOMLIST for night %d", r.nightNumber)
}

// ToString converts disponibilityRequest to string for human reading.
func (r *disponibilityRequest) ToString() string {
	return fmt.Sprintf("FREEROOM from night %d for %d night(s)", r.nightStart, r.duration)
}

// ToString converts logoutRequest to string for human reading.
func (r *logoutRequest) ToString() string {
	return "LOGOUT from: " + r.username
}