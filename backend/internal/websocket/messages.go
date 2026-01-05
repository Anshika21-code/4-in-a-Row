package websocket

type ClientMessage struct {
	Type     string `json:"type"`
	Username string `json:"username"`
	GameID   string `json:"gameId"`

	// Accept BOTH to avoid client/server mismatch.
	// Frontend should send: { type:"move", col: 3 }
	Col    int `json:"col"`
	Column int `json:"column"`

	Mode string `json:"mode"` // "human" | "bot"
}

type ServerMessage struct {
	Type string `json:"type"`

	GameID   string   `json:"gameId"`
	PlayerID string   `json:"playerId"`
	Board    [][]rune `json:"board"`

	YourSymbol  rune `json:"yourSymbol"`
	CurrentTurn rune `json:"currentTurn"` // whose turn on server now
	Winner      rune `json:"winner"`      // 0, 'X', 'O'

	// Status tells UI what happened in last action
	// "continue" | "invalid_move" | "human_won" | "bot_won" | "draw"
	Status string `json:"status"`

	// Debug message to show exactly what went wrong
	Message string `json:"message"`
}
