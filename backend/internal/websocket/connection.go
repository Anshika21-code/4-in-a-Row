package websocket

import (
	"time"

	"emitrr_assignment/backend/internal/bot"
	"emitrr_assignment/backend/internal/engine"
	"emitrr_assignment/backend/internal/game"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)


func handleConnection(conn *websocket.Conn) {
	defer conn.Close()

	var currentGameID, currentPlayerID string

	for {
		var msg ClientMessage
		if err := conn.ReadJSON(&msg); err != nil {
			if currentGameID != "" && currentPlayerID != "" {
				markDisconnected(currentGameID, currentPlayerID)
			}
			return
		}
		

		switch msg.Type {

		case "join":
	        gameID, playerID, board, symbol := handleJoin(conn, msg)

	        currentGameID = gameID
	        currentPlayerID = playerID

	        sendState(conn, gameID, playerID, board, symbol, "")



		case "move":
			session := Sessions[currentGameID]
			result := engine.PlayTurn(session.Game, msg.Column)

			broadcastState(currentGameID, result)

			// bot turn
			if result == "continue" {
				col := bot.DecideMove(session.Game, 'O', 'X')
				engine.PlayTurn(session.Game, col)
				broadcastState(currentGameID, "continue")
			}
		}
	}
}

func handleJoin(conn *websocket.Conn, msg ClientMessage) (string, string, [][]rune, rune) {
	SessionsMu.Lock()
	defer SessionsMu.Unlock()

	// --------------------
	// 1 RECONNECT FLOW
	// --------------------
	if msg.GameID != "" && msg.PlayerID != "" {
		session, ok := Sessions[msg.GameID]
		if ok {
			if last, ok := session.LastSeen[msg.PlayerID]; ok {
				if time.Since(last) <= 30*time.Second {
					session.Players[msg.PlayerID] = conn
					delete(session.LastSeen, msg.PlayerID)

					var symbol rune
					if session.Game.Player1.ID == msg.PlayerID {
						symbol = session.Game.Player1.Symbol
					} else {
						symbol = session.Game.Player2.Symbol
					}

					//  RETURN WITH BOARD
					return msg.GameID, msg.PlayerID, session.Game.Board, symbol
				}
			}
		}
	}

	// --------------------
	// 2 CREATE NEW ROOM (DEFAULT)
	// --------------------
	if msg.GameID == "" {
		gameID := uuid.NewString()
		playerID := uuid.NewString()

		player := &game.Player{
			ID:       playerID,
			Username: msg.Username,
			Symbol:   'X',
		}

		g := game.NewGame(gameID, player)

		session := &Session{
			Game:     g,
			Players:  map[string]*websocket.Conn{playerID: conn},
			LastSeen: map[string]time.Time{},
		}

		Sessions[gameID] = session

		// bot fallback
		go botFallback(gameID)

		//  RETURN WITH BOARD
		return gameID, playerID, g.Board, 'X'
	}

	// --------------------
	// 3 JOIN EXISTING ROOM
	// --------------------
	session := Sessions[msg.GameID]
	playerID := uuid.NewString()

	session.Players[playerID] = conn
	session.Game.AddSecondPlayer(&game.Player{
		ID:       playerID,
		Username: msg.Username,
		Symbol:   'O',
	})

	//  RETURN WITH BOARD
	return msg.GameID, playerID, session.Game.Board, 'O'
}




func botFallback(gameID string) {
    time.Sleep(10 * time.Second)

    SessionsMu.Lock()
    defer SessionsMu.Unlock()

    session := Sessions[gameID]
    if len(session.Players) == 1 {
        session.Game.AddSecondPlayer(&game.Player{
            ID: "bot", Username: "Bot", Symbol: 'O',
        })
        broadcastState(gameID, "bot_joined")
    }
}

func markDisconnected(gameID, playerID string) {
    session := Sessions[gameID]
    session.Mutex.Lock()
    defer session.Mutex.Unlock()

    session.LastSeen[playerID] = time.Now()
}

func broadcastState(gameID, result string) {
	session := Sessions[gameID]

	for pid, conn := range session.Players {
		var symbol rune
		if session.Game.Player1.ID == pid {
			symbol = session.Game.Player1.Symbol
		} else {
			symbol = session.Game.Player2.Symbol
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
	_ = conn.WriteJSON(state)
}




