package controller

import (
	"strconv"

	"github.com/coinflipgamesllc/api.playtest-coop.com/app"
	"github.com/coinflipgamesllc/api.playtest-coop.com/domain"
	"github.com/gin-gonic/gin"
)

// GameController handles /games routes
type GameController struct {
	GameService *app.GameService
}

// ListGamesRequest query params
type ListGamesRequest struct {
	Title       string `form:"title" example:"New Game"`
	Status      string `form:"status" example:"Prototype"`
	Designer    string `form:"designer" example:"Designer McDesignerton"`
	PlayerCount int    `form:"player_count" example:"2"`
	Age         int    `form:"age" example:"13"`
	Playtime    int    `form:"playtime" example:"30"`
	Limit       int    `form:"limit" example:"100"`
	Offset      int    `form:"offset" example:"50"`
	Sort        string `form:"sort" example:"name,desc"`
}

// ListGamesResponse paginated games list
type ListGamesResponse struct {
	Games  []domain.Game `json:"games"`
	Total  int           `json:"total" example:"1000"`
	Limit  int           `json:"limit" example:"100"`
	Offset int           `json:"offset" example:"50"`
}

// ListGames list games matching the query with pagination
// @Summary List games matching the query with pagination
// @Accept json
// @Produce json
// @Param query query ListGamesRequest false "Filters for games"
// @Success 200 {object} ListGamesResponse
// @Failure 400 {object} RequestErrorResponse
// @Failure 500 {object} ServerErrorResponse
// @Tags games
// @Router /games [get]
func (t *GameController) ListGames(c *gin.Context) {
	// Validate request
	var req ListGamesRequest
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

	c.JSON(200, ListGamesResponse{Games: games, Total: total, Limit: req.Limit, Offset: req.Offset})
}

// Stats wrapper for game stats
type Stats struct {
	MinPlayers        int `json:"min_players" binding:"min=0,ltefield=MaxPlayers" example:"1"`
	MaxPlayers        int `json:"max_players" binding:"min=0,gtefield=MinPlayers" example:"5"`
	MinAge            int `json:"min_age" binding:"min=0,max=99" example:"8"`
	EstimatedPlaytime int `json:"estimated_playtime" binding:"min=0,max=9999" example:"30"`
}

// CreateGameRequest params for creating a game
type CreateGameRequest struct {
	Title     string `json:"title" binding:"required"`
	Overview  string `json:"overview"`
	Designers []uint `json:"designers"`
	Stats     *Stats `json:"stats" binding:"omitempty,dive"`
}

// GameResponse wrapper around a game
type GameResponse struct {
	Game *domain.Game `json:"game"`
}

// CreateGame creates a new stub game
// @Summary Create a new stub game
// @Accept json
// @Produce json
// @Param game body CreateGameRequest true "Game data"
// @Success 200 {object} GameResponse
// @Failure 400 {object} RequestErrorResponse
// @Failure 500 {object} ServerErrorResponse
// @Tags games
// @Router /games [post]
func (t *GameController) CreateGame(c *gin.Context) {
	// Validate request
	var req CreateGameRequest
	if err := c.ShouldBind(&req); err != nil {
		requestErrorResponse(c, err.Error())
		return
	}

	if req.Stats == nil {
		req.Stats = &Stats{}
	}

	// Create our new game
	userID := userID(c)
	game, err := t.GameService.CreateGame(req.Title, req.Overview, req.Designers, req.Stats.MinPlayers, req.Stats.MaxPlayers, req.Stats.MinAge, req.Stats.EstimatedPlaytime, userID)
	if err != nil {
		serverErrorResponse(c, "failed to create game")
		return
	}

	c.JSON(200, GameResponse{Game: game})
}

// GetGame returns a specific game by id
// @Summary Return a specific game by id
// @Produce json
// @Param id path integer true "Game ID"
// @Success 200 {object} GameResponse
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

	c.JSON(200, GameResponse{Game: game})
}

// UpdateGameRequest params for updating a game
type UpdateGameRequest struct {
	Title     string `json:"title"`
	Overview  string `json:"overview"`
	Status    string `json:"status"`
	Designers []uint `json:"designers"`
	Stats     *Stats `json:"stats" binding:"omitempty,dive"`
}

// UpdateGame updates a specific game
// @Summary Update a specific game
// @Accept json
// @Produce json
// @Param id path integer true "Game ID"
// @Param game body UpdateGameRequest false "Game data"
// @Success 200 {object} GameResponse
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
	var req UpdateGameRequest
	if err := c.ShouldBind(&req); err != nil {
		requestErrorResponse(c, err.Error())
		return
	}

	game, err := t.GameService.UpdateGame(uint(gameID), req.Title, req.Overview, req.Status, req.Designers, req.Stats.MinPlayers, req.Stats.MaxPlayers, req.Stats.MinAge, req.Stats.EstimatedPlaytime, userID)
	if err != nil {
		serverErrorResponse(c, "failed to update game")
		return
	}

	c.JSON(200, GameResponse{Game: game})
}
