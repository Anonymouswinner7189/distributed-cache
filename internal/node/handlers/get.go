package handlers

import "net/http"

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	if key == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid key\n"))
		return
	}

	val, ok := h.cache.Get(key)
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Key not found\n"))
		return
	}

	w.Write([]byte(val))
}
