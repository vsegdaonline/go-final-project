package server

import (
	"net/http"
	"os"

	"go1f/pkg/api"

	"github.com/go-chi/chi/v5"
)

func getPort() string {
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = "7540"
	}
	return ":" + port
}

func Run() error {
	r := chi.NewRouter()
	r.Handle("/*", http.FileServer(http.Dir("web")))
	api.Init(r)
	return http.ListenAndServe(getPort(), r)
}
