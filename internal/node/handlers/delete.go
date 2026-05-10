package handlers

import (
	"net/http"
)

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	if key == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Missing key\n"))
		return
	}

	h.cache.Delete(key)
	w.Write([]byte("Key deleted\n"))
}
