package consistenthash

import (
	"hash/crc32"
	"slices"
	"sort"
	"strconv"
)

type HashRing struct {
	nodes         map[uint32]string
	keys          []uint32
	physicalNodes map[string]bool
}

func NewHashRing() *HashRing {
	return &HashRing{
		nodes:         make(map[uint32]string),
		physicalNodes: make(map[string]bool),
		keys:          []uint32{},
	}
}

func hashKey(key string) uint32 {
	return crc32.ChecksumIEEE([]byte(key))
}

func (h *HashRing) AddNode(node string) {
	h.physicalNodes[node] = true

	const virtualNodes = 100
	for i := range virtualNodes {
		virtualKey := node + "#" + strconv.Itoa(i)
		hash := hashKey(virtualKey)
		h.nodes[hash] = node
		h.keys = append(h.keys, hash)
	}

	slices.Sort(h.keys)
}

func (h *HashRing) GetNodes(key string, count int) []string {
	if len(h.nodes) == 0 {
		return nil
	}

	if count > len(h.physicalNodes) {
		count = len(h.physicalNodes)
	}

	hash := hashKey(key)
	idx := sort.Search(len(h.keys), func(i int) bool {
		return h.keys[i] >= hash
	})

	result := []string{}
	visited := make(map[string]bool)

	for len(result) < count {
		if idx == len(h.keys) {
			idx = 0
		}
		node := h.nodes[h.keys[idx]]
		// Avoid duplicate physical nodes (because of virtual nodes)
		if !visited[node] {
			result = append(result, node)
			visited[node] = true
		}
		idx++
	}

	return result
}
