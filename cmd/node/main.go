package main

import (
	"distributed-cache/internal/cache"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8001"
	}

	c := cache.NewCache()
	c.StartCleanup(10*time.Second)

	http.HandleFunc("/health",func(w http.ResponseWriter,r *http.Request){
		w.Write([]byte("OK"))
	})

	http.HandleFunc("/set", func(w http.ResponseWriter, r *http.Request){
		key := r.URL.Query().Get("key")
		value := r.URL.Query().Get("value")

		c.Set(key, value, 300*time.Second)
		w.Write([]byte("Key set"))
	})

	http.HandleFunc("/get",func(w http.ResponseWriter,r *http.Request){
		key := r.URL.Query().Get("key")

		val, ok := c.Get(key)
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("Key not found"))
			return
		}

		w.Write([]byte(val))
	})

	log.Println("Starting cache node on port",port)
	log.Fatal(http.ListenAndServe(":"+port,nil))
}