import Board from "./Board.jsx";

export default function GameScreen({ state, send }) {
  if (!state || !state.game) {
    return <div>Loading game...</div>;
  }

  const { game, yourSymbol, gameID } = state;

  // winner is null until someone wins (X/O)
  const gameOver = game.winner === "X" || game.winner === "O";

  // currentTurn might be null if backend doesn't send it yet
  const turnKnown = game.currentTurn === "X" || game.currentTurn === "O";

  const isMyTurn = !gameOver && turnKnown && yourSymbol === game.currentTurn;

  const statusText = (() => {
    if (!turnKnown) return "Waiting for server turn info...";
    if (gameOver) return `Winner: ${game.winner}`;
    if (isMyTurn) return "Your turn";
    return "Waiting for opponent/bot...";
  })();

  return (
    <div className="game-container">
      <h3>Game ID: {gameID}</h3>
      <p>You are: {yourSymbol}</p>
      <p>{statusText}</p>

      <Board
        board={game.board}
        disabled={!isMyTurn}
        onMove={(col) => {
          if (!turnKnown) {
            console.error("MOVE BLOCKED: currentTurn not received from server yet.");
            return;
          }
          if (gameOver) {
            console.error("MOVE BLOCKED: game is already over.");
            return;
          }
          if (!isMyTurn) {
            console.error("MOVE BLOCKED: not your turn.");
            return;
          }

          console.log("SENDING MOVE col =", col);

          // IMPORTANT: send `col`, not `column`
          send({
            type: "move",
            col: col,
          });
        }}
      />

      {gameOver && <h2>Winner: {game.winner}</h2>}
    </div>
  );
}
