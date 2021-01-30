package app

import (
	"fmt"
	"time"

	"github.com/coinflipgamesllc/api.playtest-coop.com/domain"
	"go.uber.org/zap"
)

type (
	// PlaytestService handles general interactions with games
	PlaytestService struct {
		EventRepository    domain.EventRepository
		GameRepository     domain.GameRepository
		PlaytestRepository domain.PlaytestRepository
		UserRepository     domain.UserRepository
		Logger             *zap.Logger
	}

	// Request DTOs

	// ListPlaytestsRequest query params
	ListPlaytestsRequest struct {
		Date    string `form:"date" binding:"required"`
		EventID uint   `form:"event_id"`
	}

	// RegisterGameRequest params required for registering for a playtest
	RegisterGameRequest struct {
		GameID              uint   `json:"game" binding:"required"`
		EventID             uint   `json:"event"`
		Date                string `json:"date" binding:"required"`
		MinNumberOfPlayers  uint   `json:"min_players" binding:"required" example:"3"`
		MaxNumberOfPlayers  uint   `json:"max_players" binding:"required" example:"5"`
		Duration            uint   `json:"duration" binding:"required" example:"60"`
		DesignerWantsToPlay bool   `json:"designer_wants_to_play" binding:"required" example:"true"`
		HopingToTest        string `json:"hoping_to_test" binding:"required" example:"Is the kerpluxic mechanic intuitive?"`
		TTSServer           string `json:"tts_server" example:"server_name"`
		TTSPassword         string `json:"tts_password" example:"password"`
	}

	// AssignPlaytestLocationRequest wraps the table assignment for a playtest
	AssignPlaytestLocationRequest struct {
		Table string `json:"table" binding:"required" example:"1"`
	}

	// Response DTOs

	// ListPlaytestsResponse playtests wrapper
	ListPlaytestsResponse struct {
		Playtests []domain.Playtest `json:"playtests"`
	}

	// PlaytestResponse playtest wrapper
	PlaytestResponse struct {
		Playtest *domain.Playtest `json:"playtest"`
	}
)

// ListPlaytests returns all the playtests scheduled on the specified date. Optionally by event.
func (s *PlaytestService) ListPlaytests(req *ListPlaytestsRequest) ([]domain.Playtest, error) {
	// Fetch playtests
	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		s.Logger.Error(err.Error())
		return nil, err
	}

	playtests, err := s.PlaytestRepository.PlaytestsOnDate(date, req.EventID)
	if err != nil {
		s.Logger.Error(err.Error())
		return nil, err
	}

	return playtests, nil
}

// RegisterGame sets up a new playtest for a game at a specific time. It can optionally be tied to an event
func (s *PlaytestService) RegisterGame(req *RegisterGameRequest, userID uint) (*domain.Playtest, error) {
	// Pull up our user & game and make sure they're compatible
	user, err := s.UserRepository.UserOfID(userID)
	if err != nil {
		s.Logger.Error(err.Error())
		return nil, err
	}

	game, err := s.GameRepository.GameOfID(req.GameID)
	if err != nil {
		s.Logger.Error(err.Error())
		return nil, err
	}

	if user == nil || game == nil || !game.MayBeUpdatedBy(user) {
		return nil, fmt.Errorf("you're not allowed to register this game for playtesting")
	}

	// Create the playtest
	var event *domain.Event
	if req.EventID != 0 {
		event, err = s.EventRepository.EventOfID(req.EventID)
		if err != nil {
			s.Logger.Error(err.Error())
			return nil, err
		}
	}

	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		s.Logger.Error(err.Error())
		return nil, err
	}

	playtest := domain.RegisterGame(
		game,
		event,
		date,
		req.MinNumberOfPlayers,
		req.MaxNumberOfPlayers,
		req.Duration,
		req.DesignerWantsToPlay,
		req.HopingToTest,
		req.TTSServer,
		req.TTSPassword,
	)

	// And save
	err = s.PlaytestRepository.Save(playtest)
	if err != nil {
		s.Logger.Error(err.Error())
		return nil, err
	}

	return playtest, nil
}

// AssignLocation assigns a playtest to a table (real or virtual)
func (s *PlaytestService) AssignLocation(playtestID uint, req *AssignPlaytestLocationRequest, userID uint) (*domain.Playtest, error) {
	// Pull up our user & make sure they're allowed to assign locations
	user, err := s.UserRepository.UserOfID(userID)
	if err != nil {
		s.Logger.Error(err.Error())
		return nil, err
	}

	playtest, err := s.PlaytestRepository.PlaytestOfID(playtestID)
	if err != nil {
		s.Logger.Error(err.Error())
		return nil, err
	}

	if user == nil || playtest == nil {
		return nil, fmt.Errorf("you're not allowed to assign locations")
	}

	playtest.AssignTable(req.Table)

	// And save
	err = s.PlaytestRepository.Save(playtest)
	if err != nil {
		s.Logger.Error(err.Error())
		return nil, err
	}

	return playtest, nil
}

