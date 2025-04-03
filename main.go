package main

import (
	"log"
	"net/http"

	"github.com/CodexRodney/WeStreamBackend/internal/rooms"
)

func main() {
	setupAPI()
	// Serve on port :8080, fudge yeah hardcoded port
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// setupAPI will start all Routes and their Handlers
func setupAPI() {
	manager := rooms.NewManager()

	// Serve the ./public directory at Route
	http.Handle("/ws", http.HandlerFunc(manager.ServeWS))
}
