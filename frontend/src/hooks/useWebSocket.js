import { useEffect, useRef } from "react";

export default function useWebSocket(onMessage) {
  const wsRef = useRef(null);

  useEffect(() => {
    if (wsRef.current) return;

    const ws = new WebSocket("ws://localhost:8080/ws");
    wsRef.current = ws;

    ws.onopen = () => console.log("WS connected");
    ws.onmessage = (e) => onMessage(JSON.parse(e.data));
    ws.onclose = () => console.log("WS closed");
    ws.onerror = (e) => console.error("WS error", e);

    //  DO NOT CLOSE in StrictMode dev
    return () => {};
  }, []);

  return wsRef;
}
