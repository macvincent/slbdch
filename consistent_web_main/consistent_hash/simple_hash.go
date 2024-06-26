package consistent_hash

import (
	"sort"
	"sync"
	"time"
)

type SimpleHash struct {
	orderedKeys    []string
	nodeMap        map[string]ServerNode
	sizeInclRepls  int
	mux            sync.RWMutex
}

func NewSimpleHash(nodeMap map[string]ServerNode) *SimpleHash {
	var size int = 0
	orderedKeys := make([]string, 0)

	for ip, node := range nodeMap  {
		size += node.Replicas
		orderedKeys = append(orderedKeys, ip)
	}
	h := &SimpleHash{
		orderedKeys:                 orderedKeys,
		nodeMap:                     nodeMap,
		sizeInclRepls:               size,
	}

	return h
}

func (h *SimpleHash) InsertNode(ip_address string, replica_count int) {
	h.mux.Lock()
	defer h.mux.Unlock()
	if entry, ok := h.nodeMap[ip_address]; ok {
		h.sizeInclRepls -= entry.Replicas
		entry.Replicas = replica_count
		h.sizeInclRepls += entry.Replicas

		// Note that this changes the distribution for other keys after the ip address
		// this is fine because regardless the other keys will be changed
		h.nodeMap[ip_address] = entry
	} else {
		timestamp := time.Now().Add(60 * time.Second)
		h.nodeMap[ip_address] = ServerNode{IP: ip_address, Timestamp: timestamp, Replicas: replica_count}
		h.orderedKeys = append(h.orderedKeys, ip_address)
		h.sizeInclRepls += replica_count
	}
}

func (h *SimpleHash) DeleteNode(ip string) {
	h.mux.Lock()
	defer h.mux.Unlock()
	h.sizeInclRepls -= h.nodeMap[ip].Replicas
	index := sort.Search(len(h.orderedKeys), func(i int) bool {
		return h.orderedKeys[i] >= ip
	})

	if index < len(h.orderedKeys) && h.orderedKeys[index] == ip {
		h.orderedKeys = append(h.orderedKeys[:index], h.orderedKeys[index+1:]...)
	}
	delete(h.nodeMap, ip)
}

func (h *SimpleHash) ValueLookup(value string) string {
	h.mux.RLock()
	defer h.mux.RUnlock()
	replica_idx := (int) (getTrieKey(value)) % h.sizeInclRepls
	var tempSize int = 0
	for _, ip := range h.orderedKeys {
		tempSize += h.nodeMap[ip].Replicas
		if replica_idx < tempSize {
			return ip
		}
	}
	return ""
}