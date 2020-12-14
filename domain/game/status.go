package game

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
