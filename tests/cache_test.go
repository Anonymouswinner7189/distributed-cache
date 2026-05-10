package tests

import (
	"distributed-cache/internal/cache"
	"testing"
	"time"
)

func TestCacheSetGetDelete(t *testing.T) {
	c := cache.NewCache()

	c.Set("user:1", "Yash", time.Minute)

	value, ok := c.Get("user:1")
	if !ok {
		t.Fatal("expected key to exist")
	}
	if value != "Yash" {
		t.Fatalf("expected value %q, got %q", "Yash", value)
	}

	c.Delete("user:1")

	if value, ok := c.Get("user:1"); ok {
		t.Fatalf("expected key to be deleted, got value %q", value)
	}
}

func TestCacheTTLExpiration(t *testing.T) {
	c := cache.NewCache()

	c.Set("short", "lived", 10*time.Millisecond)
	time.Sleep(20 * time.Millisecond)

	if value, ok := c.Get("short"); ok {
		t.Fatalf("expected key to expire, got value %q", value)
	}
}

func TestCacheZeroTTLDoesNotExpire(t *testing.T) {
	c := cache.NewCache()

	c.Set("persistent", "value", 0)
	time.Sleep(20 * time.Millisecond)

	value, ok := c.Get("persistent")
	if !ok {
		t.Fatal("expected key without TTL to remain")
	}
	if value != "value" {
		t.Fatalf("expected value %q, got %q", "value", value)
	}
}
