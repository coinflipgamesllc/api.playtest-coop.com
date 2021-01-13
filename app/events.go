package app

import (
	"errors"

	"github.com/coinflipgamesllc/api.playtest-coop.com/domain"
	"github.com/coinflipgamesllc/api.playtest-coop.com/domain/event"
	"go.uber.org/zap"
)

type (
	// EventService handles general interactions with events
	EventService struct {
		EventRepository domain.EventRepository
		UserRepository  domain.UserRepository
		Logger          *zap.Logger
	}

	// Request DTOs

	// CreateEventRequest params for creating an event
	CreateEventRequest struct {
		Title    string `json:"title" binding:"required"`
		Details  string `json:"details" binding:"required"`
		Type     string `json:"type" binding:"required"`
		URL      string `json:"url" binding:"omitempty,url"`
		Location string `json:"location"`
		Duration int64  `json:"duration"`
		RRule    string `json:"rrule" binding:"required"`
	}

	// UpdateEventRequest params for updating an event
	UpdateEventRequest struct {
		Title        string `json:"title"`
		Type         string `json:"type"`
		Details      string `json:"details"`
		Facilitators []uint `json:"facilitators"`
		URL          string `json:"url" binding:"omitempty,url"`
		Location     string `json:"location"`
		Duration     int64  `json:"duration"`
		RRule        string `json:"rrule"`
	}

	// Response DTOs

	// ListEventsResponse paginated events list
	ListEventsResponse struct {
		Events []domain.Event `json:"events"`
	}

	// EventResponse wrapper around an event
	EventResponse struct {
		Event *domain.Event `json:"event"`
	}
)

// ListEvents returns all events matching the specified query. The results are paginated
func (s *EventService) ListEvents() ([]domain.Event, error) {
	// Fetch events
	events, err := s.EventRepository.ListEvents()

	if err != nil {
		s.Logger.Error(err.Error())
		return nil, err
	}

	return events, nil
}

// CreateEvent creates a new stub event
func (s *EventService) CreateEvent(req *CreateEventRequest, userID uint) (*domain.Event, error) {
	user, err := s.UserRepository.UserOfID(userID)
	if err != nil {
		s.Logger.Error(err.Error())
		return nil, err
	}

	var e *domain.Event
	switch req.Type {
	case string(event.Remote):
		e = domain.NewRemoteEvent(req.Title, req.Details, req.URL, req.Duration, req.RRule, *user)

	case string(event.InPerson):
		e = domain.NewInPersonEvent(req.Title, req.Details, req.Location, req.Duration, req.RRule, *user)
	}

	// And save
	err = s.EventRepository.Save(e)
	if err != nil {
		s.Logger.Error(err.Error())
		return nil, err
	}

	return e, nil
}

// GetEvent returns a specific event
func (s *EventService) GetEvent(eventID uint) (*domain.Event, error) {
	e, err := s.EventRepository.EventOfID(eventID)
	if err != nil {
		s.Logger.Error(err.Error())
		return nil, err
	}

	return e, nil
}

// UpdateEvent updates a specific event
func (s *EventService) UpdateEvent(eventID uint, req *UpdateEventRequest, userID uint) (*domain.Event, error) {
	e, err := s.EventRepository.EventOfID(eventID)
	if err != nil {
		s.Logger.Error(err.Error())
		return nil, err
	}

	if e == nil {
		return nil, errors.New("event not found")
	}

	// Ensure that our current user is allowed to edit the event
	user, err := s.UserRepository.UserOfID(userID)
	if err != nil {
		s.Logger.Error(err.Error())
		return nil, err
	}

	if !e.MayBeUpdatedBy(user) {
		return nil, errors.New("you may not edit this event")
	}

	// Update event
	if req.Title != "" {
		e.Rename(req.Title)
	}

	if req.Details != "" {
		e.UpdateDetails(req.Details)
	}

	if len(req.Facilitators) > 0 {
		facilitators := []domain.User{}
		for _, facilitatorID := range req.Facilitators {
			facilitator, err := s.UserRepository.UserOfID(facilitatorID)
			if err != nil {
				s.Logger.Error(err.Error())
				return nil, err
			}

			facilitators = append(facilitators, *facilitator)
		}

		e.ReplaceFacilitators(facilitators)
	}

	if req.Type != "" {
		err := e.UpdateType(req.Type)
		if err != nil {
			s.Logger.Error(err.Error())
			return nil, err
		}
	}

	if req.URL != "" {
		e.UpdateURL(req.URL)
	}

	if req.Location != "" {
		e.UpdateLocation(req.Location)
	}

	if req.Duration != 0 {
		e.UpdateDuration(req.Duration)
	}

	if req.RRule != "" {
		e.UpdateRRule(req.RRule)
	}

	// And save
	err = s.EventRepository.Save(e)
	if err != nil {
		s.Logger.Error(err.Error())
		return nil, err
	}

	return e, nil
}
