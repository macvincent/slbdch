package consistent_hash

import (
	"fmt"
	"math/rand/v2"
	"sort"
)

type consistentHash struct {
	// maps virtual nodes hash to their IP addresses
	vnodeHashToAddress map[uint32]string
	// sorted list of virtual nodes
	sortedVnodeHash []uint32
	replicaPerNode  int
}

func NewConsistentHash(ipAddresses []string, replicaPerNode int) *consistentHash {
	ch := &consistentHash{
		vnodeHashToAddress: make(map[uint32]string),
		sortedVnodeHash:    make([]uint32, 0),
		replicaPerNode:     replicaPerNode,
	}

	// Add IP addresses to the hash table
	for _, ip := range ipAddresses {
		for replicaNumber := 0; replicaNumber < replicaPerNode; replicaNumber++ {
			replica_hash := getTrieKey(fmt.Sprintf("%s-%d", ip, replicaNumber))
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

func (ch *consistentHash) DeleteNode(ip string) {
	for replicaNumber := 0; replicaNumber < ch.replicaPerNode; replicaNumber++ {
		replica_hash := getTrieKey(fmt.Sprintf("%s-%d", ip, replicaNumber))
		delete(ch.vnodeHashToAddress, replica_hash)

		// Delete from sortedVnodeHash
		index := sort.Search(len(ch.sortedVnodeHash), func(i int) bool {
			return ch.sortedVnodeHash[i] >= replica_hash
		})

		if index < len(ch.sortedVnodeHash) && ch.sortedVnodeHash[index] == replica_hash {
			ch.sortedVnodeHash = append(ch.sortedVnodeHash[:index], ch.sortedVnodeHash[index+1:]...)
		}
	}
}

// Thus funciton is used soley for testing purposes
func CycleMain() {
	nodeAddresses := []string{"8.8.8.8", "1.1.1.1", "208.67.222.222", "208.67.220.220", "9.9.9.9"}
	replicaPerNode := 10

	consistentHash := NewConsistentHash(nodeAddresses, replicaPerNode)

	ipAddressCount := make(map[string]int)

	for i := 0; i < replicaPerNode*1000; i++ {
		url := fmt.Sprintf("www.%v.com", rand.IntN(100000))
		ip := consistentHash.ValueLookup(url)
		ipAddressCount[ip]++
	}

	// print the number of url values per IP address
	for ip, count := range ipAddressCount {
		fmt.Printf("IP: %v, Count: %v\n", ip, count)
	}
	fmt.Println("Ideal Count Per Node: ", replicaPerNode*1000/len(nodeAddresses))

	for _, ip := range nodeAddresses {
		consistentHash.DeleteNode(ip)
	}

	fmt.Println(consistentHash.ValueLookup("www.google.com"))
}
