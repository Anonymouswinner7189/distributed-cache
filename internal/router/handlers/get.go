package handlers

import (
	"distributed-cache/internal/httpclient"
	"io"
	"log"
	"net/http"
	"net/url"
)

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
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

	checkedNodes := []string{}
	failedNodes := []string{}

	for _, node := range nodes {
		requestURL := node + "/get?key=" + url.QueryEscape(key)

		resp, err := httpclient.Client.Get(requestURL)
		if err != nil {
			failedNodes = append(failedNodes, node)
			log.Printf("Read failed key=%s node=%s err=%v", key, node, err)
			continue
		}

		checkedNodes = append(checkedNodes, node)

		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			log.Printf(
				"GET key=%s found_node=%s checked_nodes=%v failed_nodes=%v",
				key,
				node,
				checkedNodes,
				failedNodes,
			)

			w.WriteHeader(http.StatusOK)
			w.Write(body)
			return
		}

		log.Printf("GET key=%s not_found_node=%s status=%d", key, node, resp.StatusCode)
	}

	if len(checkedNodes) == 0 {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("All nodes are unreachable\n"))
		return
	}

	log.Printf(
		"GET key=%s not_found checked_nodes=%v failed_nodes=%v",
		key,
		checkedNodes,
		failedNodes,
	)

	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("Key not found\n"))
}
