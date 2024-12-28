package main

import (
	"fmt"
	"log"
	"net/http"

	InfoHandler "github.com/BlueSpadeXchain/blp-api/api/info"
	OrderHandler "github.com/BlueSpadeXchain/blp-api/api/orders"
	WebSocket "github.com/BlueSpadeXchain/blp-api/ws"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	http.HandleFunc("/api/order", OrderHandler.Handler)
	http.HandleFunc("/api/info", InfoHandler.Handler)
	http.HandleFunc("/ws", WebSocket.Handler)

	log.Println("Starting server on :8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
