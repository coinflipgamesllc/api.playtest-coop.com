package domain

// Event is a thing that happened in the application
type Event struct {
	Name string
	Data map[string]interface{}
}
