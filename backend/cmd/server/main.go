package main

import (
	"log"
	"net/http"

	ws "emitrr_assignment/backend/internal/websocket"
)

func main() {
	http.HandleFunc("/ws", ws.HandleWS)

	log.Println(" WebSocket server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
