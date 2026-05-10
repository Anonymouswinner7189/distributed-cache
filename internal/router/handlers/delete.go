package handlers

import (
	"distributed-cache/internal/httpclient"
	"log"
	"net/http"
	"net/url"
)

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	if key == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Missing key\n"))
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
		requestURL := node + "/delete?key=" + url.QueryEscape(key)

		resp, err := httpclient.Client.Get(requestURL)
		if err != nil {
			failedNodes = append(failedNodes, node)
			log.Printf("Delete failed key=%s node=%s err=%v\n", key, node, err)
			continue
		}

		if resp.StatusCode == http.StatusOK {
			successfulNodes = append(successfulNodes, node)
		} else {
			failedNodes = append(failedNodes, node)
			log.Printf("Delete failed key=%s node=%s status=%d\n", key, node, resp.StatusCode)
		}

		resp.Body.Close()
	}

	if len(successfulNodes) == 0 {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Delete failed on all nodes\n"))
		return
	}

	if len(successfulNodes) < len(nodes) {
		log.Printf(
			"Partial delete key=%s successful_nodes=%v failed_nodes=%v\n",
			key,
			successfulNodes,
			failedNodes,
		)
	}

	log.Printf(
		"DELETE key=%s successful_nodes=%v failed_nodes=%v\n",
		key,
		successfulNodes,
		failedNodes,
	)

	w.Write([]byte("Deleted successfully\n"))
}
