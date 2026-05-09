package handlers

import "distributed-cache/internal/cache"

type Handler struct {
	cache *cache.Cache
}

func NewHandler(c *cache.Cache) *Handler {
	return &Handler{
		cache: c,
	}
}
