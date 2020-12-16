package app

import (
	"errors"

	"github.com/coinflipgamesllc/api.playtest-coop.com/domain"
	"go.uber.org/zap"
)

// GameService handles general interactions with games
type GameService struct {
	GameRepository domain.GameRepository
	UserRepository domain.UserRepository
	Logger         *zap.SugaredLogger
}

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
