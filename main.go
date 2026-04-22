package main

import (
	"log"
	"net/http"
)

func main() {
	http.Handle("/", http.FileServer(http.Dir("web")))

	if err := http.ListenAndServe(":7540", nil); err != nil {
		log.Printf("failed to start server: %v", err)
	}
}
