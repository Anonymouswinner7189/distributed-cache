package tests

import (
	"net/http"
	"net/http/httptest"
	"sync"
)

type fakeCacheNode struct {
	server *httptest.Server

	mu               sync.Mutex
	store            map[string]string
	setFailureStatus int
	getFailureStatus int
	delFailureStatus int
}

func newFakeCacheNode() *fakeCacheNode {
	node := &fakeCacheNode{
		store: make(map[string]string),
	}

	node.server = httptest.NewServer(http.HandlerFunc(node.handle))
	return node
}

func (n *fakeCacheNode) URL() string {
	return n.server.URL
}

func (n *fakeCacheNode) Close() {
	n.server.Close()
}

func (n *fakeCacheNode) SetValue(key, value string) {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.store[key] = value
}

func (n *fakeCacheNode) HasKey(key string) bool {
	n.mu.Lock()
	defer n.mu.Unlock()
	_, ok := n.store[key]
	return ok
}

func (n *fakeCacheNode) Value(key string) (string, bool) {
	n.mu.Lock()
	defer n.mu.Unlock()
	value, ok := n.store[key]
	return value, ok
}

func (n *fakeCacheNode) handle(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/health":
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK\n"))
	case "/set":
		n.handleSet(w, r)
	case "/get":
		n.handleGet(w, r)
	case "/delete":
		n.handleDelete(w, r)
	default:
		http.NotFound(w, r)
	}
}

func (n *fakeCacheNode) handleSet(w http.ResponseWriter, r *http.Request) {
	if n.setFailureStatus != 0 {
		w.WriteHeader(n.setFailureStatus)
		return
	}

	key := r.URL.Query().Get("key")
	value := r.URL.Query().Get("value")
	if key == "" || value == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	n.SetValue(key, value)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Key set\n"))
}

func (n *fakeCacheNode) handleGet(w http.ResponseWriter, r *http.Request) {
	if n.getFailureStatus != 0 {
		w.WriteHeader(n.getFailureStatus)
		return
	}

	key := r.URL.Query().Get("key")
	value, ok := n.Value(key)
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Key not found\n"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(value))
}

func (n *fakeCacheNode) handleDelete(w http.ResponseWriter, r *http.Request) {
	if n.delFailureStatus != 0 {
		w.WriteHeader(n.delFailureStatus)
		return
	}

	key := r.URL.Query().Get("key")
	if key == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	n.mu.Lock()
	delete(n.store, key)
	n.mu.Unlock()

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Key deleted\n"))
}
