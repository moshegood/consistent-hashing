package algos

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
)

func RingHash(vNodeMultiplier int) []int {
	ring := MakeRing(NumNodes, vNodeMultiplier)
	r := rand.Int63()
	nodeCounts := make([]int, NumNodes)
	for lease := 0; lease < NumLeases; lease++ {
		value := HashingFunction(fmt.Sprintf("lease %d-%d", r, lease))
		nodeCounts[getCircularHashOwner(ring, value).node]++
	}
	return nodeCounts
}

// Helper functionality - used for both RingHashes and MPH
type circularHashEntry struct {
	hashValue uint64
	node      int
}

func (entry circularHashEntry) String() string {
	return fmt.Sprintf("%d: %x", entry.node, entry.hashValue)
}

func MakeRing(nodeCount, vNodeMultiplier int) []circularHashEntry {
	r := rand.Int63()
	var nodes []circularHashEntry
	for node := 0; node < nodeCount; node++ {
		for v := 0; v < vNodeMultiplier; v++ {
			value := HashingFunction(fmt.Sprintf("node %d-%d-%d", r, node, v))
			nodes = append(nodes, circularHashEntry{value, node})
		}
	}
	sort.Slice(nodes, func(i, j int) bool {
		return nodes[i].hashValue < nodes[j].hashValue
	})
	return nodes
}

func FindOwnerOfLeaseInRing(ring []circularHashEntry, lease int, vNodeMultiplier int) int {
	r := rand.Int63()
	finalOwner := 0
	var minDistance uint64 = math.MaxUint64
	for v := 0; v < vNodeMultiplier; v++ {
		value := HashingFunction(fmt.Sprintf("lease %d-%d-%d", r, lease, v))
		entry := getCircularHashOwner(ring, value)
		d := distance(entry.hashValue, value)
		if d < minDistance {
			// fmt.Printf("Lease %d: Picked %x(%d) for %x with distance %d\n", lease, entry.hashValue, entry.node, str, d)
			minDistance = d
			finalOwner = entry.node
		}
	}
	return finalOwner
}

func getCircularHashOwner(nodes []circularHashEntry, hashValue uint64) circularHashEntry {
	// Find the first entry that is greater than the hash value.
	// TODO: Optimize with binary search? It's such a small list.
	for _, entry := range nodes {
		if entry.hashValue >= hashValue {
			return entry
		}
	}
	return nodes[0]
}

func distance(a, b uint64) uint64 {
	return a - b
}
