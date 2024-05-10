package consistent_hash

import (
	"crypto/sha256"
	"encoding/binary"
	"log"
	"strconv"
)

type TrieNode struct {
	children  [2]*TrieNode
	isServer  bool
	ipAddress string
}

type Trie struct {
	root *TrieNode
}

func newNode() *TrieNode {
	return &TrieNode{}
}

func NewTrie(ipAddresses []string, replicaPerNode int) *Trie {
	trie := &Trie{root: newNode()}

	// Add IP addresses to the hash table
	for _, ip := range ipAddresses {
		for replica_number := 0; replica_number < replicaPerNode; replica_number++ {
			trie.insert(ip, replica_number)
			log.Printf("Inserted IP %v, replica number %v\n", ip, replica_number)
		}
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

func (t *Trie) Search(key string) string {
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
