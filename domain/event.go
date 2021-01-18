package domain

import (
	"time"

	"github.com/coinflipgamesllc/api.playtest-coop.com/domain/event"
	"gorm.io/gorm"
)

// Event holds metadata for an organized playtesting event. The actual scheduling data is
// stored in ical format.
type Event struct {
	ID        uint           `json:"id" gorm:"primarykey" example:"123"`
	CreatedAt time.Time      `json:"created_at" example:"2020-12-11T15:29:49.321629-08:00"`
	UpdatedAt time.Time      `json:"updated_at" example:"2020-12-13T15:42:40.578904-08:00"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	Title        string     `json:"title" example:"Seattle Wednesday Night Playtesting"`
	Type         event.Type `json:"type"`
	Facilitators []User     `json:"facilitators" gorm:"many2many:event_facilitators;"`
	Details      string     `json:"details" example:"Get together and test out some games!"`
	Location     string     `json:"location,omitempty" example:"123 Fake St..."`
	URL          string     `json:"url,omitempty" example:"https://discord.gg/ABC1234"`

	Duration time.Duration `json:"duration" example:"14400000"`
	RRule    string        `json:"rrule"`
}

// EventRepository defines how to interact with events in database
type EventRepository interface {
	ListEvents() ([]Event, error)
	EventOfID(id uint) (*Event, error)
	Save(*Event) error
}

// NewRemoteEvent creates a remote playtesting event
func NewRemoteEvent(title, details, url string, duration int64, rrule string, primaryFacilitator User) *Event {
	return &Event{
		Type:         event.Remote,
		Title:        title,
		Details:      details,
		Facilitators: []User{primaryFacilitator},
		URL:          url,
		Duration:     time.Duration(duration),
		RRule:        rrule,
	}
}

// NewInPersonEvent creates an in-person playtesting event
func NewInPersonEvent(title, details, location string, duration int64, rrule string, primaryFacilitator User) *Event {
	return &Event{
		Type:         event.InPerson,
		Title:        title,
		Details:      details,
		Facilitators: []User{primaryFacilitator},
		Location:     location,
		Duration:     time.Duration(duration),
		RRule:        rrule,
	}
}

// MayBeUpdatedBy checks if the given user has permission to update the event.
// Currently, only designers may update events they facilitate.
func (e *Event) MayBeUpdatedBy(user *User) bool {
	if user == nil {
		return false
	}

	for _, facilitator := range e.Facilitators {
		if facilitator.ID == user.ID {
			return true
		}
	}

	return false
}

// Rename will change the title of the event. Blank names are not allowed.
func (e *Event) Rename(newTitle string) {
	if newTitle != "" && e.Title != newTitle {
		e.Title = newTitle
	}
}

// UpdateType will change the type for the event. Invalid types are not allowed.
func (e *Event) UpdateType(nt string) error {
	newType, err := event.TypeFromString(nt)
	if err != nil {
		return err
	}

	e.Type = newType

	switch e.Type {
	case event.Remote:
		e.Location = ""
	case event.InPerson:
		e.URL = ""
	}

	return nil
}

// UpdateDetails will change the details for the event. Blank details are not allowed.
func (e *Event) UpdateDetails(newDetails string) {
	if newDetails != "" && e.Details != newDetails {
		e.Details = newDetails
	}
}

// AddFacilitator will include the provider user as a facilitator on this event.
func (e *Event) AddFacilitator(facilitator *User) {
	if facilitator == nil {
		return
	}

	if e.Facilitators == nil {
		e.Facilitators = []User{}
	}

	for _, d := range e.Facilitators {
		if d.ID == facilitator.ID {
			return
		}
	}

	e.Facilitators = append(e.Facilitators, *facilitator)
}

// ReplaceFacilitators will overwrite the existing facilitator list with the newly provided one
func (e *Event) ReplaceFacilitators(facilitators []User) {
	e.Facilitators = nil
	for _, facilitator := range facilitators {
		e.AddFacilitator(&facilitator)
	}
}

// UpdateURL replaces the existing URL
func (e *Event) UpdateURL(newURL string) {
	e.URL = newURL
}

// UpdateLocation replaces the existing Location
func (e *Event) UpdateLocation(newLocation string) {
	e.Location = newLocation
}

// UpdateDuration replaces the existing Duration
func (e *Event) UpdateDuration(newDuration int64) {
	e.Duration = time.Duration(newDuration)
}

// UpdateRRule replaces the existing RRule
func (e *Event) UpdateRRule(newRRule string) {
	e.RRule = newRRule
}
