package bot

import "emitrr_assignment/backend/internal/game"

// DecideMove → bot wins if possible, else blocks, else center, else first valid
func DecideMove(g *game.Game, botSymbol rune, playerSymbol rune) int {

	// 1 Try to WIN
	for col := 0; col < game.Columns; col++ {
		if g.CanPlay(col) {
			clone := g.Clone()
			clone.ApplyMove(col)
			if clone.Winner == botSymbol {
				return col
			}
		}
	}

	// 2 BLOCK opponent win
	for col := 0; col < game.Columns; col++ {
		if g.CanPlay(col) {
			clone := g.Clone()
			clone.Turn = playerSymbol
			clone.ApplyMove(col)
			if clone.Winner == playerSymbol {
				return col
			}
		}
	}

	// 3 Play center
	center := game.Columns / 2
	if g.CanPlay(center) {
		return center
	}

	// 4 First valid column
	for col := 0; col < game.Columns; col++ {
		if g.CanPlay(col) {
			return col
		}
	}

	return -1
}
