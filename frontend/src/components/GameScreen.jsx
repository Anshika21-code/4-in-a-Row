import Board from "./Board"

export default function GameScreen({ state, send }) {
  if (!state || !state.board || state.board.length === 0) {
    return <p>Waiting for game to start...</p>
  }

  // ✅ FIXED TURN LOGIC
  const isYourTurn =
    state.winner === "continue" && state.yourSymbol === 88

  const handleMove = (col) => {
    if (!isYourTurn) return

    send({
      type: "move",
      column: col
    })
  }

  return (
    <div>
      <h3>Game ID: {state.gameID}</h3>
      <p>You are: {state.yourSymbol === 88 ? "X" : "O"}</p>

      <Board
        board={state.board}
        onMove={handleMove}
        disabled={!isYourTurn}
      />

      {state.winner && state.winner !== "continue" && (
        <h2>{state.winner}</h2>
      )}
    </div>
  )
}
