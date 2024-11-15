package algos

import (
	"crypto/md5"
	"fmt"
	"math/rand"
	"sort"
)

func RingHash(vNodeMultiplier int) []int {
	r := rand.Int63()
	var nodes []circularHashEntry
	for node := 0; node < NumNodes; node++ {
		for v := 0; v < vNodeMultiplier; v++ {
			value := md5.Sum([]byte(fmt.Sprintf("node %d-%d-%d", r, node, v)))
			str := string(value[:])
			nodes = append(nodes, circularHashEntry{str, node})
		}
	}
	sort.Slice(nodes, func(i, j int) bool {
		return nodes[i].hashValue < nodes[j].hashValue
	})
	nodeCounts := make([]int, NumNodes)
	leases := []string{}
	for lease := 0; lease < NumLeases; lease++ {
		value := md5.Sum([]byte(fmt.Sprintf("lease %d-%d", r, lease)))
		str := string(value[:])
		nodeCounts[getCircularHashOwner(nodes, str).node]++
		leases = append(leases, fmt.Sprintf("%x", str))
	}
	return nodeCounts
}

func getCircularHashOwner(nodes []circularHashEntry, hashValue string) circularHashEntry {
	// Find the first entry that is greater than the hash value.
	// TODO: Optimize with binary search? It's such a small list.
	for _, entry := range nodes {
		if entry.hashValue >= hashValue {
			return entry
		}
	}
	return nodes[0]
}

type circularHashEntry struct {
	hashValue string
	node      int
}

func (entry circularHashEntry) String() string {
	return fmt.Sprintf("%d: %x", entry.node, entry.hashValue)
}
