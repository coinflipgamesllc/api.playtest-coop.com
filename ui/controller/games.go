package controller

import (
	"strconv"

	"github.com/coinflipgamesllc/api.playtest-coop.com/app"
	"github.com/gin-gonic/gin"
)

// GameController handles /games routes
type GameController struct {
	GameService *app.GameService
}

// ListGames list games matching the query with pagination
// @Summary List games matching the query with pagination
// @Accept json
// @Produce json
// @Param query query app.ListGamesRequest false "Filters for games"
// @Success 200 {object} app.ListGamesResponse
// @Failure 400 {object} RequestErrorResponse
// @Failure 500 {object} ServerErrorResponse
// @Tags games
// @Router /games [get]
func (t *GameController) ListGames(c *gin.Context) {
	// Validate request
	var req app.ListGamesRequest
	if err := c.ShouldBind(&req); err != nil {
		requestErrorResponse(c, err.Error())
		return
	}

	// Fetch games
	games, total, err := t.GameService.ListGames(
		req.Title,
		req.Status,
		req.Designer,
		req.PlayerCount,
		req.Age,
		req.Playtime,
		req.Limit,
		req.Offset,
		req.Sort,
	)

	if err != nil {
		serverErrorResponse(c, "failed to fetch games")
		return
	}

	c.JSON(200, app.ListGamesResponse{Games: games, Total: total, Limit: req.Limit, Offset: req.Offset})
}

// CreateGame creates a new stub game
// @Summary Create a new stub game
// @Accept json
// @Produce json
// @Param game body app.CreateGameRequest true "Game data"
// @Success 200 {object} app.GameResponse
// @Failure 400 {object} RequestErrorResponse
// @Failure 500 {object} ServerErrorResponse
// @Tags games
// @Router /games [post]
func (t *GameController) CreateGame(c *gin.Context) {
	// Validate request
	var req app.CreateGameRequest
	if err := c.ShouldBind(&req); err != nil {
		requestErrorResponse(c, err.Error())
		return
	}

	if req.Stats == nil {
		req.Stats = &app.Stats{}
	}

	// Create our new game
	userID := userID(c)
	game, err := t.GameService.CreateGame(req.Title, req.Overview, req.Designers, req.Stats.MinPlayers, req.Stats.MaxPlayers, req.Stats.MinAge, req.Stats.EstimatedPlaytime, userID)
	if err != nil {
		serverErrorResponse(c, "failed to create game")
		return
	}

	c.JSON(200, app.GameResponse{Game: game})
}

// GetGame returns a specific game by id
// @Summary Return a specific game by id
// @Produce json
// @Param id path integer true "Game ID"
// @Success 200 {object} app.GameResponse
// @Failure 400 {object} RequestErrorResponse
// @Failure 500 {object} ServerErrorResponse
// @Tags games
// @Router /games/:id [get]
func (t *GameController) GetGame(c *gin.Context) {
	// Validate request
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		requestErrorResponse(c, err.Error())
		return
	}

	game, err := t.GameService.GetGame(uint(id))
	if err != nil {
		serverErrorResponse(c, "failed to fetch game")
		return
	}

	if game == nil {
		notFoundResponse(c, "game not found")
		return
	}

	c.JSON(200, app.GameResponse{Game: game})
}

// UpdateGame updates a specific game
// @Summary Update a specific game
// @Accept json
// @Produce json
// @Param id path integer true "Game ID"
// @Param game body app.UpdateGameRequest false "Game data"
// @Success 200 {object} app.GameResponse
// @Failure 400 {object} RequestErrorResponse
// @Failure 500 {object} ServerErrorResponse
// @Tags games
// @Router /games/:id [put]
func (t *GameController) UpdateGame(c *gin.Context) {
	// Pull game by ID
	gameID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		requestErrorResponse(c, err.Error())
		return
	}

	userID := userID(c)

	// Validate the request itself
	var req app.UpdateGameRequest
	if err := c.ShouldBind(&req); err != nil {
		requestErrorResponse(c, err.Error())
		return
	}

	game, err := t.GameService.UpdateGame(uint(gameID), req.Title, req.Overview, req.Status, req.Designers, req.Stats.MinPlayers, req.Stats.MaxPlayers, req.Stats.MinAge, req.Stats.EstimatedPlaytime, userID)
	if err != nil {
		serverErrorResponse(c, "failed to update game")
		return
	}

	c.JSON(200, app.GameResponse{Game: game})
}
