import { useEffect, useRef } from "react";

export default function useWebSocket(onMessage) {
  const wsRef = useRef(null);
  const onMessageRef = useRef(onMessage);

  useEffect(() => {
    onMessageRef.current = onMessage;
  }, [onMessage]);

  useEffect(() => {
    const ws = new WebSocket("ws://localhost:8080/ws");
    wsRef.current = ws;

    ws.onopen = () => console.log("WebSocket connected");

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
