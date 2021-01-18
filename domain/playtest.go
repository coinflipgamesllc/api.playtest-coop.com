package domain

import (
	"database/sql"
	"time"

	"github.com/coinflipgamesllc/api.playtest-coop.com/domain/playtest"
	"gorm.io/gorm"
)

// Playtest is an individual test of a game
type Playtest struct {
	ID        uint           `json:"id" gorm:"primarykey" example:"123"`
	CreatedAt time.Time      `json:"created_at" example:"2020-12-11T15:29:49.321629-08:00"`
	UpdatedAt time.Time      `json:"updated_at" example:"2020-12-13T15:42:40.578904-08:00"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	Game   Game `json:"game"`
	GameID uint `json:"-"`

	Event         *Event    `json:"-"`
	EventID       *uint     `json:"-"`
	ScheduledDate time.Time `json:"-"`

	Requirements playtest.Requirements `json:"requirements" gorm:"embedded"`
	Location     *playtest.Location    `json:"location,omitempty" gorm:"embedded"`
	StartTime    sql.NullTime          `json:"start_time"`
	FeedbackTime sql.NullTime          `json:"feedback_time"`
	EndTime      sql.NullTime          `json:"end_time"`
	Players      []User                `json:"players" gorm:"many2many:playtesters;"`
}

// PlaytestRepository defines how to interact with playtests in database
type PlaytestRepository interface {
	PlaytestsOnDate(time.Time, uint) ([]Playtest, error)
	PlaytestOfID(id uint) (*Playtest, error)
	Save(*Playtest) error
}

// RegisterGame sets up a new playtest for a game at a specific time. It can optionally be tied to an event
func RegisterGame(game *Game, event *Event, sched time.Time, minPlayers, maxPlayers, duration uint, designerWantsToPlay bool, hopeToTest, ttsServer, ttsPassword string) *Playtest {
	sched = sched.Truncate(time.Hour * 24) // We only want the date

	return &Playtest{
		GameID:        game.ID,
		EventID:       &event.ID,
		ScheduledDate: sched,
		Requirements: playtest.Requirements{
			MinPlayers:          minPlayers,
			MaxPlayers:          maxPlayers,
			Duration:            duration,
			DesignerWantsToPlay: designerWantsToPlay,
			HopingToTest:        hopeToTest,
		},
		Location: &playtest.Location{
			TTSServer:   ttsServer,
			TTSPassword: ttsPassword,
		},
	}
}

// AssignTable will place the playtest at a table (real or virtual)
func (p *Playtest) AssignTable(table string) {
	if p.Location == nil {
		p.Location = &playtest.Location{
			Table: table,
		}
	} else {
		p.Location.Table = table
	}
}

// Start will set the time the playtest started to now
func (p *Playtest) Start() {
	p.StartTime = sql.NullTime{Time: time.Now(), Valid: true}
}

// StartFeedback will set the time feedback started to now
func (p *Playtest) StartFeedback() {
	p.FeedbackTime = sql.NullTime{Time: time.Now(), Valid: true}
}

// Finish will set the time the playtest ended to now
func (p *Playtest) Finish() {
	p.EndTime = sql.NullTime{Time: time.Now(), Valid: true}
}

// AddPlayer adds a new player to the test
func (p *Playtest) AddPlayer(player *User) {
	if player == nil {
		return
	}

	if p.Players == nil {
		p.Players = []User{}
	}

	for _, u := range p.Players {
		if u.ID == player.ID {
			return
		}
	}

	p.Players = append(p.Players, *player)
}

// RemovePlayer removes the specified player from the test
func (p *Playtest) RemovePlayer(player *User) {
	for i, u := range p.Players {
		if u.ID == player.ID {
			copy(p.Players[i:], p.Players[i+1:])
			p.Players = p.Players[:len(p.Players)-1]

			return
		}
	}
}
