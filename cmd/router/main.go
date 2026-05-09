package main

import (
	"distributed-cache/internal/router"
	"log"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "9000"
	}

	nodes := []string{
		"http://node1:8001",
		"http://node2:8002",
		"http://node3:8003",
	}

	server := router.NewServer(nodes)

	log.Println("Router running on port " + port)
	log.Fatal(http.ListenAndServe(":"+port, server.Routes()))
}
