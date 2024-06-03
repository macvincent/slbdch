package consistent_hash

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"log"
	"math/rand/v2"
	"strconv"
	"sync"
	"time"
)

type ServerNode struct {
	IP        string
	Timestamp time.Time
	Replicas  int
}

type TrieNode struct {
	children  [2]*TrieNode
	isServer  bool
	ipAddress string
}

type Trie struct {
	root    *TrieNode
	nodeMap map[string]ServerNode
	mux sync.RWMutex
}

func newNode() *TrieNode {
	return &TrieNode{}
}

func NewTrie(nodeMap map[string]ServerNode) *Trie {
	trie := &Trie{
		root:    newNode(),
		nodeMap: nodeMap,
	}
	// Note no other thread has access to trie yet so we don't need a lock here

	// Add IP addresses to the hash table
	for ip := range nodeMap {
		trie.InsertNode(ip, nodeMap[ip].Replicas)
	}
	return trie
}

func getTrieKey(key string) uint32 {
	hashFunction := sha256.New()
	hashFunction.Write([]byte(key))
	hashValue := hashFunction.Sum(nil)
	trie_key := binary.LittleEndian.Uint32(hashValue[:4])
	hashFunction.Reset()
	return trie_key
}

func (t *Trie) insert(ip_address string, replica_number int) {
	trie_key := getTrieKey(ip_address + strconv.Itoa(replica_number))
	node := t.root
	for i := 31; i >= 0; i-- {
		index := (trie_key >> i) & 1
		if node.children[index] == nil {
			node.children[index] = newNode()
		}
		node = node.children[index]
	}
	node.isServer = true
	node.ipAddress = ip_address
}

func (t *Trie) ValueLookup(key string) string {
	t.mux.RLock()
	defer t.mux.RUnlock()
	trie_key := getTrieKey(key)
	node := t.root
	for i := 31; i >= 0; i-- {
		index := (trie_key >> i) & 1
		if node.children[index] == nil {
			node = node.children[1-index]
		} else {
			node = node.children[index]
		}
	}
	return node.ipAddress
}

func (t *Trie) DeleteNode(ip_address string) {
	t.mux.Lock()
	defer t.mux.Unlock()
	replica_count := t.nodeMap[ip_address].Replicas
	for replica_number := 0; replica_number < replica_count; replica_number++ {
		trie_key := getTrieKey(ip_address + strconv.Itoa(replica_number))
		t.root = t.deleteRecursive(t.root, trie_key, 31)
	}
	delete(t.nodeMap, ip_address)
}

func (t *Trie) deleteRecursive(node *TrieNode, trie_key uint32, bitIndex int) *TrieNode {
	if node == nil {
		return nil
	}
	// Find the index (0 or 1) where the node needs to go
	index := (trie_key >> bitIndex) & 1

	// If the other child is also nil, then we return nil
	if bitIndex == 0 {
		// If bitIndex is 0, we are at the parent of the leaf node.
		// We simply set the node's child to nil. It will be automatically
		// garbage collected.
		node.children[index] = nil
	} else {
		// Otherwise, we recursively call the delete function on the child
		// node
		node.children[index] = t.deleteRecursive(node.children[index], trie_key, bitIndex-1)
	}

	// If both children of a node are nil, we simply return nil.
	if node.children[index] == nil && node.children[1-index] == nil && node != t.root {
		return nil
	}
	return node
}

func (t *Trie) InsertNode(ip_address string, replica_count int) {
	t.mux.Lock()
	defer t.mux.Unlock()
	// Upadte replica count if the node already exists
	if entry, ok := t.nodeMap[ip_address]; ok {
		entry.Replicas = replica_count
		t.nodeMap[ip_address] = entry
	} else {
		timestamp := time.Now().Add(60 * time.Second)
		entry := ServerNode{IP: ip_address, Timestamp: timestamp, Replicas: replica_count}
		t.nodeMap[ip_address] = entry
	}
	for replica_number := 0; replica_number < replica_count; replica_number++ {
		t.insert(ip_address, replica_number)
		log.Printf("Inserted IP %v, replica number %v\n", ip_address, replica_number)
	}
}

// Thus function is used solely for testing purposes
func KademliaMain() {
	timestamp := time.Now().Add(60 * time.Second)
	nodeList := []ServerNode{{IP: "localhost", Timestamp: timestamp, Replicas: 10}, {IP: "10.30.147.20", Timestamp: timestamp, Replicas: 3}}
	replica_count := 0

	nodeMap := make(map[string]ServerNode)
	for _, node := range nodeList {
		nodeMap[node.IP] = node
		replica_count += node.Replicas
	}
	consistentHash := NewTrie(nodeMap)

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

	// Delete a node
	for ip := range nodeMap {
		consistentHash.DeleteNode(ip)
	}

	consistentHash.InsertNode("localhost2", 1)
	fmt.Println(consistentHash.ValueLookup("www.google.com"))
}
