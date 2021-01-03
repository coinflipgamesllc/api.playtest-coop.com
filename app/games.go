package app

import (
	"errors"

	"github.com/coinflipgamesllc/api.playtest-coop.com/domain"
	"github.com/coinflipgamesllc/api.playtest-coop.com/domain/game"
	"go.uber.org/zap"
)

type (
	// GameService handles general interactions with games
	GameService struct {
		GameRepository domain.GameRepository
		UserRepository domain.UserRepository
		Logger         *zap.Logger
	}

	// Request DTOs

	// ListGamesRequest query params
	ListGamesRequest struct {
		Title       string `form:"title" example:"New Game"`
		Status      string `form:"status" example:"Prototype"`
		Designer    string `form:"designer" example:"Designer McDesignerton"`
		Owner       uint   `form:"owner" example:"123"`
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
		Title     string   `json:"title"`
		Overview  string   `json:"overview"`
		Status    string   `json:"status"`
		Designers []uint   `json:"designers"`
		Stats     *Stats   `json:"stats" binding:"omitempty,dive"`
		Mechanics []string `json:"mechanics" example:"['Hidden Movement', 'Worker Placement']"`
		TTSMod    int      `json:"tts_mod" example:"12345678"`
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

	// RulesResponse wrapper around a collection of rules sections
	RulesResponse struct {
		Rules []game.RulesSection `json:"rules"`
	}

	// ListMechanicsResponse wrapper for a listing of mechanics
	ListMechanicsResponse struct {
		Mechanics []string `json:"mechanics" example:"['trick-taking', 'worker placement', ...]"`
	}
)

// ListGames returns all games matching the specified query. The results are paginated
func (s *GameService) ListGames(req *ListGamesRequest) ([]domain.Game, int, error) {
	// Limit our limit
	if req.Limit == 0 {
		req.Limit = 10
	}
	if req.Limit > 100 {
		req.Limit = 100
	}

	// Fetch games
	games, total, err := s.GameRepository.ListGames(
		req.Title,
		req.Status,
		req.Designer,
		req.Owner,
		req.PlayerCount,
		req.Age,
		req.Playtime,
		req.Limit,
		req.Offset,
		req.Sort,
	)

	if err != nil {
		s.Logger.Error(err.Error())
		return nil, 0, err
	}

	return games, total, nil
}

// CreateGame creates a new stub game
func (s *GameService) CreateGame(req *CreateGameRequest, userID uint) (*domain.Game, error) {
	user, err := s.UserRepository.UserOfID(userID)
	if err != nil {
		s.Logger.Error(err.Error())
		return nil, err
	}

	game := domain.NewGame(req.Title, *user)

	// If the request included optional information, add it now
	if req.Overview != "" {
		game.UpdateOverview(req.Overview)
	}

	if len(req.Designers) > 1 { // Index 0 is always the current user, which is included already
		for _, designerID := range req.Designers {
			if designerID == user.ID {
				continue
			}

			designer, err := s.UserRepository.UserOfID(designerID)
			if err != nil {
				s.Logger.Error(err.Error())
				return nil, err
			}

			game.AddDesigner(designer)
		}
	}

	if req.Stats != nil {
		game.UpdateStats(req.Stats.MinPlayers, req.Stats.MaxPlayers, req.Stats.MinAge, req.Stats.EstimatedPlaytime)
	}

	// And save
	err = s.GameRepository.Save(game)
	if err != nil {
		s.Logger.Error(err.Error())
		return nil, err
	}

	return game, nil
}

// GetGame returns a specific game
func (s *GameService) GetGame(gameID uint) (*domain.Game, error) {
	game, err := s.GameRepository.GameOfID(gameID)
	if err != nil {
		s.Logger.Error(err.Error())
		return nil, err
	}

	return game, nil
}

// GetRules returns rules for a specific game
func (s *GameService) GetRules(gameID uint) ([]game.RulesSection, error) {
	rules, err := s.GameRepository.RulesOfGame(gameID)
	if err != nil {
		s.Logger.Error(err.Error())
		return nil, err
	}

	return rules, nil
}

// UpdateGame updates a specific game
func (s *GameService) UpdateGame(gameID uint, req *UpdateGameRequest, userID uint) (*domain.Game, error) {
	game, err := s.GameRepository.GameOfID(gameID)
	if err != nil {
		s.Logger.Error(err.Error())
		return nil, err
	}

	if game == nil {
		return nil, errors.New("game not found")
	}

	// Ensure that our current user is allowed to edit the game
	user, err := s.UserRepository.UserOfID(userID)
	if err != nil {
		s.Logger.Error(err.Error())
		return nil, err
	}

	if !game.MayBeUpdatedBy(user) {
		return nil, errors.New("you may not edit this game")
	}

	// Update game
	if req.Title != "" {
		game.Rename(req.Title)
	}

	if req.Overview != "" {
		game.UpdateOverview(req.Overview)
	}

	if req.Status != "" {
		err := game.UpdateStatus(req.Status)
		if err != nil {
			s.Logger.Error(err.Error())
			return nil, err
		}
	}

	if len(req.Designers) > 0 {
		designers := []domain.User{}
		for _, designerID := range req.Designers {
			designer, err := s.UserRepository.UserOfID(designerID)
			if err != nil {
				s.Logger.Error(err.Error())
				return nil, err
			}

			designers = append(designers, *designer)
		}

		game.ReplaceDesigners(designers)
	}

	if req.Stats != nil {
		game.UpdateStats(req.Stats.MinPlayers, req.Stats.MaxPlayers, req.Stats.MinAge, req.Stats.EstimatedPlaytime)
	}

	if req.Mechanics != nil {
		game.ReplaceMechanics(req.Mechanics)
	}

	if req.TTSMod != 0 {
		game.LinkTabletopSimulatorMod(req.TTSMod)
	}

	// And save
	err = s.GameRepository.Save(game)
	if err != nil {
		s.Logger.Error(err.Error())
		return nil, err
	}

	return game, nil
}

func (s *GameService) ListAvailableMechanics() []string {
	return game.AvailableMechanics()
}
