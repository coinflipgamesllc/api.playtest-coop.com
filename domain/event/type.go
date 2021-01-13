package event

import "fmt"

// Type tracks what kind of event is happening, which changes how the event will be run
type Type string

const (
	// Remote free-for-all event
	Remote Type = "Remote"

	// InPerson free-for-all event
	InPerson = "In-person"

	// Not implemented yet
	// RemoteRegistered is a remote event requiring registration and timeslot/table signup
	// RemoteRegistered = "Remote Registered"

	// InPersonRegistered is an in-person event requiring registration and timeslot/table signup
	// InPersonRegistered = "In-person Registered"
)

// TypeFromString returns the Type corresponding to the provided string
func TypeFromString(s string) (Type, error) {
	switch s {
	case "Remote":
		return Remote, nil
	case "InPerson":
		return InPerson, nil
	default:
		return "", InvalidType{s}
	}
}

// InvalidType returned for strings that don't match a type we're tracking
type InvalidType struct {
	PassedValue string
}

func (e InvalidType) Error() string {
	return fmt.Sprintf("invalid type '%s'", e.PassedValue)
}
