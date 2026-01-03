import { useEffect, useRef } from "react"

export function useWebSocket(onMessage) {
  const ws = useRef(null)

  useEffect(() => {
    ws.current = new WebSocket("ws://localhost:8080/ws")

    ws.current.onmessage = (event) => {
      const data = JSON.parse(event.data)
      onMessage(data)
    }

    return () => ws.current.close()
  }, [])

  const send = (payload) => {
    if (ws.current?.readyState === 1) {
      ws.current.send(JSON.stringify(payload))
    }
  }

  return send
}
