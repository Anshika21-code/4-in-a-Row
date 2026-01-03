import Cell from "./Cell"

export default function Board({ board, onMove, disabled }) {
  return (
    <div className="board">
      {board[0].map((_, col) => (
        <div
          key={col}
          className="column"
          onClick={() => !disabled && onMove(col)}
        >
          {board.map((row, r) => (
            <Cell key={r} value={board[r][col]} />
          ))}
        </div>
      ))}
    </div>
  )
}
