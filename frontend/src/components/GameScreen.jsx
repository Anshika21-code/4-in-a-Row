import Board from "./Board"

export default function GameScreen({ state, send }) {
  const { board, gameID, symbol, status, currentTurn, message } = state

  const isYourTurn = currentTurn === symbol

  return (
    <div>
      <h3>
        Game ID: {gameID}
        <button
          onClick={() => navigator.clipboard.writeText(gameID)}
          style={{ marginLeft: "10px" }}
        >
          Copy
        </button>
      </h3>

      <p>You are: {symbol}</p>

      {/* 🔴 DISCONNECT / RECONNECT BANNER */}
      {message && (
        <div className="banner">
          {message === "opponent_disconnected" && "⚠️ Opponent disconnected"}
          {message === "opponent_reconnected" && "✅ Opponent reconnected"}
          {message === "reconnected" && "🔄 You reconnected successfully"}
        </div>
      )}

      <p className={isYourTurn ? "your-turn" : "wait-turn"}>
        {status !== "continue"
          ? status
          : isYourTurn
          ? "Your Turn 🟢"
          : "Opponent's Turn ⏳"}
      </p>

      <Board
        board={board}
        onMove={(col) => send({ type: "move", column: col })}
        disabled={!isYourTurn || status !== "continue"}
      />
    </div>
  )
}
