package handlers

import (
	"distributed-cache/internal/httpclient"
	"log"
	"net/http"
	"net/url"
)

func (h *Handler) Set(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	value := r.URL.Query().Get("value")
	if key == "" || value == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Missing key or value\n"))
		return
	}

	nodes := h.ring.GetNodes(key, 2) // replication factor of 2
	if len(nodes) == 0 {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("No nodes available\n"))
		return
	}

	successfulNodes := []string{}
	failedNodes := []string{}

	for _, node := range nodes {
		requestURL := node + "/set?key=" + url.QueryEscape(key) + "&value=" + url.QueryEscape(value)

		resp, err := httpclient.Client.Get(requestURL)
		if err != nil {
			failedNodes = append(failedNodes, node)
			log.Printf("Write failed key=%s node=%s err=%v", key, node, err)
			continue
		}

		if resp.StatusCode == http.StatusOK {
			successfulNodes = append(successfulNodes, node)
		} else {
			failedNodes = append(failedNodes, node)
			log.Printf("Write failed key=%s node=%s status=%d", key, node, resp.StatusCode)
		}

		resp.Body.Close()
	}

	if len(successfulNodes) == 0 {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Write failed on all nodes\n"))
		return
	}

	if len(successfulNodes) < len(nodes) {
		log.Printf(
			"Partial replication key=%s successful_nodes=%v failed_nodes=%v",
			key,
			successfulNodes,
			failedNodes,
		)
	}

	log.Printf(
		"SET key=%s successful_nodes=%v failed_nodes=%v",
		key,
		successfulNodes,
		failedNodes,
	)
	w.Write([]byte("Replicated successfully\n"))
}
