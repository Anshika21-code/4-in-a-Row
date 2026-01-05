package game

const (
	Rows    = 6
	Columns = 7
)

type Player struct {
	ID       string
	Username string
	Symbol   rune // 'X' or 'O'
}

type Game struct {
	Board   [][]rune
	Player1 *Player
	Player2 *Player
	Turn    rune
	Winner  rune // 'X', 'O', or 0 (no winner yet)
}

func NewGame(id string, p1 *Player) *Game {
	board := make([][]rune, Rows)
	for i := range board {
		board[i] = make([]rune, Columns)
	}

	return &Game{
		Board:   board,
		Player1: p1,
		Turn:    p1.Symbol,
		Winner:  0,
	}
}

func (g *Game) AddSecondPlayer(p2 *Player) {
	g.Player2 = p2
}

func (g *Game) CanPlay(col int) bool {
	return col >= 0 && col < Columns && g.Board[0][col] == 0
}

func (g *Game) ApplyMove(col int) bool {
	if g.Winner != 0 || !g.CanPlay(col) {
		return false
	}
	if g.Player1 == nil || g.Player2 == nil {
		return false
	}

	placedRow := -1
	for row := Rows - 1; row >= 0; row-- {
		if g.Board[row][col] == 0 {
			g.Board[row][col] = g.Turn
			placedRow = row
			break
		}
	}
	if placedRow == -1 {
		return false
	}

	if g.checkWinnerFrom(placedRow, col, g.Turn) {
		g.Winner = g.Turn
		return true
	}

	// switch turn
	if g.Turn == g.Player1.Symbol {
		g.Turn = g.Player2.Symbol
	} else {
		g.Turn = g.Player1.Symbol
	}

	return true
}

func (g *Game) checkWinnerFrom(r, c int, sym rune) bool {
	if r < 0 || r >= Rows || c < 0 || c >= Columns {
		return false
	}
	if g.Board[r][c] != sym {
		return false
	}

	dirs := [][2]int{{0, 1}, {1, 0}, {1, 1}, {1, -1}}
	for _, d := range dirs {
		total := 1
		total += g.countDir(r, c, d[0], d[1], sym)
		total += g.countDir(r, c, -d[0], -d[1], sym)
		if total >= 4 {
			return true
		}
	}
	return false
}

func (g *Game) countDir(r, c, dr, dc int, sym rune) int {
	count := 0
	r += dr
	c += dc
	for r >= 0 && r < Rows && c >= 0 && c < Columns && g.Board[r][c] == sym {
		count++
		r += dr
		c += dc
	}
	return count
}

func (g *Game) Clone() *Game {
	board := make([][]rune, Rows)
	for i := range board {
		board[i] = make([]rune, Columns)
		copy(board[i], g.Board[i])
	}

	return &Game{
		Board:   board,
		Player1: g.Player1,
		Player2: g.Player2,
		Turn:    g.Turn,
		Winner:  g.Winner,
	}
}

func (g *Game) IsDraw() bool {
	for col := 0; col < Columns; col++ {
		if g.CanPlay(col) {
			return false
		}
	}
	return g.Winner == 0
}
