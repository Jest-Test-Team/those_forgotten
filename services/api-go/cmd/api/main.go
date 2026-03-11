package main

import (
	"log"

	"github.com/dennislee928/those_forgotten/services/api-go/internal/app"
)

func main() {
	server, err := app.NewServer()
	if err != nil {
		log.Fatalf("bootstrap api: %v", err)
	}

	if err := server.Start(); err != nil {
		log.Fatalf("start api: %v", err)
	}
}
