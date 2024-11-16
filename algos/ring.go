package algos

import (
	"cmp"
	"fmt"
	"math/rand"
	"slices"
	"sort"
)

func RingHash(vNodeMultiplier int) []int {
	r := rand.Int63()
	var nodes []circularHashEntry
	for node := 0; node < NumNodes; node++ {
		for v := 0; v < vNodeMultiplier; v++ {
			value := HashBinaryKey("node", r, node, v)
			nodes = append(nodes, circularHashEntry{value, node})
		}
	}
	sort.Slice(nodes, func(i, j int) bool {
		return nodes[i].hashValue < nodes[j].hashValue
	})
	nodeCounts := make([]int, NumNodes)
	for lease := 0; lease < NumLeases; lease++ {
		value := HashUnaryKey("lease", r, lease)
		nodeCounts[getCircularHashOwner(nodes, value).node]++
	}
	return nodeCounts
}

func getCircularHashOwner(nodes []circularHashEntry, hashValue uint64) circularHashEntry {
	// Find the first entry that is greater than the hash value.
	index, _ := slices.BinarySearchFunc(nodes, hashValue, func(entry circularHashEntry, hashValue uint64) int {
		return cmp.Compare(entry.hashValue, hashValue)
	})
	if index == len(nodes) {
		return nodes[0]
	}
	return nodes[index]
}

type circularHashEntry struct {
	hashValue uint64
	node      int
}

func (entry circularHashEntry) String() string {
	return fmt.Sprintf("%d: %d", entry.node, entry.hashValue)
}
