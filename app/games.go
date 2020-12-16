package app

import (
	"errors"

	"github.com/coinflipgamesllc/api.playtest-coop.com/domain"
	"go.uber.org/zap"
)

type (
	// GameService handles general interactions with games
	GameService struct {
		GameRepository domain.GameRepository
		UserRepository domain.UserRepository
		Logger         *zap.SugaredLogger
	}

	// Request DTOs

	// ListGamesRequest query params
	ListGamesRequest struct {
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

	// Stats wrapper for game stats
	Stats struct {
		MinPlayers        int `json:"min_players" binding:"min=0,ltefield=MaxPlayers" example:"1"`
		MaxPlayers        int `json:"max_players" binding:"min=0,gtefield=MinPlayers" example:"5"`
		MinAge            int `json:"min_age" binding:"min=0,max=99" example:"8"`
		EstimatedPlaytime int `json:"estimated_playtime" binding:"min=0,max=9999" example:"30"`
	}

	// CreateGameRequest params for creating a game
	CreateGameRequest struct {
		Title     string `json:"title" binding:"required"`
		Overview  string `json:"overview"`
		Designers []uint `json:"designers"`
		Stats     *Stats `json:"stats" binding:"omitempty,dive"`
	}

	// UpdateGameRequest params for updating a game
	UpdateGameRequest struct {
		Title     string `json:"title"`
		Overview  string `json:"overview"`
		Status    string `json:"status"`
		Designers []uint `json:"designers"`
		Stats     *Stats `json:"stats" binding:"omitempty,dive"`
	}

	// Response DTOs

	// ListGamesResponse paginated games list
	ListGamesResponse struct {
		Games  []domain.Game `json:"games"`
		Total  int           `json:"total" example:"1000"`
		Limit  int           `json:"limit" example:"100"`
		Offset int           `json:"offset" example:"50"`
	}

	// GameResponse wrapper around a game
	GameResponse struct {
		Game *domain.Game `json:"game"`
	}
)

// ListGames returns all games matching the specified query. The results are paginated
func (s *GameService) ListGames(title, status, designer string, playerCount, age, playtime, limit, offset int, sort string) ([]domain.Game, int, error) {
	// Limit our limit
	if limit == 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	// Fetch games
	games, total, err := s.GameRepository.ListGames(
		title,
		status,
		designer,
		playerCount,
		age,
		playtime,
		limit,
		offset,
		sort,
	)

	if err != nil {
		s.Logger.Error(err)
		return nil, 0, err
	}

	return games, total, nil
}

// CreateGame creates a new stub game
func (s *GameService) CreateGame(title, overview string, designers []uint, minPlayers, maxPlayers, minAge, estimatedPlaytime int, userID uint) (*domain.Game, error) {
	user, err := s.UserRepository.UserOfID(userID)
	if err != nil {
		s.Logger.Error(err)
		return nil, err
	}

	game := domain.NewGame(title, *user)

	// If the request included optional information, add it now
	if overview != "" {
		game.UpdateOverview(overview)
	}

	if len(designers) > 1 { // Index 0 is always the current user, which is included already
		for _, designerID := range designers {
			if designerID == user.ID {
				continue
			}

			designer, err := s.UserRepository.UserOfID(designerID)
			if err != nil {
				s.Logger.Error(err)
				return nil, err
			}

			game.AddDesigner(designer)
		}
	}

	game.UpdateStats(minPlayers, maxPlayers, minAge, estimatedPlaytime)

	// And save
	err = s.GameRepository.Save(game)
	if err != nil {
		s.Logger.Error(err)
		return nil, err
	}

	return game, nil
}

// GetGame returns a specific game
func (s *GameService) GetGame(gameID uint) (*domain.Game, error) {
	game, err := s.GameRepository.GameOfID(gameID)
	if err != nil {
		s.Logger.Error(err)
		return nil, err
	}

	return game, nil
}

// UpdateGame updates a specific game
func (s *GameService) UpdateGame(gameID uint, title, overview, status string, designers []uint, minPlayers, maxPlayers, minAge, estimatedPlaytime int, userID uint) (*domain.Game, error) {
	game, err := s.GameRepository.GameOfID(gameID)
	if err != nil {
		s.Logger.Error(err)
		return nil, err
	}

	if game == nil {
		return nil, errors.New("game not found")
	}

	// Ensure that our current user is allowed to edit the game
	user, err := s.UserRepository.UserOfID(userID)
	if err != nil {
		s.Logger.Error(err)
		return nil, err
	}

	if !game.MayBeUpdatedBy(user) {
		return nil, errors.New("you may not edit this game")
	}

	// Update game
	if title != "" {
		game.Rename(title)
	}

	if overview != "" {
		game.UpdateOverview(overview)
	}

	if status != "" {
		err := game.UpdateStatus(status)
		if err != nil {
			s.Logger.Error(err)
			return nil, err
		}
	}

	if len(designers) > 0 {
		des := []domain.User{}
		for _, designerID := range designers {
			designer, err := s.UserRepository.UserOfID(designerID)
			if err != nil {
				s.Logger.Error(err)
				return nil, err
			}

			des = append(des, *designer)
		}

		game.ReplaceDesigners(des)
	}

	game.UpdateStats(minPlayers, maxPlayers, minAge, estimatedPlaytime)

	// And save
	err = s.GameRepository.Save(game)
	if err != nil {
		s.Logger.Error(err)
		return nil, err
	}

	return game, nil
}
