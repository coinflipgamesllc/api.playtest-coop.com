package game

import (
	"fmt"
)

// Status tracks where a game is in the lifecycle
type Status string

const (
	// Prototype games are actively being worked on by designers
	Prototype Status = "Prototype"

	// Signed games are under contract by a publisher, and sometimes in development, but not available for purchase
	Signed = "Signed"

	// Published games are available for purchase
	Published = "Published"

	// Archived games are any that are no longer publicly visible
	Archived = "Archived"
)

// StatusFromString returns the Status type corresponding to the provided string
func StatusFromString(s string) (Status, error) {
	switch s {
	case "Prototype":
		return Prototype, nil
	case "Signed":
		return Signed, nil
	case "Published":
		return Published, nil
	case "Archived":
		return Archived, nil
	default:
		return "", InvalidStatus{s}
	}
}

// InvalidStatus returned for strings that don't match a status we're tracking
type InvalidStatus struct {
	PassedValue string
}

func (e InvalidStatus) Error() string {
	return fmt.Sprintf("invalid status '%s'", e.PassedValue)
}
