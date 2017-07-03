package message

import "time"

// StuffHappenedEvent is what it is, yo
type StuffHappenedEvent struct {
	Stuff string
	When  time.Time
}
