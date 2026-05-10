package router

import (
	"distributed-cache/internal/consistenthash"
	"distributed-cache/internal/router/handlers"
	"net/http"
)

type Server struct {
	mux  *http.ServeMux
	ring *consistenthash.HashRing
}

func NewServer(nodes []string) *Server {
	ring := consistenthash.NewHashRing()
	for _, node := range nodes {
		ring.AddNode(node)
	}

	h := handlers.NewHandler(ring)
	mux := http.NewServeMux()
	mux.HandleFunc("/health", h.Health)
	mux.HandleFunc("/set", h.Set)
	mux.HandleFunc("/get", h.Get)
	mux.HandleFunc("/delete", h.Delete)

	return &Server{
		mux:  mux,
		ring: ring,
	}
}

func (s *Server) Routes() http.Handler {
	return s.mux
}
