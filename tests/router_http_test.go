package tests

import (
	"distributed-cache/internal/router"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestRouterHealth(t *testing.T) {
	node1 := newFakeCacheNode()
	defer node1.Close()

	server := httptest.NewServer(router.NewServer([]string{node1.URL()}).Routes())
	defer server.Close()

	resp, body := get(t, server.URL+"/health")

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}
	if body != "OK\n" {
		t.Fatalf("unexpected response body: %q", body)
	}
}

func TestRouterSetGetDeleteFlow(t *testing.T) {
	node1 := newFakeCacheNode()
	defer node1.Close()
	node2 := newFakeCacheNode()
	defer node2.Close()

	server := httptest.NewServer(router.NewServer([]string{node1.URL(), node2.URL()}).Routes())
	defer server.Close()

	key := "user:1"
	value := "Yash"

	resp, body := get(t, server.URL+"/set?key="+url.QueryEscape(key)+"&value="+url.QueryEscape(value))
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected set status %d, got %d", http.StatusOK, resp.StatusCode)
	}
	if body != "Replicated successfully\n" {
		t.Fatalf("unexpected set response body: %q", body)
	}

	assertNodeValue(t, node1, key, value)
	assertNodeValue(t, node2, key, value)

	resp, body = get(t, server.URL+"/get?key="+url.QueryEscape(key))
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected get status %d, got %d", http.StatusOK, resp.StatusCode)
	}
	if body != value {
		t.Fatalf("expected value %q, got %q", value, body)
	}

	resp, body = get(t, server.URL+"/delete?key="+url.QueryEscape(key))
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected delete status %d, got %d", http.StatusOK, resp.StatusCode)
	}
	if body != "Deleted successfully\n" {
		t.Fatalf("unexpected delete response body: %q", body)
	}

	assertNodeMissing(t, node1, key)
	assertNodeMissing(t, node2, key)

	resp, body = get(t, server.URL+"/get?key="+url.QueryEscape(key))
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected get-after-delete status %d, got %d", http.StatusNotFound, resp.StatusCode)
	}
	if body != "Key not found\n" {
		t.Fatalf("unexpected get-after-delete response body: %q", body)
	}
}

func TestRouterRejectsMissingKeys(t *testing.T) {
	node1 := newFakeCacheNode()
	defer node1.Close()

	server := httptest.NewServer(router.NewServer([]string{node1.URL()}).Routes())
	defer server.Close()

	tests := []struct {
		name string
		path string
		body string
	}{
		{name: "set missing value", path: "/set?key=user:1", body: "Missing key or value\n"},
		{name: "get missing key", path: "/get", body: "Missing key\n"},
		{name: "delete missing key", path: "/delete", body: "Missing key\n"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, body := get(t, server.URL+tt.path)
			if resp.StatusCode != http.StatusBadRequest {
				t.Fatalf("expected status %d, got %d", http.StatusBadRequest, resp.StatusCode)
			}
			if body != tt.body {
				t.Fatalf("unexpected response body: %q", body)
			}
		})
	}
}

func TestRouterReturnsServerErrorWhenNoNodesAvailable(t *testing.T) {
	server := httptest.NewServer(router.NewServer(nil).Routes())
	defer server.Close()

	tests := []string{
		"/set?key=user:1&value=Yash",
		"/get?key=user:1",
		"/delete?key=user:1",
	}

	for _, path := range tests {
		t.Run(path, func(t *testing.T) {
			resp, _ := get(t, server.URL+path)
			if resp.StatusCode != http.StatusInternalServerError {
				t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, resp.StatusCode)
			}
		})
	}
}

func TestRouterSetSucceedsWithPartialReplication(t *testing.T) {
	node1 := newFakeCacheNode()
	defer node1.Close()
	node2 := newFakeCacheNode()
	node2.setFailureStatus = http.StatusInternalServerError
	defer node2.Close()

	server := httptest.NewServer(router.NewServer([]string{node1.URL(), node2.URL()}).Routes())
	defer server.Close()

	key := "user:partial-set"
	value := "partial"

	resp, body := get(t, server.URL+"/set?key="+url.QueryEscape(key)+"&value="+url.QueryEscape(value))
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}
	if body != "Replicated successfully\n" {
		t.Fatalf("unexpected response body: %q", body)
	}

	assertNodeValue(t, node1, key, value)
	assertNodeMissing(t, node2, key)
}

func TestRouterSetFailsWhenAllReplicasFail(t *testing.T) {
	node1 := newFakeCacheNode()
	node1.setFailureStatus = http.StatusInternalServerError
	defer node1.Close()
	node2 := newFakeCacheNode()
	node2.setFailureStatus = http.StatusInternalServerError
	defer node2.Close()

	server := httptest.NewServer(router.NewServer([]string{node1.URL(), node2.URL()}).Routes())
	defer server.Close()

	resp, body := get(t, server.URL+"/set?key=user:fail&value=value")
	if resp.StatusCode != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, resp.StatusCode)
	}
	if body != "Write failed on all nodes\n" {
		t.Fatalf("unexpected response body: %q", body)
	}
}

func TestRouterDeleteSucceedsWithPartialReplication(t *testing.T) {
	node1 := newFakeCacheNode()
	defer node1.Close()
	node2 := newFakeCacheNode()
	node2.delFailureStatus = http.StatusInternalServerError
	defer node2.Close()

	key := "user:partial-delete"
	node1.SetValue(key, "value")
	node2.SetValue(key, "value")

	server := httptest.NewServer(router.NewServer([]string{node1.URL(), node2.URL()}).Routes())
	defer server.Close()

	resp, body := get(t, server.URL+"/delete?key="+url.QueryEscape(key))
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}
	if body != "Deleted successfully\n" {
		t.Fatalf("unexpected response body: %q", body)
	}

	assertNodeMissing(t, node1, key)
	assertNodeValue(t, node2, key, "value")
}

func TestRouterDeleteFailsWhenAllReplicasFail(t *testing.T) {
	node1 := newFakeCacheNode()
	node1.delFailureStatus = http.StatusInternalServerError
	defer node1.Close()
	node2 := newFakeCacheNode()
	node2.delFailureStatus = http.StatusInternalServerError
	defer node2.Close()

	server := httptest.NewServer(router.NewServer([]string{node1.URL(), node2.URL()}).Routes())
	defer server.Close()

	resp, body := get(t, server.URL+"/delete?key=user:fail")
	if resp.StatusCode != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, resp.StatusCode)
	}
	if body != "Delete failed on all nodes\n" {
		t.Fatalf("unexpected response body: %q", body)
	}
}

func assertNodeValue(t *testing.T, node *fakeCacheNode, key, want string) {
	t.Helper()

	got, ok := node.Value(key)
	if !ok {
		t.Fatalf("expected node %s to contain key %q", node.URL(), key)
	}
	if got != want {
		t.Fatalf("expected node %s key %q value %q, got %q", node.URL(), key, want, got)
	}
}

func assertNodeMissing(t *testing.T, node *fakeCacheNode, key string) {
	t.Helper()

	if node.HasKey(key) {
		t.Fatalf("expected node %s to not contain key %q", node.URL(), key)
	}
}
