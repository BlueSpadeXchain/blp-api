package main

import (
	"fmt"
	"log"
	"net/http"

	InfoHandler "blp-api/api/info.go"
	OrderHandler "blp-api/api/order.go"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	http.HandleFunc("/api/order", OrderHandler.Handler)
	http.HandleFunc("/api/info", InfoHandler.Handler)

	log.Println("Starting server on :8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
