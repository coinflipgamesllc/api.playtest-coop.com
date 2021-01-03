package game

import "time"

// RulesSection is a single section in a rule book
type RulesSection struct {
	ID        uint      `json:"id" gorm:"primarykey" example:"123"`
	CreatedAt time.Time `json:"created_at" example:"2020-12-11T15:29:49.321629-08:00"`
	UpdatedAt time.Time `json:"updated_at" example:"2020-12-13T15:42:40.578904-08:00"`

	GameID uint `json:"-"`

	Title   string `json:"title" gorm:"not null" example:"Components"`
	Content string `json:"content" example:"<ul><li>52 Cards</li><li>10 dice</li>..."`

	OrderBy uint `json:"order" example:"0"`
}

// NewRulesSection creates a new section attached to the provided game
func NewRulesSection(gameID uint, title, content string, order uint) *RulesSection {
	return &RulesSection{
		GameID:  gameID,
		Title:   title,
		Content: content,
		OrderBy: order,
	}
}

// UpdateTitle will replace the title
func (s *RulesSection) UpdateTitle(newTitle string) {
	s.Title = newTitle
}

// UpdateContent will replace the content
func (s *RulesSection) UpdateContent(newContent string) {
	s.Content = newContent
}

// UpdateOrder will simply accept the new order provided.
// Calling code is responsible for re-sorting the collection this section appears in.
func (s *RulesSection) UpdateOrder(order uint) {
	s.OrderBy = order
}
