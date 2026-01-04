import Cell from "./Cell"

export default function Board({ board, onMove, disabled }) {
  if (!board || board.length === 0) {
    return <p>Waiting for game to start...</p>
  }

  const rows = board.length
  const cols = board[0].length

  return (
    <div className="board">
      {Array.from({ length: cols }).map((_, colIndex) => (
        <div
          key={colIndex}
          className="column"
          onClick={() => !disabled && onMove(colIndex)}
        >
          {Array.from({ length: rows }).map((_, rowIndex) => (
            <Cell
              key={rowIndex}
              value={board[rowIndex][colIndex]}
            />
          ))}
        </div>
      ))}
    </div>
  )
}
