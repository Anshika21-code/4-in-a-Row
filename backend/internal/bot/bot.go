package bot

import "emitrr_assignment/backend/internal/game"

// DecideMove returns the column number the bot should play
func DecideMove(g *game.Game, botSymbol rune, playerSymbol rune) int {

	// 1 If bot can win in this move → WIN
	for col := 0; col < game.Columns; col++ {
		if canPlay(g, col) && isWinningMove(g, col, botSymbol) {
			return col
		}
	}

	// 2 If player can win next move → BLOCK
	for col := 0; col < game.Columns; col++ {
		if canPlay(g, col) && isWinningMove(g, col, playerSymbol) {
			return col
		}
	}

	// 3 Play center if possible
	center := game.Columns / 2
	if canPlay(g, center) {
		return center
	}

	// 4 Play first valid column
	for col := 0; col < game.Columns; col++ {
		if canPlay(g, col) {
			return col
		}
	}

	return -1 // should never happen
}

// --------------------
// Helper Functions
// --------------------

func canPlay(g *game.Game, column int) bool {
	return g.Board[column][0] == game.Empty
}

func isWinningMove(g *game.Game, column int, symbol rune) bool {
	clone := cloneGame(g)
	clone.CurrentTurn = symbol
	_, err := clone.DropDisc(column)
	return err == nil && clone.Status == game.Finished
}

func cloneGame(g *game.Game) *game.Game {
	boardCopy := make([][]rune, game.Columns)
	for c := 0; c < game.Columns; c++ {
		boardCopy[c] = make([]rune, game.Rows)
		copy(boardCopy[c], g.Board[c])
	}

	return &game.Game{
		ID:          g.ID,
		Board:       boardCopy,
		Players:     g.Players,
		CurrentTurn: g.CurrentTurn,
		Status:      g.Status,
		Winner:      nil,
	}
}
