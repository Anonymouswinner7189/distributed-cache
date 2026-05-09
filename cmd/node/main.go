package main

import (
	"distributed-cache/internal/node"
	"log"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8001"
	}

	server := node.NewServer()

	log.Println("Starting cache node on port: ", port)
	log.Fatal(http.ListenAndServe(":"+port, server.Routes()))
}
