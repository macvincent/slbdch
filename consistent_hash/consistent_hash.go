package main

import (
	"crypto/sha256"
	"fmt"
	"math/rand/v2"
	"sort"
)

type consistentHash struct {
	// maps virtual nodes hash to their IP addresses
	vnodeHashToAddress map[byte]string
	// sorted list of virtual nodes
	sortedVnodeHash []byte
}

func newConsistentHash(ipAddresses []string, replicaPerNode int) *consistentHash {
	ch := &consistentHash{
		vnodeHashToAddress: make(map[byte]string),
		sortedVnodeHash:    make([]byte, 0),
	}

	hashFunction := sha256.New()

	// Add IP addresses to the hash table
	for _, ip := range ipAddresses {
		hashFunction.Write([]byte(ip))
		virtualNodePositions := hashFunction.Sum(nil)
		hashFunction.Reset()

		for i := 31; i > 31-replicaPerNode; i-- {
			ch.sortedVnodeHash = append(ch.sortedVnodeHash, virtualNodePositions[i])
			ch.vnodeHashToAddress[virtualNodePositions[i]] = ip
			fmt.Printf("IP: %v, Virtual Node: %v\n", ip, virtualNodePositions[i])
		}
	}

	// Sort the virtual nodes for easy lookup
	sort.Slice(ch.sortedVnodeHash, func(i, j int) bool {
		return ch.sortedVnodeHash[i] < ch.sortedVnodeHash[j]
	})

	return ch
}

func (ch *consistentHash) valueLookup(value string) string {
	hashFunction := sha256.New()
	hashFunction.Write([]byte(value))
	hash := hashFunction.Sum(nil)[31]

	// find the next virtual node that is clockwise to the given hash
	index := sort.Search(len(ch.sortedVnodeHash), func(i int) bool {
		return ch.sortedVnodeHash[i] >= hash
	})

	if index == len(ch.sortedVnodeHash) {
		index = 0
	}

	return ch.vnodeHashToAddress[ch.sortedVnodeHash[index]]
}

func main() {
	fmt.Println("*** Consistent Hash Initialization***")
	nodeAddresses := []string{"8.8.8.8", "1.1.1.1", "208.67.222.222", "208.67.220.220", "9.9.9.9"}
	replicaPerNode := 10
	consistentHash := newConsistentHash(nodeAddresses, replicaPerNode)

	fmt.Println("\n*** Sanity Check ***")
	// a map to count the number of url valuwa per IP address
	ipAddressCount := make(map[string]int)

	for i := 0; i < replicaPerNode*1000; i++ {
		url := fmt.Sprintf("www.%v.com", rand.IntN(100000))
		ip := consistentHash.valueLookup(url)
		ipAddressCount[ip]++
	}

	// print the number of url values per IP address
	for ip, count := range ipAddressCount {
		fmt.Printf("IP: %v, Count: %v\n", ip, count)
	}
	fmt.Println("Ideal Count Per Node: ", replicaPerNode*1000/len(nodeAddresses))
}
