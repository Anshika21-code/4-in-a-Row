import Cell from "./Cell";
import "../styles/index.css";

export default function Board({ board, onMove, disabled }) {
  if (!board || board.length === 0) {
    return <p>Waiting for game to start...</p>;
  }

  const rows = board.length;      // 6
  const cols = board[0].length;   // 7

  return (
    <div className="board">
      {Array.from({ length: cols }).map((_, col) => (
        <div
          key={col}
          className={`column ${!disabled ? "clickable" : ""}`}
          onClick={() => {
            if (disabled) return;
            onMove(col);
          }}
        >
          {Array.from({ length: rows }).map((_, row) => (
            <Cell
              key={row}
              value={board[row][col]}   // IMPORTANT: no row reversal
            />
          ))}
        </div>
      ))}
    </div>
  );
}
