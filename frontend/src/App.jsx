import { useState } from "react";
import useWebSocket from "./hooks/useWebSocket";
import JoinScreen from "./components/JoinScreen";
import GameScreen from "./components/GameScreen";

export default function App() {
  const [state, setState] = useState(null);

  const wsRef = useWebSocket((data) => {
    console.log("WS MESSAGE RAW:", data);

    if (data.type === "state") {
      // 🔥 MAP backend → frontend format
      const mappedBoard = data.board.map(row =>
        row.map(cell => {
          if (cell === 88) return "X";
          if (cell === 79) return "O";
          return null;
        })
      );

      setState({
        gameID: data.gameId,
        yourSymbol: data.yourSymbol === 88 ? "X" : "O",
        game: {
          board: mappedBoard,
          winner: data.winner,
          currentTurn: data.yourSymbol === 88 ? "X" : "O",
        },
      });
    }
  });

  const send = (msg) => {
    if (wsRef.current?.readyState === WebSocket.OPEN) {
      wsRef.current.send(JSON.stringify(msg));
    }
  };

  const handleJoin = ({ username, mode }) => {
    send({ type: "join", username, mode });
  };

  if (!state) return <JoinScreen onJoin={handleJoin} />;

  return <GameScreen state={state} send={send} />;
}
