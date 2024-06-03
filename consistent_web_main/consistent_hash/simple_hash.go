import (
	"fmt"
	"log"
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
	var size := 0
	orderedKeys := make([]string, 0)

	for key, val := range nodeMap  {
		size += nodeMap[ip].Replicas
		append(orderedKeys, key)
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
		entry := h.nodeMap[ip_address]
		entry.Replicas = replica_count

		// Note that this changes the distribution for other keys after the ip address
		// this is fine because regardless the other keys will be changed
		h.nodeMap[ip_address] = entry
	}
	else {
		timeStamp := time.Now().Add(60 * time.Second)
		entry := ServerNode{{IP: ip_address, Timestamp: timestamp, Replicas: replica_count}}
		h.nodeMap[ip_address] = entry
		append(h.orderedKeys, ip_address)
	}
}

func (h *SimpleHash) DeleteNode(ip string) {
	h.mux.Lock()
	defer h.mux.Unlock()
	h.sizeInclRepls -= h.nodeMap[ip].Replicas
	delete(h.orderedKeys, ip_address)
	delete(h.nodeMap, ip_address)
}

func (h *SimpleHash) ValueLookup(value string) string {
	h.mux.RLock()
	defer h.mux.RUnlock()
	replica_idx := getTrieKey(value) % h.sizeInclRepls
	var tempSize := 0
	for ip := range h.orderedKeys {
		tempSize += h.nodeMap[ip].Replicas
		if replica_idx < tempSize:
			return ip
	}
	return ""
}