package main

import (
	"distributed-cache/internal/hashing"
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "9000"
	}

	nodes := []string {
		"http://localhost:8001",
		"http://localhost:8002",
		"http://localhost:8003",
	}

	ring := hashing.NewHashRing()
	for _, node := range nodes {
		ring.AddNode(node)
	}

	http.HandleFunc("/",func(w http.ResponseWriter, r *http.Request){
		w.Write([]byte("Router running"))
	})

	http.HandleFunc("/set", func(w http.ResponseWriter, r *http.Request){
		key := r.URL.Query().Get("key")
		value := r.URL.Query().Get("value")

		node := ring.GetNode(key)
		url := node + "/set?key=" + key + "&value=" + value

		resp, err := http.Get(url)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error setting key"))
			return
		}
		defer resp.Body.Close()

		log.Println("SET key:", key, "→ node:", node)
	})

	http.HandleFunc("/get", func(w http.ResponseWriter, r *http.Request){
		key := r.URL.Query().Get("key")
		node := ring.GetNode(key)
		url := node + "/get?key=" + key

		resp, err := http.Post(url, "application/json", nil)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error getting key"))
			return
		}
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)
		w.WriteHeader(resp.StatusCode)
		w.Write(body)

		log.Println("GET key:", key, "→ node:", node)
	})

	log.Println("Router running on port 9000")
	log.Fatal(http.ListenAndServe(":9000",nil))
}