package main

import (
	"os"
	"log"
	"net/http"

	ws "emitrr_assignment/backend/internal/websocket"
)

func main() {
	http.HandleFunc("/ws", ws.HandleWS)

	fs := http.FileServer(http.Dir("../frontend/dist"))
	http.Handle("/", fs)

	log.Println(" WebSocket server running on :8080")
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("Server running on :" + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
