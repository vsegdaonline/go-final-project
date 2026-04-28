package server

import (
	"net/http"
	"os"

	"go1f/pkg/api"
)

func getPort() string {
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = "7540"
	}
	return ":" + port
}

func Run() error {
	http.Handle("/", http.FileServer(http.Dir("web")))
	api.Init()
	return http.ListenAndServe(getPort(), nil)
}
