package main

import (
	"final-project/pkg/server"
	"log"
	"net/http"
)

func main() {
	port := server.GetPort()
	http.Handle("/", http.FileServer(http.Dir("web")))

	if err := http.ListenAndServe(port, nil); err != nil {
		log.Printf("failed to start server: %v", err)
	}
}
