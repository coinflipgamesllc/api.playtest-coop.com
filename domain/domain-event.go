package domain

// DomainEvent is a thing that happened in the application
type DomainEvent struct {
	Name string
	Data map[string]interface{}
}
