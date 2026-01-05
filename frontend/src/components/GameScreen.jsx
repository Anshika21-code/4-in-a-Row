import Board from "./Board.jsx";

export default function GameScreen({ state, send }) {
  // 🛑 jab tak websocket state nahi aati
  if (!state || !state.game) {
    return <div>Loading game...</div>;
  }

  const { game, yourSymbol, gameID } = state;

  const isMyTurn =
    game.winner === "continue" &&
    yourSymbol === game.currentTurn;

  return (
    <div className="game-container">
      <h3>Game ID: {gameID}</h3>
      <p>You are: {yourSymbol}</p>

      <Board
        board={game.board}
        disabled={!isMyTurn}
        onMove={(col) =>
          send({
            type: "move",
            column: col,
          })
        }
      />

      {!isMyTurn && game.winner === "continue" && (
        <p>⏳ Waiting for opponent...</p>
      )}

      {game.winner !== "continue" && (
        <h2>{game.winner}</h2>
      )}
    </div>
  );
}
