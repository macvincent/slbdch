package main

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
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

func NewTrie() *Trie {
	return &Trie{root: newNode()}
}

func getTrieKey(key string) uint32 {
	hashFunction := sha256.New()
	hashFunction.Write([]byte(key))
	hashValue := hashFunction.Sum(nil)
	trie_key := binary.LittleEndian.Uint32(hashValue[:4])
	hashFunction.Reset()
	return trie_key
}

func (t *Trie) insert(key string) {
	trie_key := getTrieKey(key)
	node := t.root
	for i := 0; i < 32; i++ {
		index := trie_key & 1
		if node.children[index] == nil {
			node.children[index] = newNode()
		}
		node = node.children[index]
		trie_key = (trie_key >> 1)
	}
	node.isServer = true
	node.ipAddress = key
}

func (t *Trie) search(key string) string {
	trie_key := getTrieKey(key)
	node := t.root
	for i := 0; i < 32; i++ {
		index := trie_key & 1
		if node.children[index] == nil {
			node = node.children[^index]
		} else {
			node = node.children[index]
		}
	}
	return node.ipAddress
}

func main() {
	trie := NewTrie()
	// Insert binary numbers
	trie.insert("1010")
	trie.insert("110")
	trie.insert("1001")
	// Search for binary numbers
	fmt.Println(trie.search("1010")) // true
	fmt.Println(trie.search("1101")) // false
}
