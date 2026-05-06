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
		"http://node1:8001",
		"http://node2:8002",
		"http://node3:8003",
	}

	ring := hashing.NewHashRing()
	for _, node := range nodes {
		ring.AddNode(node)
	}

	http.HandleFunc("/",func(w http.ResponseWriter, r *http.Request){
		w.Write([]byte("Router running\n"))
	})

	http.HandleFunc("/set", func(w http.ResponseWriter, r *http.Request){
		key := r.URL.Query().Get("key")
		value := r.URL.Query().Get("value")

		nodes := ring.GetNodes(key, 2) // replication factor of 2

		for _, node := range nodes {
			url := node + "/set?key=" + key + "&value=" + value

			resp, err := http.Post(url, "application/json", nil)
			if err != nil {
				log.Println("Write failed for node:", node)
				continue
			}
			resp.Body.Close()
		}

		log.Println("SET key:", key, "→ nodes:", nodes)
		w.Write([]byte("Replicated successfully\n"))
	})

	http.HandleFunc("/get", func(w http.ResponseWriter, r *http.Request){
		key := r.URL.Query().Get("key")
		nodes := ring.GetNodes(key,2)

		var lastErr error

		for _, node := range nodes {
			url := node + "/get?key=" + key

			resp, err := http.Get(url)
			if err != nil {
				log.Println("Node unreachable:", node)
				lastErr = err
				continue
			}

			body, _ := io.ReadAll(resp.Body)
			resp.Body.Close()

			if resp.StatusCode == http.StatusOK {
				log.Println("GET key:", key, "FOUND at node:", node)
				w.WriteHeader(http.StatusOK)
				w.Write(body)
				return
			}

			log.Println("GET key:", key, "NOT FOUND at node:", node)
		}

		if lastErr != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("All nodes unreachable\n"))
			return
		}

		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Key not found\n"))
	})

	log.Println("Router running on port 9000")
	log.Fatal(http.ListenAndServe(":9000",nil))
}