package models

type Server struct {
	Addr       string   `json:"addr"`
	GamePort   int      `json:"gameport"`
	Name       string   `json:"name"`
	Map        string   `json:"map"`
	Players    int      `json:"players"`
	MaxPlayers int      `json:"max_players"`
	Region     int      `json:"region"`
	Tags       []string `json:"tags"`
	// Add other fields as needed
}
