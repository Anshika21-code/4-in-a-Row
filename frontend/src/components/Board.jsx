import Cell from "./Cell"
import "../styles/index.css"

export default function Board({ board, onMove, disabled }) {
  if (!board || board.length === 0) {
    return <p>Waiting for game to start...</p>
  }

  const rows = board.length
  const cols = board[0].length

  return (
    <div className="board">
      {Array.from({ length: cols }).map((_, col) => (
        <div
          key={col}
          className={`column ${!disabled ? "clickable" : ""}`}
          onClick={() => !disabled && onMove(col)}
        >
          {Array.from({ length: rows }).map((_, row) => (
            <Cell
              key={row}
              value={board[rows - 1 - row][col]}  
            />
          ))}
        </div>
      ))}
    </div>
  )
}
