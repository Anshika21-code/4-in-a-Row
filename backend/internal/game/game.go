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
	if !g.CanPlay(col) || g.Winner != 0 {
		return false
	}

	for row := Rows - 1; row >= 0; row-- {
		if g.Board[row][col] == 0 {
			g.Board[row][col] = g.Turn
			break
		}
	}

	if g.checkWinner(g.Turn) {
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

func (g *Game) checkWinner(sym rune) bool {
	dirs := [][]int{
		{0, 1}, {1, 0}, {1, 1}, {1, -1},
	}

	for r := 0; r < Rows; r++ {
		for c := 0; c < Columns; c++ {
			if g.Board[r][c] != sym {
				continue
			}
			for _, d := range dirs {
				count := 1
				for k := 1; k < 4; k++ {
					nr := r + d[0]*k
					nc := c + d[1]*k
					if nr < 0 || nr >= Rows || nc < 0 || nc >= Columns {
						break
					}
					if g.Board[nr][nc] != sym {
						break
					}
					count++
				}
				if count == 4 {
					return true
				}
			}
		}
	}
	return false
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
