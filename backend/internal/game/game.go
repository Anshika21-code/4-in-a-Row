package game

import "errors"

const (
	Rows    = 6
	Columns = 7
	Empty   = '.'
)

type GameStatus string

const (
	Waiting    GameStatus = "WAITING"
	InProgress GameStatus = "IN_PROGRESS"
	Finished   GameStatus = "FINISHED"
)

type Game struct {
	Board      [][]rune
	CurrentTurn rune
	Status     GameStatus
	Winner     rune
}

// --------------------
// Game Initialization
// --------------------

func NewGame() *Game {
	board := make([][]rune, Columns)
	for c := 0; c < Columns; c++ {
		board[c] = make([]rune, Rows)
		for r := 0; r < Rows; r++ {
			board[c][r] = Empty
		}
	}

	return &Game{
		Board:       board,
		CurrentTurn: 'X',
		Status:      InProgress,
		Winner:      Empty,
	}
}

// --------------------
// Drop Disc Logic
// --------------------

func (g *Game) DropDisc(column int) (int, error) {
	if g.Status != InProgress {
		return -1, errors.New("game is not active")
	}

	if column < 0 || column >= Columns {
		return -1, errors.New("invalid column")
	}

	for row := Rows - 1; row >= 0; row-- {
		if g.Board[column][row] == Empty {
			g.Board[column][row] = g.CurrentTurn

			if g.checkWin(column, row) {
				g.Status = Finished
				g.Winner = g.CurrentTurn
			} else if g.isDraw() {
				g.Status = Finished
				g.Winner = Empty
			} else {
				g.switchTurn()
			}

			return row, nil
		}
	}

	return -1, errors.New("column is full")
}

// --------------------
// Turn Handling
// --------------------

func (g *Game) switchTurn() {
	if g.CurrentTurn == 'X' {
		g.CurrentTurn = 'O'
	} else {
		g.CurrentTurn = 'X'
	}
}

// --------------------
// Win Detection
// --------------------

func (g *Game) checkWin(col, row int) bool {
	directions := [][2]int{
		{1, 0},  // Horizontal
		{0, 1},  // Vertical
		{1, 1},  // Diagonal \
		{1, -1}, // Diagonal /
	}

	for _, d := range directions {
		count := 1
		count += g.countDirection(col, row, d[0], d[1])
		count += g.countDirection(col, row, -d[0], -d[1])

		if count >= 4 {
			return true
		}
	}

	return false
}

func (g *Game) countDirection(col, row, dc, dr int) int {
	count := 0
	symbol := g.CurrentTurn

	c := col + dc
	r := row + dr

	for c >= 0 && c < Columns && r >= 0 && r < Rows {
		if g.Board[c][r] != symbol {
			break
		}
		count++
		c += dc
		r += dr
	}

	return count
}

// --------------------
// Draw Detection
// --------------------

func (g *Game) isDraw() bool {
	for c := 0; c < Columns; c++ {
		if g.Board[c][0] == Empty {
			return false
		}
	}
	return true
}
