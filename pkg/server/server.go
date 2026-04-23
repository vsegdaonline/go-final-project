package server

import "os"

func GetPort() string {
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = "7540"
	}
	return ":" + port
}
