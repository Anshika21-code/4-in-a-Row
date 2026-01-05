import { useEffect, useRef } from "react";

export default function useWebSocket(onMessage) {
  const wsRef = useRef(null);
  const onMessageRef = useRef(onMessage);

  useEffect(() => {
    onMessageRef.current = onMessage;
  }, [onMessage]);

  useEffect(() => {
    // Use environment variable for WebSocket URL
    // In production, this will use the Render backend URL
    // In development, it falls back to localhost
    const wsUrl = import.meta.env.VITE_WS_URL || "ws://localhost:8080";
    const ws = new WebSocket(`${wsUrl}/ws`);
    wsRef.current = ws;

    ws.onopen = () => console.log("WebSocket connected to:", wsUrl);

    ws.onmessage = (event) => {
      const raw = event.data;
      // Helpful when backend accidentally sends plain text
      // or sends numbers/bytes that aren't JSON.
      try {
        const data = JSON.parse(raw);
        onMessageRef.current?.(data);
      } catch (e) {
        console.error("WS PARSE ERROR: message is not valid JSON:", raw);
      }
    };

    ws.onerror = (e) => console.error("WebSocket error", e);
    ws.onclose = (e) => console.log("WebSocket closed", e.code, e.reason);

    return () => {
      try {
        ws.close();
      } catch (_) {}
    };
  }, []);

  return wsRef;
}