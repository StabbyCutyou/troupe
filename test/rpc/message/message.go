package message

import "time"

// StuffHappenedEvent is what it is, yo
type StuffHappenedEvent struct {
	Stuff string
	When  time.Time
}

// StuffHappenedEventError is an error thrown from inside of the SHE handler
type StuffHappenedEventError string

// Error implements the interface
func (e StuffHappenedEventError) Error() string {
	return string(e)
}
