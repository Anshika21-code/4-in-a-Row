package websocket

import (
	"log"
	"time"

	"emitrr_assignment/backend/internal/bot"
	"emitrr_assignment/backend/internal/engine"
	"emitrr_assignment/backend/internal/game"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

func handleConnection(conn *websocket.Conn) {
	defer conn.Close()

	var currentGameID string
	var currentPlayerID string

	log.Println(" New connection established")

	for {
		var msg ClientMessage
		if err := conn.ReadJSON(&msg); err != nil {
			log.Println(" Error reading message:", err)
			return
		}

		log.Printf("Received message type: %s", msg.Type)

		switch msg.Type {

		case "join":
			gameID, playerID, board, symbol := handleJoin(conn, msg)

			currentGameID = gameID
			currentPlayerID = playerID

			log.Printf(" Player joined - GameID: %s, PlayerID: %s, Mode: %s", gameID, playerID, msg.Mode)

			// Send initial state to the player
			sendState(
				conn,
				gameID,
				currentPlayerID,
				board,
				symbol,
				"continue",
			)

			// If bot mode, bot makes first move
			session := Sessions[gameID]
			if session != nil && session.Game.Player2 != nil && session.Game.Player2.ID == "BOT" {
				log.Println(" Bot making first move...")
				time.Sleep(500 * time.Millisecond) // Small delay for better UX
				
				col := bot.DecideMove(session.Game, 'O', 'X')
				result := engine.PlayTurn(session.Game, col)
				
				log.Printf(" Bot placed disc in column %d, result: %s", col, result)
				
				sendState(
					conn,
					gameID,
					currentPlayerID,
					session.Game.Board,
					symbol,
					result,
				)
			}

		case "move":
			log.Printf("🎮 Player move - Column: %d", msg.Column)
			
			session := Sessions[currentGameID]
			if session == nil {
				log.Println(" Session not found")
				return
			}

			// Player's move
			result := engine.PlayTurn(session.Game, msg.Column)
			log.Printf(" Player move result: %s", result)
			
			broadcastState(currentGameID, result)

			// Bot's turn (only if game is still continuing and it's bot mode)
			if session.Game.Player2 != nil && session.Game.Player2.ID == "BOT" && result == "continue" {
				log.Println("🤖 Bot's turn...")
				time.Sleep(500 * time.Millisecond) // Delay for realism
				
				col := bot.DecideMove(session.Game, 'O', 'X')
				botResult := engine.PlayTurn(session.Game, col)
				
				log.Printf("🤖 Bot placed disc in column %d, result: %s", col, botResult)
				
				broadcastState(currentGameID, botResult)
			}
		}
	}
}

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
	}

	session.Players[playerID] = conn
	session.LastSeen[playerID] = time.Now()
	Sessions[gameID] = session

	// Bot mode - add bot as Player2
	if msg.Mode == "bot" {
		log.Println("🤖 Adding bot to game")
		botPlayer := &game.Player{
			ID:       "BOT",
			Username: "Bot",
			Symbol:   'O',
		}
		g.Player2 = botPlayer
	}

	return gameID, playerID, g.Board, symbol
}

func broadcastState(gameID string, result string) {
	session := Sessions[gameID]
	if session == nil {
		return
	}

	for pid, conn := range session.Players {
		var symbol rune
		if session.Game.Player1 != nil && session.Game.Player1.ID == pid {
			symbol = session.Game.Player1.Symbol
		} else if session.Game.Player2 != nil && session.Game.Player2.ID == pid {
			symbol = session.Game.Player2.Symbol
		} else {
			symbol = 'X'
		}

		sendState(conn, gameID, pid, session.Game.Board, symbol, result)
	}
}

func sendState(
	conn *websocket.Conn,
	gameID string,
	playerID string,
	board [][]rune,
	symbol rune,
	result string,
) {
	state := ServerMessage{
		Type:       "state",
		GameID:     gameID,
		PlayerID:   playerID,
		Board:      board,
		YourSymbol: symbol,
		Winner:     result,
	}
	
	if err := conn.WriteJSON(state); err != nil {
		log.Printf(" Error sending state: %v", err)
	} else {
		log.Printf(" State sent to player %s", playerID)
	}
}
