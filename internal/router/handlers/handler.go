package handlers

import "distributed-cache/internal/consistenthash"

type Handler struct {
	ring *consistenthash.HashRing
}

func NewHandler(ring *consistenthash.HashRing) *Handler {
	return &Handler{
		ring: ring,
	}
}
