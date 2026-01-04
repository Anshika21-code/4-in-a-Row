package websocket

type ClientMessage struct {
	Type     string `json:"type"`
	Username string `json:"username"`
	GameID   string `json:"gameId"`
	Column   int    `json:"column"`
	Mode     string `json:"mode"` // "human" | "bot"
}

type ServerMessage struct {
	Type       string   `json:"type"`
	GameID     string   `json:"gameId"`
	PlayerID   string   `json:"playerId"`
	Board      [][]rune `json:"board"`
	YourSymbol rune     `json:"yourSymbol"`
	Winner     string   `json:"winner"`
}
