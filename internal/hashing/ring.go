package hashing

import (
	"hash/crc32"
	"slices"
	"sort"
	"strconv"
)

type HashRing struct {
	nodes map[uint32]string
	keys []uint32
}

func NewHashRing() *HashRing {
	return &HashRing {
		nodes : make(map[uint32]string),
		keys: []uint32{},
	}
}

func hashKey(key string) uint32 {
	return crc32.ChecksumIEEE([]byte(key))
}

func (h *HashRing) AddNode(node string) {
	const virtualNodes = 100
	for i:=range virtualNodes {
		virtualKey := node + "#" + strconv.Itoa(i)
		hash := hashKey(virtualKey)
		h.nodes[hash] = node
		h.keys = append(h.keys, hash)
	}

	slices.Sort(h.keys)
}

func (h *HashRing) GetNode(key string) string {
	if len(h.nodes) == 0 {
		return ""
	}

	hash := hashKey(key)

	idx := sort.Search(len(h.keys), func(i int) bool {
		return h.keys[i] >= hash
	})

	if(idx == len(h.keys)) {
		idx = 0
	}

	return h.nodes[h.keys[idx]]
}