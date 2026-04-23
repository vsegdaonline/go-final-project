package main

import (
	"go1f/pkg/server"
	"log"
)

func main() {
	if err := server.Run(); err != nil {
		log.Printf("failed to start server: %v", err)
	}
}
