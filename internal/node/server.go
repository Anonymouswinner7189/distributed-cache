package node

import (
	"distributed-cache/internal/cache"
	"distributed-cache/internal/node/handlers"
	"net/http"
	"time"
)

type Server struct {
	mux *http.ServeMux
}

func NewServer() *Server {
	c := cache.NewCache()
	c.StartCleanup(10 * time.Second)

	h := handlers.NewHandler(c)

	mux := http.NewServeMux()
	mux.HandleFunc("/health", h.Health)
	mux.HandleFunc("/set", h.Set)
	mux.HandleFunc("/get", h.Get)

	return &Server{
		mux: mux,
	}
}

func (s *Server) Routes() http.Handler {
	return s.mux
}
