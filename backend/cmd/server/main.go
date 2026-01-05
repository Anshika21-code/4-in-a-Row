package main

import (
	"log"
	"net/http"

	ws "emitrr_assignment/backend/internal/websocket"
)

func main() {
	http.HandleFunc("/ws", ws.HandleWS)

	fs := http.FileServer(http.Dir("../frontend/dist"))
	http.Handle("/", fs)

	log.Println(" WebSocket server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
