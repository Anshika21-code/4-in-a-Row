package websocket

import (
	"log"
	"time"

	"emitrr_assignment/backend/internal/bot"
	"emitrr_assignment/backend/internal/game"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

func handleConnection(conn *websocket.Conn) {
	defer conn.Close()

	var currentGameID string
	var currentPlayerID string

	log.Println("New connection established")

	for {
		var msg ClientMessage
		if err := conn.ReadJSON(&msg); err != nil {
			log.Println("Error reading message:", err)
			return
		}

		log.Printf("Received message type: %s", msg.Type)

		switch msg.Type {

		// ================= JOIN =================
		case "join":
			gameID, playerID, _, _ := handleJoin(conn, msg)

			currentGameID = gameID
			currentPlayerID = playerID

			log.Printf("Player joined - GameID: %s, PlayerID: %s, Mode: %s", gameID, playerID, msg.Mode)

			session := Sessions[currentGameID]
			if session == nil {
				log.Println("JOIN ERROR: session not found after create")
				continue
			}

			sendState(conn, session, currentPlayerID, "continue", "joined_ok")

		// ================= MOVE =================
		case "move":
			if currentGameID == "" || currentPlayerID == "" {
				log.Println("MOVE ERROR: player not joined yet")
				continue
			}

			session := Sessions[currentGameID]
			if session == nil {
				log.Println("MOVE ERROR: session not found:", currentGameID)
				continue
			}

			// choose column from either field
			col := msg.Col
			if col == 0 && msg.Column != 0 {
				col = msg.Column
			}
			// Note: col==0 is valid column; so we must distinguish between "not provided"
			// vs "provided as 0". If both fields are 0, it's still valid (col=0).
			// We'll validate by range anyway.

			session.Mu.Lock()
			deferUnlock := true
			defer func() {
				if deferUnlock {
					session.Mu.Unlock()
				}
			}()

			// Validate player exists
			if session.Game.Player1 == nil {
				sendState(conn, session, currentPlayerID, "invalid_move", "server_error: player1 missing")
				continue
			}
			if session.Game.Player2 == nil {
				sendState(conn, session, currentPlayerID, "invalid_move", "server_error: player2 missing (bot not added?)")
				continue
			}

			// Game over
			if session.Game.Winner != 0 {
				sendState(conn, session, currentPlayerID, "invalid_move", "game_already_finished")
				continue
			}

			// Only Player1 (human) can move in bot mode
			if currentPlayerID != session.Game.Player1.ID {
				sendState(conn, session, currentPlayerID, "invalid_move", "not_allowed: only player1 can move vs bot")
				continue
			}

			// Bot busy guard (prevents spam clicks while bot is moving)
			if session.BotBusy {
				sendState(conn, session, currentPlayerID, "invalid_move", "bot_busy_wait")
				continue
			}

			// Turn check
			if session.Game.Turn != session.Game.Player1.Symbol {
				sendState(conn, session, currentPlayerID, "invalid_move", "not_your_turn")
				continue
			}

			// Column validity
			if col < 0 || col >= game.Columns {
				sendState(conn, session, currentPlayerID, "invalid_move", "invalid_column_out_of_range")
				continue
			}
			if !session.Game.CanPlay(col) {
				sendState(conn, session, currentPlayerID, "invalid_move", "invalid_column_full")
				continue
			}

			// Apply HUMAN move
			ok := session.Game.ApplyMove(col)
			if !ok {
				sendState(conn, session, currentPlayerID, "invalid_move", "apply_move_failed")
				continue
			}

			// Broadcast after human move
			if session.Game.Winner != 0 {
				deferUnlock = false
				session.Mu.Unlock()
				broadcastState(session, "human_won", "human_won")
				continue
			}
			if session.Game.IsDraw() {
				deferUnlock = false
				session.Mu.Unlock()
				broadcastState(session, "draw", "draw")
				continue
			}

			// Now BOT turn
			session.BotBusy = true
			deferUnlock = false
			session.Mu.Unlock()

			// Optional small delay for UX
			time.Sleep(500 * time.Millisecond)

			session.Mu.Lock()
			// decide bot column using actual symbols
			botCol := bot.DecideMove(session.Game, session.Game.Player2.Symbol, session.Game.Player1.Symbol)
			if botCol < 0 {
				session.BotBusy = false
				session.Mu.Unlock()
				broadcastState(session, "draw", "bot_no_moves")
				continue
			}

			// Ensure it's bot's turn (should be, after human move)
			if session.Game.Turn != session.Game.Player2.Symbol {
				session.BotBusy = false
				session.Mu.Unlock()
				broadcastState(session, "invalid_move", "server_error: expected bot turn but turn mismatch")
				continue
			}

			_ = session.Game.ApplyMove(botCol)

			session.BotBusy = false

			// Compute result
			status := "continue"
			msgText := "continue"
			if session.Game.Winner != 0 {
				status = "bot_won"
				msgText = "bot_won"
			} else if session.Game.IsDraw() {
				status = "draw"
				msgText = "draw"
			}

			session.Mu.Unlock()
			broadcastState(session, status, msgText)

		default:
			log.Printf("Unknown message type: %s", msg.Type)
		}
	}
}

// ================= JOIN HANDLER =================

func handleJoin(conn *websocket.Conn, msg ClientMessage) (string, string, [][]rune, rune) {
	gameID := msg.GameID
	if gameID == "" {
		gameID = uuid.New().String()
	}

	playerID := uuid.New().String()
	symbol := 'X'

	player := &game.Player{
		ID:       playerID,
		Username: msg.Username,
		Symbol:   symbol,
	}

	g := game.NewGame(gameID, player)

	session := &Session{
		GameID:   gameID,
		Game:     g,
		Players:  make(map[string]*websocket.Conn),
		LastSeen: make(map[string]time.Time),
		BotBusy:  false,
	}

	session.Players[playerID] = conn
	session.LastSeen[playerID] = time.Now()
	Sessions[gameID] = session

	// Bot mode
	if msg.Mode == "bot" {
		log.Println("🤖 Adding bot to game")
		botPlayer := &game.Player{
			ID:       "BOT",
			Username: "Bot",
			Symbol:   'O',
		}
		g.AddSecondPlayer(botPlayer)
	}

	return gameID, playerID, g.Board, symbol
}

// ================= STATE BROADCAST =================

func broadcastState(session *Session, status string, message string) {
	for pid, conn := range session.Players {
		sendState(conn, session, pid, status, message)
	}
}

// ================= SEND STATE =================

func sendState(conn *websocket.Conn, session *Session, playerID string, status string, message string) {
	var symbol rune = 'X'
	if session.Game.Player1 != nil && session.Game.Player1.ID == playerID {
		symbol = session.Game.Player1.Symbol
	} else if session.Game.Player2 != nil && session.Game.Player2.ID == playerID {
		symbol = session.Game.Player2.Symbol
	}

	state := ServerMessage{
		Type:        "state",
		GameID:      session.GameID,
		PlayerID:    playerID,
		Board:       session.Game.Board,
		YourSymbol:  symbol,
		CurrentTurn: session.Game.Turn,
		Winner:      session.Game.Winner,
		Status:      status,
		Message:     message,
	}

	if err := conn.WriteJSON(state); err != nil {
		log.Printf("Error sending state to player %s: %v", playerID, err)
	} else {
		log.Printf("State sent to player %s (status=%s message=%s)", playerID, status, message)
	}
}
