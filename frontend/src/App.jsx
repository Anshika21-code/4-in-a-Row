import { useCallback, useMemo, useState } from "react";
import useWebSocket from "./hooks/useWebSocket";
import JoinScreen from "./components/JoinScreen";
import GameScreen from "./components/GameScreen";

function mapCell(cell) {
  // Backend sends rune codes; 88='X', 79='O'
  if (cell === 88) return "X";
  if (cell === 79) return "O";
  return null;
}

function mapWinner(w) {
  if (w === 88) return "X";
  if (w === 79) return "O";
  return null;
}

export default function App() {
  const [state, setState] = useState(null);

  const onMessage = useCallback((data) => {
    console.log("WS MESSAGE RAW:", data);

    if (!data || typeof data !== "object") {
      console.error("WS ERROR: message is not an object:", data);
      return;
    }

    if (data.type !== "state") return;

    if (!Array.isArray(data.board)) {
      console.error("WS ERROR: state.board missing/invalid:", data);
      return;
    }

    const mappedBoard = data.board.map((row, r) => {
      if (!Array.isArray(row)) {
        console.error(`WS ERROR: state.board[${r}] not array`, row);
        return [];
      }
      return row.map((cell) => mapCell(cell));
    });

    setState({
      gameID: data.gameId ?? null,
      yourSymbol: data.yourSymbol === 88 ? "X" : "O",
      game: {
        board: mappedBoard,
        winner: mapWinner(data.winner),
        // Expect backend to send currentTurn as 88/79 or 0; otherwise null
        currentTurn:
          data.currentTurn === 88 || data.currentTurn === 79
            ? mapCell(data.currentTurn)
            : null,
        // Optional message from backend (recommended)
        status: data.status ?? null,
      },
    });
  }, []);

  const wsRef = useWebSocket(onMessage);

  const send = useCallback(
    (msg) => {
      const ws = wsRef.current;
      if (!ws) {
        console.error("WS SEND FAILED: socket not created yet.");
        return;
      }
      if (ws.readyState !== WebSocket.OPEN) {
        console.error("WS SEND FAILED: socket not open. readyState =", ws.readyState);
        return;
      }
      ws.send(JSON.stringify(msg));
    },
    [wsRef]
  );

  const api = useMemo(() => ({ send }), [send]);

  const handleJoin = useCallback(
    ({ username, mode }) => {
      send({ type: "join", username, mode });
    },
    [send]
  );

  if (!state) return <JoinScreen onJoin={handleJoin} />;

  return <GameScreen state={state} send={api.send} />;
}
