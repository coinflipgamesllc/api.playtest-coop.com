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

// NewGame creates a bare-bones game with a title and designer
func NewGame(title string, primaryDesigner User) *Game {
	return &Game{
		Title:     title,
		Designers: []User{primaryDesigner},
	}
}

// Rename will change the name of the game. Blank names are not allowed.
func (g *Game) Rename(newTitle string) {
	if newTitle != "" && g.Title != newTitle {
		g.Title = newTitle
	}
}
