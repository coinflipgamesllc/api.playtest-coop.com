package domain

import (
	"time"

	"github.com/coinflipgamesllc/api.playtest-coop.com/domain/game"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

// Game is the root structure containing all the information about a game
type Game struct {
	ID        uint           `json:"id" gorm:"primarykey" example:"123"`
	CreatedAt time.Time      `json:"created_at" example:"2020-12-11T15:29:49.321629-08:00"`
	UpdatedAt time.Time      `json:"updated_at" example:"2020-12-13T15:42:40.578904-08:00"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	Title     string              `json:"title" gorm:"not null" example:"The Best Game"`
	Overview  string              `json:"overview" example:"In the Best Game, players take on the role of ..."`
	Status    game.Status         `json:"status" example:"Prototype"`
	Stats     game.Stats          `json:"stats" gorm:"embedded"`
	Mechanics pq.StringArray      `json:"mechanics" gorm:"type:text[]" example:"['Hidden Movement', 'Worker Placement']"`
	Designers []User              `json:"designers" gorm:"many2many:game_designers;"`
	Files     []File              `json:"files"`
	Rules     []game.RulesSection `json:"-"`

	TabletopSimulatorMod int `json:"tts_mod" example:"2247242964"`
}

// GameRepository defines how to interact with games in database
type GameRepository interface {
	ListGames(title, status, designer string, owner uint, playerCount, age, playtime, limit, offset int, sort string) ([]Game, int, error)
	GameOfID(id uint) (*Game, error)
	RulesOfGame(id uint) ([]game.RulesSection, error)
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

// UpdateStatus will set the status of the game, provided the status exists
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
	if minPlayers != 0 {
		g.Stats.MinPlayers = minPlayers
	}
	if maxPlayers != 0 {
		g.Stats.MaxPlayers = maxPlayers
	}
	if minAge != 0 {
		g.Stats.MinAge = minAge
	}
	if estimatedPlaytime != 0 {
		g.Stats.EstimatedPlaytime = estimatedPlaytime
	}
}

// ReplaceMechanics will overwite the existing mechanics list with the new one
func (g *Game) ReplaceMechanics(mechanics []string) {
	g.Mechanics = nil
	for _, mechanic := range mechanics {
		g.Mechanics = append(g.Mechanics, mechanic)
	}
}

// LinkTabletopSimulatorMod will link the specified mod to this game
func (g *Game) LinkTabletopSimulatorMod(mod int) {
	g.TabletopSimulatorMod = mod
}
