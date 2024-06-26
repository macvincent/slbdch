package consistent_hash

import (
	"fmt"
	"log"
	"math/rand/v2"
	"sort"
	"sync"
	"time"
)

type consistentHash struct {
	// maps virtual nodes value in cycle to their IP addresses
	vnodeHashToAddress map[uint32]string
	// sorted list of virtual nodes
	sortedVnodeHash []uint32
	nodeMap         map[string]ServerNode
	mux             sync.RWMutex
}

func (ch *consistentHash) getReplicaHashValues(ip string) []uint32 {
	ch.mux.RLock()
	defer ch.mux.RUnlock()
	hashValues := make([]uint32, 0)
	replica_count := ch.nodeMap[ip].Replicas
	for replicaNumber := 0; replicaNumber < replica_count; replicaNumber++ {
		hashValues = append(hashValues, getTrieKey(fmt.Sprintf("%s-%d", ip, replicaNumber)))
	}
	return hashValues
}

func NewConsistentHash(nodeMap map[string]ServerNode) *consistentHash {
	ch := &consistentHash{
		vnodeHashToAddress: make(map[uint32]string),
		sortedVnodeHash:    make([]uint32, 0),
		nodeMap:            nodeMap,
	}
	// Note no other thread has access to ch yet so we don't need a lock here
	// Add IP addresses to the hash table
	for ip := range nodeMap {
		for _, replica_hash := range ch.getReplicaHashValues(ip) {
			ch.sortedVnodeHash = append(ch.sortedVnodeHash, replica_hash)
			ch.vnodeHashToAddress[replica_hash] = ip
		}
	}

	// Sort the virtual nodes for easy lookup
	sort.Slice(ch.sortedVnodeHash, func(i, j int) bool {
		return ch.sortedVnodeHash[i] < ch.sortedVnodeHash[j]
	})

	return ch
}

func (ch *consistentHash) ValueLookup(value string) string {
	ch.mux.RLock()
	defer ch.mux.RUnlock()
	if len(ch.sortedVnodeHash) == 0 {
		return fmt.Errorf("no nodes available").Error()
	}

	hash := getTrieKey(value)

	// find the next virtual node that is clockwise to the given hash
	index := sort.Search(len(ch.sortedVnodeHash), func(i int) bool {
		return ch.sortedVnodeHash[i] >= hash
	})

	if index == len(ch.sortedVnodeHash) {
		index = 0
	}

	return ch.vnodeHashToAddress[ch.sortedVnodeHash[index]]
}

func (ch *consistentHash) InsertNode(ip_address string, replica_count int) {
	ch.mux.Lock()
	defer ch.mux.Unlock()
	// Update replica count if the node already exists
	if entry, exists := ch.nodeMap[ip_address]; exists {
		entry.Replicas = replica_count
		ch.nodeMap[ip_address] = entry
	} else {
		timeStamp := time.Now().Add(60 * time.Second)
		entry := ServerNode{IP: ip_address, Timestamp: timeStamp, Replicas: replica_count}
		ch.nodeMap[ip_address] = entry
	}

	// Insert the virtual nodes
	for replica_number, replica_hash := range ch.getReplicaHashValues(ip_address) {
		ch.sortedVnodeHash = append(ch.sortedVnodeHash, replica_hash)
		ch.vnodeHashToAddress[replica_hash] = ip_address
		log.Printf("Inserted IP %v, replica number %v\n", ip_address, replica_number)
	}

	// Sort the virtual nodes for easy lookup
	sort.Slice(ch.sortedVnodeHash, func(i, j int) bool {
		return ch.sortedVnodeHash[i] < ch.sortedVnodeHash[j]
	})
}

func (ch *consistentHash) DeleteNode(ip string) {
	ch.mux.Lock()
	defer ch.mux.Unlock()
	for _, replica_hash := range ch.getReplicaHashValues(ip) {
		delete(ch.vnodeHashToAddress, replica_hash)

		// Delete from sortedVnodeHash
		index := sort.Search(len(ch.sortedVnodeHash), func(i int) bool {
			return ch.sortedVnodeHash[i] >= replica_hash
		})

		if index < len(ch.sortedVnodeHash) && ch.sortedVnodeHash[index] == replica_hash {
			ch.sortedVnodeHash = append(ch.sortedVnodeHash[:index], ch.sortedVnodeHash[index+1:]...)
		}
	}
	// Delete from nodeMap
	delete(ch.nodeMap, ip)
}

// Thus funciton is used soley for testing purposes
func CycleMain() {
	timestamp := time.Now().Add(60 * time.Second)
	nodeList := []ServerNode{{IP: "localhost", Timestamp: timestamp, Replicas: 10}, {IP: "10.30.147.20", Timestamp: timestamp, Replicas: 3}}
	replica_count := 0

	nodeMap := make(map[string]ServerNode)
	for _, node := range nodeList {
		nodeMap[node.IP] = node
		replica_count += node.Replicas
	}
	consistentHash := NewConsistentHash(nodeMap)

	ipAddressCount := make(map[string]int)
	numCalls := 10000
	for i := 0; i < numCalls; i++ {
		url := fmt.Sprintf("www.%v.com", rand.IntN(100000))
		ip := consistentHash.ValueLookup(url)
		ipAddressCount[ip]++
	}

	fmt.Println("Expected vs True Count Per Node: ")
	for ip, node := range nodeMap {
		fmt.Printf("IP: %v, Expected Count: %v, True Count: %v\n", ip, node.Replicas*numCalls/replica_count, ipAddressCount[ip])
	}

	// Delete all nodes
	for ip := range nodeMap {
		consistentHash.DeleteNode(ip)
	}

	// Insert a new node
	consistentHash.InsertNode("localhost2", 2)

	// Search for a value
	fmt.Println(consistentHash.ValueLookup("www.google.com"))
}