// AddPlayer adds a player to the playtest
func (s *PlaytestService) AddPlayer(playtestID uint, userID uint) (*domain.Playtest, error) {
	// Pull up our user & make sure they're allowed to join
	user, err := s.UserRepository.UserOfID(userID)
	if err != nil {
		s.Logger.Error(err.Error())
		return nil, err
	}

	playtest, err := s.PlaytestRepository.PlaytestOfID(playtestID)
	if err != nil {
		s.Logger.Error(err.Error())
		return nil, err
	}

	if user == nil || playtest == nil {
		return nil, fmt.Errorf("you're not allowed to join this playtest")
	}

	playtest.AddPlayer(user)

	// And save
	err = s.PlaytestRepository.Save(playtest)
	if err != nil {
		s.Logger.Error(err.Error())
		return nil, err
	}

	return playtest, nil
}

// RemovePlayer removes a player from a playtest
func (s *PlaytestService) RemovePlayer(playtestID uint, userID uint) (*domain.Playtest, error) {
	// Pull up our user & make sure they're allowed to assign locations
	user, err := s.UserRepository.UserOfID(userID)
	if err != nil {
		s.Logger.Error(err.Error())
		return nil, err
	}

	playtest, err := s.PlaytestRepository.PlaytestOfID(playtestID)
	if err != nil {
		s.Logger.Error(err.Error())
		return nil, err
	}

	if user == nil || playtest == nil {
		return nil, fmt.Errorf("you're not allowed to leave this playtest")
	}

	playtest.RemovePlayer(user)

	// And save
	err = s.PlaytestRepository.Save(playtest)
	if err != nil {
		s.Logger.Error(err.Error())
		return nil, err
	}

	return playtest, nil
}

// StartPlaytest will set the time the playtest started to now
func (s *PlaytestService) StartPlaytest(playtestID uint, userID uint) (*domain.Playtest, error) {
	// Pull up our user & make sure they're allowed to assign locations
	user, err := s.UserRepository.UserOfID(userID)
	if err != nil {
		s.Logger.Error(err.Error())
		return nil, err
	}

	playtest, err := s.PlaytestRepository.PlaytestOfID(playtestID)
	if err != nil {
		s.Logger.Error(err.Error())
		return nil, err
	}

	if user == nil || playtest == nil {
		return nil, fmt.Errorf("you're not allowed to assign locations")
	}

	playtest.Start()

	// And save
	err = s.PlaytestRepository.Save(playtest)
	if err != nil {
		s.Logger.Error(err.Error())
		return nil, err
	}

	return playtest, nil
}

// StartFeedback will set the feedback time to now
func (s *PlaytestService) StartFeedback(playtestID uint, userID uint) (*domain.Playtest, error) {
	// Pull up our user & make sure they're allowed to assign locations
	user, err := s.UserRepository.UserOfID(userID)
	if err != nil {
		s.Logger.Error(err.Error())
		return nil, err
	}

	playtest, err := s.PlaytestRepository.PlaytestOfID(playtestID)
	if err != nil {
		s.Logger.Error(err.Error())
		return nil, err
	}

	if user == nil || playtest == nil {
		return nil, fmt.Errorf("you're not allowed to assign locations")
	}

	playtest.StartFeedback()

	// And save
	err = s.PlaytestRepository.Save(playtest)
	if err != nil {
		s.Logger.Error(err.Error())
		return nil, err
	}

	return playtest, nil
}

// FinishPlaytest will set the time the playtest finished to now
func (s *PlaytestService) FinishPlaytest(playtestID uint, userID uint) (*domain.Playtest, error) {
	// Pull up our user & make sure they're allowed to assign locations
	user, err := s.UserRepository.UserOfID(userID)
	if err != nil {
		s.Logger.Error(err.Error())
		return nil, err
	}

	playtest, err := s.PlaytestRepository.PlaytestOfID(playtestID)
	if err != nil {
		s.Logger.Error(err.Error())
		return nil, err
	}

	if user == nil || playtest == nil {
		return nil, fmt.Errorf("you're not allowed to assign locations")
	}

	playtest.Finish()

	// And save
	err = s.PlaytestRepository.Save(playtest)
	if err != nil {
		s.Logger.Error(err.Error())
		return nil, err
	}

	return playtest, nil
}
