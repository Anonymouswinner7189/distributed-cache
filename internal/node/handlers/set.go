package handlers

import (
	"net/http"
	"time"
)

func (h *Handler) Set(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	value := r.URL.Query().Get("value")
	if key == "" || value == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid key or value\n"))
		return
	}

	h.cache.Set(key, value, 300*time.Second)
	w.Write([]byte("Key set\n"))
}
