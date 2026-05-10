package tests

import (
	"distributed-cache/internal/node"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNodeHealth(t *testing.T) {
	server := httptest.NewServer(node.NewServer().Routes())
	defer server.Close()

	resp, body := get(t, server.URL+"/health")

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}
	if body != "OK\n" {
		t.Fatalf("unexpected response body: %q", body)
	}
}

func TestNodeSetGetDeleteFlow(t *testing.T) {
	server := httptest.NewServer(node.NewServer().Routes())
	defer server.Close()

	resp, body := get(t, server.URL+"/set?key=user:1&value=Yash")
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected set status %d, got %d", http.StatusOK, resp.StatusCode)
	}
	if body != "Key set\n" {
		t.Fatalf("unexpected set response body: %q", body)
	}

	resp, body = get(t, server.URL+"/get?key=user:1")
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected get status %d, got %d", http.StatusOK, resp.StatusCode)
	}
	if body != "Yash" {
		t.Fatalf("expected value %q, got %q", "Yash", body)
	}

	resp, body = get(t, server.URL+"/delete?key=user:1")
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected delete status %d, got %d", http.StatusOK, resp.StatusCode)
	}
	if body != "Key deleted\n" {
		t.Fatalf("unexpected delete response body: %q", body)
	}

	resp, body = get(t, server.URL+"/get?key=user:1")
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected get-after-delete status %d, got %d", http.StatusNotFound, resp.StatusCode)
	}
	if body != "Key not found\n" {
		t.Fatalf("unexpected get-after-delete response body: %q", body)
	}
}

func TestNodeRejectsMissingKeys(t *testing.T) {
	server := httptest.NewServer(node.NewServer().Routes())
	defer server.Close()

	tests := []struct {
		name string
		path string
		body string
	}{
		{name: "set missing value", path: "/set?key=user:1", body: "Invalid key or value\n"},
		{name: "get missing key", path: "/get", body: "Invalid key\n"},
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

func get(t *testing.T, target string) (*http.Response, string) {
	t.Helper()

	resp, err := http.Get(target)
	if err != nil {
		t.Fatalf("GET %s failed: %v", target, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("reading response body failed: %v", err)
	}

	return resp, string(body)
}
