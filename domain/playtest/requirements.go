package playtest

// Requirements are the conditions that the designer(s) need for the playtest to succeed
type Requirements struct {
	MinPlayers          uint   `json:"min_players" example:"3"`
	MaxPlayers          uint   `json:"max_players" example:"5"`
	Duration            uint   `json:"duration" example:"60"`
	DesignerWantsToPlay bool   `json:"designer_wants_to_play" example:"true"`
	HopingToTest        string `json:"hoping_to_test" example:"Is the kerpluxic mechanic intuitive?"`
}
