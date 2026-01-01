package engine

import (
	"emitrr_assignment/backend/internal/bot"
	"emitrr_assignment/backend/internal/game"
)

func PlayTurn(g *game.Game, humanCol int) string {

	// Human move
	if !g.ApplyMove(humanCol) {
		return "invalid_move"
	}

	if g.Winner != 0 {
		return "human_won"
	}

	if g.IsDraw() {
		return "draw"
	}

	// Bot move
	botCol := bot.DecideMove(g, g.Player2.Symbol, g.Player1.Symbol)
	g.ApplyMove(botCol)

	if g.Winner != 0 {
		return "bot_won"
	}

	if g.IsDraw() {
		return "draw"
	}

	return "continue"
}
