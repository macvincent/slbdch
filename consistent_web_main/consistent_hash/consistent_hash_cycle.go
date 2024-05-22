package consistent_hash

import (
	"fmt"
	"math/rand/v2"
	"sort"
	"time"
)

type consistentHash struct {
	// maps virtual nodes value in cycle to their IP addresses
	vnodeHashToAddress map[uint32]string
	// sorted list of virtual nodes
	sortedVnodeHash []uint32
	nodeMap         map[string]ServerNode
}

func (ch consistentHash) getReplicaHashValues(ip string) []uint32 {
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
	for _, replica_hash := range ch.getReplicaHashValues(ip) {
		delete(ch.vnodeHashToAddress, replica_hash)

		// Delete from sortedVnodeHash
		index := sort.Search(len(ch.sortedVnodeHash), func(i int) bool {
			return ch.sortedVnodeHash[i] >= replica_hash
		})

		if index < len(ch.sortedVnodeHash) && ch.sortedVnodeHash[index] == replica_hash {
			ch.sortedVnodeHash = append(ch.sortedVnodeHash[:index], ch.sortedVnodeHash[index+1:]...)
		}

		delete(ch.nodeMap, ip)
	}
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

	for i := 0; i < 1000; i++ {
		url := fmt.Sprintf("www.%v.com", rand.IntN(100000))
		ip := consistentHash.ValueLookup(url)
		ipAddressCount[ip]++
	}

	// print the number of url values per IP address
	for ip, count := range ipAddressCount {
		fmt.Printf("IP: %v, Count: %v\n", ip, count)
	}
	fmt.Println("Ideal Count Per Node: ")
	for ip, node := range nodeMap {
		fmt.Printf("IP: %v, Count: %v\n", ip, node.Replicas*1000/replica_count)
	}

	// Delete a node
	for ip := range nodeMap {
		consistentHash.DeleteNode(ip)
	}

	fmt.Println(consistentHash.ValueLookup("www.google.com"))
}
