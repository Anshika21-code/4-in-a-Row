package websocket

type ClientMessage struct {
	Type     string `json:"type"` // join | move | reconnect
	GameID   string `json:"gameId,omitempty"`
	PlayerID string `json:"playerId,omitempty"`
	Username string `json:"username,omitempty"`
	Column   int    `json:"column,omitempty"`
}

type ServerMessage struct {
	Type       string        `json:"type"`
	GameID     string        `json:"gameId,omitempty"`
	PlayerID   string        `json:"playerId,omitempty"`
	Board      [][]rune      `json:"board,omitempty"`
	YourSymbol rune          `json:"yourSymbol,omitempty"`
	Winner     string        `json:"winner,omitempty"`
}
