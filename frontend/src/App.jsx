import { useState, useEffect } from "react";
import JoinScreen from "./components/JoinScreen";
import GameScreen from "./components/GameScreen";
import useWebSocket from "./hooks/useWebSocket";

export default function App() {
  const [gameState, setGameState] = useState(null);

  // ✅ single message handler
  const onMessage = (data) => {
    setGameState(data);

    if (data.gameID && data.playerID) {
      localStorage.setItem("gameID", data.gameID);
      localStorage.setItem("playerID", data.playerID);
    }

    if (data.message) {
      setTimeout(() => {
        setGameState((prev) => ({ ...prev, message: null }));
      }, 3000);
    }
  };

  // ✅ ONLY ONE websocket
  const wsRef = useWebSocket(onMessage);

  // ✅ JOIN
  const handleJoin = (payload) => {
    if (!wsRef.current || wsRef.current.readyState !== WebSocket.OPEN) {
      console.log("WS not ready");
      return;
    }

    wsRef.current.send(JSON.stringify(payload));
    console.log("✅ join sent", payload);
  };

  // ❌ Remove reconnect for now (backend doesn't handle it yet)
  // useEffect(() => {}, []);

  if (!gameState) {
    return <JoinScreen onJoin={handleJoin} />;
  }

  return <GameScreen state={gameState} send={handleJoin} />;
}
