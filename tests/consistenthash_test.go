package tests

import (
	"distributed-cache/internal/consistenthash"
	"reflect"
	"testing"
)

func TestHashRingReturnsNilWhenEmpty(t *testing.T) {
	ring := consistenthash.NewHashRing()

	nodes := ring.GetNodes("user:1", 2)
	if nodes != nil {
		t.Fatalf("expected nil nodes for empty ring, got %v", nodes)
	}
}

func TestHashRingReturnsRequestedDistinctNodes(t *testing.T) {
	ring := consistenthash.NewHashRing()
	ring.AddNode("node1")
	ring.AddNode("node2")
	ring.AddNode("node3")

	nodes := ring.GetNodes("user:1", 2)
	if len(nodes) != 2 {
		t.Fatalf("expected 2 nodes, got %d: %v", len(nodes), nodes)
	}
	if nodes[0] == nodes[1] {
		t.Fatalf("expected distinct physical nodes, got %v", nodes)
	}
}

func TestHashRingCapsReplicaCountAtPhysicalNodeCount(t *testing.T) {
	ring := consistenthash.NewHashRing()
	ring.AddNode("node1")
	ring.AddNode("node2")

	nodes := ring.GetNodes("user:1", 5)
	if len(nodes) != 2 {
		t.Fatalf("expected replica count to be capped at 2, got %d: %v", len(nodes), nodes)
	}
}

func TestHashRingLookupIsStable(t *testing.T) {
	ring := consistenthash.NewHashRing()
	ring.AddNode("node1")
	ring.AddNode("node2")
	ring.AddNode("node3")

	first := ring.GetNodes("user:1", 2)
	second := ring.GetNodes("user:1", 2)

	if !reflect.DeepEqual(first, second) {
		t.Fatalf("expected stable lookup, first=%v second=%v", first, second)
	}
}
