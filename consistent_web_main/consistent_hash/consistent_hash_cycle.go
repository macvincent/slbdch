package consistent_hash

import (
	"crypto/sha256"
	"sort"
)

type consistentHash struct {
	// maps virtual nodes hash to their IP addresses
	vnodeHashToAddress map[byte]string
	// sorted list of virtual nodes
	sortedVnodeHash []byte
}

func NewConsistentHash(ipAddresses []string, replicaPerNode int) *consistentHash {
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
		}
	}

	// Sort the virtual nodes for easy lookup
	sort.Slice(ch.sortedVnodeHash, func(i, j int) bool {
		return ch.sortedVnodeHash[i] < ch.sortedVnodeHash[j]
	})

	return ch
}

func (ch *consistentHash) ValueLookup(value string) string {
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
