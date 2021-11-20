// Package clock implements lamport clock system
package clock

// timestamp is the local time
var timestamp uint = 0

// GetTimestamp returns the local timestamp
func GetTimestamp() uint {
	return timestamp
}

// SyncTimestamp adjusts the local timestamp with receivedTimestamp
func SyncTimestamp(receivedTimestamp uint) {
	timestamp = maxTimeStamp(timestamp, receivedTimestamp) + 1
}

// IncTimestamp increments timestamp
func IncTimestamp() {
	timestamp += 1
}

// maxTimeStamp returns the biggest timestamp between timestamp1 and timestamp2
func maxTimeStamp(timestamp1 uint, timestamp2 uint) uint{
	if timestamp1 >= timestamp2 {
		return timestamp1
	}

	return timestamp2
}