package controller

import (
	"strconv"

	"github.com/coinflipgamesllc/api.playtest-coop.com/app"
	"github.com/gin-gonic/gin"
)

// EventController handles /events routes
type EventController struct {
	EventService *app.EventService
}

// ListEvents list all events
// @Summary List all events
// @Accept json
// @Produce json
// @Success 200 {object} app.ListEventsResponse
// @Failure 400 {object} ValidationErrorResponse
// @Failure 400 {object} RequestErrorResponse
// @Failure 500 {object} ServerErrorResponse
// @Tags events
// @Router /events [get]
func (t *EventController) ListEvents(c *gin.Context) {
	// Fetch events
	events, err := t.EventService.ListEvents()

	if err != nil {
		serverErrorResponse(c, "failed to fetch events")
		return
	}

	c.JSON(200, app.ListEventsResponse{Events: events})
}

// CreateEvent creates a new stub event
// @Summary Create a new stub event
// @Accept json
// @Produce json
// @Param event body app.CreateEventRequest true "Event data"
// @Success 201 {object} app.EventResponse
// @Failure 400 {object} ValidationErrorResponse
// @Failure 400 {object} RequestErrorResponse
// @Failure 500 {object} ServerErrorResponse
// @Tags events
// @Router /events [post]
func (t *EventController) CreateEvent(c *gin.Context) {
	// Validate request
	var req app.CreateEventRequest
	if err := c.ShouldBind(&req); err != nil {
		validationErrorResponse(c, err)
		return
	}

	// Create our new event
	userID := userID(c)
	event, err := t.EventService.CreateEvent(&req, userID)
	if err != nil {
		serverErrorResponse(c, "failed to create event")
		return
	}

	c.JSON(201, app.EventResponse{Event: event})
}

// GetEvent returns a specific event by id
// @Summary Return a specific event by id
// @Produce json
// @Param id path integer true "Event ID"
// @Success 200 {object} app.EventResponse
// @Failure 400 {object} RequestErrorResponse
// @Failure 500 {object} ServerErrorResponse
// @Tags events
// @Router /events/:id [get]
func (t *EventController) GetEvent(c *gin.Context) {
	// Validate request
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		requestErrorResponse(c, err.Error())
		return
	}

	event, err := t.EventService.GetEvent(uint(id))
	if err != nil {
		serverErrorResponse(c, "failed to fetch event")
		return
	}

	if event == nil {
		notFoundResponse(c, "event not found")
		return
	}

	c.JSON(200, app.EventResponse{Event: event})
}

// UpdateEvent updates a specific event
// @Summary Update a specific event
// @Accept json
// @Produce json
// @Param id path integer true "Event ID"
// @Param event body app.UpdateEventRequest false "Event data"
// @Success 200 {object} app.EventResponse
// @Failure 400 {object} ValidationErrorResponse
// @Failure 400 {object} RequestErrorResponse
// @Failure 500 {object} ServerErrorResponse
// @Tags events
// @Router /events/:id [put]
func (t *EventController) UpdateEvent(c *gin.Context) {
	// Pull event by ID
	eventID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		requestErrorResponse(c, err.Error())
		return
	}

	userID := userID(c)

	// Validate the request itself
	var req app.UpdateEventRequest
	if err := c.ShouldBind(&req); err != nil {
		validationErrorResponse(c, err)
		return
	}

	event, err := t.EventService.UpdateEvent(uint(eventID), &req, userID)
	if err != nil {
		serverErrorResponse(c, "failed to update event")
		return
	}

	c.JSON(200, app.EventResponse{Event: event})
}
