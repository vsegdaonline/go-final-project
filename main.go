package main

import (
	"go1f/pkg/db"
	"go1f/pkg/server"
	"log"
)

func main() {
	if err := db.Init(db.GetDBFile()); err != nil {
		log.Fatal(err.Error())
	}
	if err := server.Run(); err != nil {
		log.Printf("failed to start server: %v", err)
	}
}
