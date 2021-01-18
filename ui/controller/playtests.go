package controller

import (
	"strconv"

	"github.com/coinflipgamesllc/api.playtest-coop.com/app"
	"github.com/gin-gonic/gin"
)

// PlaytestController handles /playtests routes
type PlaytestController struct {
	PlaytestService *app.PlaytestService
}

// PlaytestsOnDate returns all the playtests scheduled for the provided date. Optionally by event.
// @Summary Return all the playtests scheduled for the provided date. Optionally by event.
// @Accept json
// @Produce json
// @Param query query app.ListPlaytestsRequest false "Filters for playtests"
// @Success 200 {object} app.ListPlaytestsResponse
// @Failure 400 {object} ValidationErrorResponse
// @Failure 400 {object} RequestErrorResponse
// @Failure 500 {object} ServerErrorResponse
// @Tags playtests
// @Router /playtests [get]
func (t *PlaytestController) PlaytestsOnDate(c *gin.Context) {
	// Validate request
	var req app.ListPlaytestsRequest
	if err := c.ShouldBind(&req); err != nil {
		validationErrorResponse(c, err)
		return
	}

	// Fetch playtests
	playtests, err := t.PlaytestService.ListPlaytests(&req)

	if err != nil {
		serverErrorResponse(c, "failed to fetch playtests")
		return
	}

	c.JSON(200, app.ListPlaytestsResponse{Playtests: playtests})
}

// RegisterGame schedules a playtest for a particular game
// @Summary Schedule a playtest for a particular game
// @Accept json
// @Produce json
// @Param game body app.RegisterGameRequest true "Playtest registration data"
// @Success 201 {object} app.PlaytestResponse
// @Failure 400 {object} ValidationErrorResponse
// @Failure 400 {object} RequestErrorResponse
// @Failure 500 {object} ServerErrorResponse
// @Tags playtests
// @Router /playtests/register-game [post]
func (t *PlaytestController) RegisterGame(c *gin.Context) {
	// Validate request
	var req app.RegisterGameRequest
	if err := c.ShouldBind(&req); err != nil {
		validationErrorResponse(c, err)
		return
	}

	// Register our game
	userID := userID(c)
	playtest, err := t.PlaytestService.RegisterGame(&req, userID)
	if err != nil {
		serverErrorResponse(c, "failed to register game")
		return
	}

	c.JSON(201, app.PlaytestResponse{Playtest: playtest})
}

// AssignLocation assigns a playtest to a table (real or virtual)
// @Summary Assign a playtest to a table (real or virtual)
// @Accept json
// @Produce json
// @Param id path integer true "Playtest ID"
// @Param event body app.AssignPlaytestLocationRequest false "Table"
// @Success 200 {object} app.PlaytestResponse
// @Failure 400 {object} ValidationErrorResponse
// @Failure 400 {object} RequestErrorResponse
// @Failure 500 {object} ServerErrorResponse
// @Tags playtests
// @Router /playtests/:id/location [put]
func (t *PlaytestController) AssignLocation(c *gin.Context) {
	// Pull playtest by ID
	playtestID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		requestErrorResponse(c, err.Error())
		return
	}

	// Validate request
	var req app.AssignPlaytestLocationRequest
	if err := c.ShouldBind(&req); err != nil {
		validationErrorResponse(c, err)
		return
	}

	// Assign the location
	userID := userID(c)
	playtest, err := t.PlaytestService.AssignLocation(uint(playtestID), &req, userID)
	if err != nil {
		serverErrorResponse(c, "failed to assign location")
		return
	}

	c.JSON(200, app.PlaytestResponse{Playtest: playtest})
}

// AddPlayer adds a player to the playtest
// @Summary adds a player to the playtest
// @Accept json
// @Produce json
// @Param id path integer true "Playtest ID"
// @Success 200 {object} app.PlaytestResponse
// @Failure 400 {object} ValidationErrorResponse
// @Failure 400 {object} RequestErrorResponse
// @Failure 500 {object} ServerErrorResponse
// @Tags playtests
// @Router /playtests/:id/location [put]
func (t *PlaytestController) AddPlayer(c *gin.Context) {
	// Pull playtest by ID
	playtestID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		requestErrorResponse(c, err.Error())
		return
	}

	// Assign the location
	userID := userID(c)
	playtest, err := t.PlaytestService.AddPlayer(uint(playtestID), userID)
	if err != nil {
		serverErrorResponse(c, "failed to join playtest")
		return
	}

	c.JSON(200, app.PlaytestResponse{Playtest: playtest})
}

