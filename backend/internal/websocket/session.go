package websocket

import (
	"time"

	"emitrr_assignment/backend/internal/game"

	"github.com/gorilla/websocket"
)

type Session struct {
	GameID   string
	Game     *game.Game
	Players  map[string]*websocket.Conn
	LastSeen map[string]time.Time
}

var Sessions = make(map[string]*Session)
