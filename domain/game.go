package domain

import (
	"time"

	"github.com/coinflipgamesllc/api.playtest-coop.com/domain/game"
	"gorm.io/gorm"
)

// Game is the root structure containing all the information about a game
type Game struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	Title     string      `json:"title" gorm:"not null"`
	Overview  string      `json:"overview"`
	Status    game.Status `json:"status"`
	Designers []User      `json:"designers" gorm:"many2many:game_designers;"`
	Stats     game.Stats  `json:"stats" gorm:"embedded"`
}

// GameRepository defines how to interact with games in database
type GameRepository interface {
	ListGames(title, status, designer string, playerCount, age, playtime, limit, offset int, sort string) ([]Game, int, error)
	GameOfID(id uint) (*Game, error)
	Save(*Game) error
}

// NewGame creates a bare-bones game with a title and designer
func NewGame(title string, primaryDesigner User) *Game {
	return &Game{
		Title:     title,
		Status:    game.Prototype,
		Designers: []User{primaryDesigner},
		Stats: game.Stats{
			MinPlayers:        1,
			MaxPlayers:        5,
			MinAge:            8,
			EstimatedPlaytime: 30,
		},
	}
}

// MayBeUpdatedBy checks if the given user has permission to update the game.
// Currently, only designers may update games they own.
func (g *Game) MayBeUpdatedBy(user *User) bool {
	if user == nil {
		return false
	}

	for _, designer := range g.Designers {
		if designer.ID == user.ID {
			return true
		}
	}

	return false
}

// Rename will change the name of the game. Blank names are not allowed.
func (g *Game) Rename(newTitle string) {
	if newTitle != "" && g.Title != newTitle {
		g.Title = newTitle
	}
}

// UpdateOverview will change the overview for the game. Blank overviews are not allowed.
func (g *Game) UpdateOverview(newOverview string) {
	if newOverview != "" && g.Overview != newOverview {
		g.Overview = newOverview
	}
}

func (g *Game) UpdateStatus(ns string) error {
	newStatus, err := game.StatusFromString(ns)
	if err != nil {
		return err
	}

	g.Status = newStatus

	return nil
}

// AddDesigner will include the provider user as a designer on this game.
func (g *Game) AddDesigner(designer *User) {
	if designer == nil {
		return
	}

	if g.Designers == nil {
		g.Designers = []User{}
	}

	for _, d := range g.Designers {
		if d.ID == designer.ID {
			return
		}
	}

	g.Designers = append(g.Designers, *designer)
}

// ReplaceDesigners will overwrite the existing designer list with the newly provided one
func (g *Game) ReplaceDesigners(designers []User) {
	g.Designers = nil
	for _, designer := range designers {
		g.AddDesigner(&designer)
	}
}

// UpdateStats will replace the existing game stats with the provided values
func (g *Game) UpdateStats(minPlayers, maxPlayers, minAge, estimatedPlaytime int) {
	g.Stats = game.Stats{
		MinPlayers:        minPlayers,
		MaxPlayers:        maxPlayers,
		MinAge:            minAge,
		EstimatedPlaytime: estimatedPlaytime,
	}
}
