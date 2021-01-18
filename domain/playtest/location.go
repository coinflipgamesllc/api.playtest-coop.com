package playtest

// Location describes where a playtest is taking place
type Location struct {
	Table       string `json:"table,omitempty"`
	TTSServer   string `json:"tts_server,omitempty"`
	TTSPassword string `json:"tts_password,omitempty"`
}