// RemovePlayer removes a player from a playtest
// @Summary removes a player from a playtest
// @Accept json
// @Produce json
// @Param id path integer true "Playtest ID"
// @Success 200 {object} app.PlaytestResponse
// @Failure 400 {object} ValidationErrorResponse
// @Failure 400 {object} RequestErrorResponse
// @Failure 500 {object} ServerErrorResponse
// @Tags playtests
// @Router /playtests/:id/location [put]
func (t *PlaytestController) RemovePlayer(c *gin.Context) {
	// Pull playtest by ID
	playtestID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		requestErrorResponse(c, err.Error())
		return
	}

	// Assign the location
	userID := userID(c)
	playtest, err := t.PlaytestService.RemovePlayer(uint(playtestID), userID)
	if err != nil {
		serverErrorResponse(c, "failed to leave playtest")
		return
	}

	c.JSON(200, app.PlaytestResponse{Playtest: playtest})
}

// Start marks the time the playtest started
// @Summary  marks the time the playtest started
// @Accept json
// @Produce json
// @Param id path integer true "Playtest ID"
// @Success 200 {object} app.PlaytestResponse
// @Failure 400 {object} ValidationErrorResponse
// @Failure 400 {object} RequestErrorResponse
// @Failure 500 {object} ServerErrorResponse
// @Tags playtests
// @Router /playtests/:id/start [put]
func (t *PlaytestController) Start(c *gin.Context) {
	// Pull playtest by ID
	playtestID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		requestErrorResponse(c, err.Error())
		return
	}

	// Assign the location
	userID := userID(c)
	playtest, err := t.PlaytestService.StartPlaytest(uint(playtestID), userID)
	if err != nil {
		serverErrorResponse(c, "failed to assign location")
		return
	}

	c.JSON(200, app.PlaytestResponse{Playtest: playtest})
}

// StartFeedback marks the time the playtesters started giving feedback
// @Summary  marks the time the playtesters started giving feedback
// @Accept json
// @Produce json
// @Param id path integer true "Playtest ID"
// @Success 200 {object} app.PlaytestResponse
// @Failure 400 {object} ValidationErrorResponse
// @Failure 400 {object} RequestErrorResponse
// @Failure 500 {object} ServerErrorResponse
// @Tags playtests
// @Router /playtests/:id/start-feedback [put]
func (t *PlaytestController) StartFeedback(c *gin.Context) {
	// Pull playtest by ID
	playtestID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		requestErrorResponse(c, err.Error())
		return
	}

	// Assign the location
	userID := userID(c)
	playtest, err := t.PlaytestService.StartFeedback(uint(playtestID), userID)
	if err != nil {
		serverErrorResponse(c, "failed to assign location")
		return
	}

	c.JSON(200, app.PlaytestResponse{Playtest: playtest})
}

// Finish marks the time the playtest ended
// @Summary  marks the time the playtest ended
// @Accept json
// @Produce json
// @Param id path integer true "Playtest ID"
// @Success 200 {object} app.PlaytestResponse
// @Failure 400 {object} ValidationErrorResponse
// @Failure 400 {object} RequestErrorResponse
// @Failure 500 {object} ServerErrorResponse
// @Tags playtests
// @Router /playtests/:id/finish [put]
func (t *PlaytestController) Finish(c *gin.Context) {
	// Pull playtest by ID
	playtestID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		requestErrorResponse(c, err.Error())
		return
	}

	// Assign the location
	userID := userID(c)
	playtest, err := t.PlaytestService.FinishPlaytest(uint(playtestID), userID)
	if err != nil {
		serverErrorResponse(c, "failed to assign location")
		return
	}

	c.JSON(200, app.PlaytestResponse{Playtest: playtest})
}
