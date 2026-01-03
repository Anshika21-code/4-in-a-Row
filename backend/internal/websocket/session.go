package websocket

import (
	"sync"
	"time"

	"emitrr_assignment/backend/internal/game"
	"github.com/gorilla/websocket"
)

type Session struct {
	Game     *game.Game
	Players  map[string]*websocket.Conn // playerId -> conn
	LastSeen map[string]time.Time
	Mutex    sync.Mutex
}

var Sessions = map[string]*Session{}
var SessionsMu sync.Mutex
