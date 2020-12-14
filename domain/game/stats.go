package game

// Stats represent the basic information for a playing game
type Stats struct {
	MinPlayers        int `json:"min_players" gorm:"not null;default:1"`
	MaxPlayers        int `json:"max_players" gorm:"not null;default:5"`
	MinAge            int `json:"min_age" gorm:"not null;default:8"`
	EstimatedPlaytime int `json:"estimated_playtime" gorm:"not null;default:30"`
}
