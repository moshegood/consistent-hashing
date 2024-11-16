package algos

import (
	"math"
	"math/rand"
	"sort"
)

func MultiProbeHashing(vNodeMultiplier int) []int {
	r := rand.Int63()
	nodeCounts := make([]int, NumNodes)

	var nodes []circularHashEntry
	for node := 0; node < NumNodes; node++ {
		value := HashUnaryKey("node", r, node)
		nodes = append(nodes, circularHashEntry{value, node})
	}
	sort.Slice(nodes, func(i, j int) bool {
		return nodes[i].hashValue < nodes[j].hashValue
	})
	for lease := 0; lease < NumLeases; lease++ {
		finalOwner := 0
		var minDistance uint64 = math.MaxUint64
		for v := 0; v < vNodeMultiplier; v++ {
			value := HashBinaryKey("lease", r, lease, v)
			entry := getCircularHashOwner(nodes, value)
			d := entry.hashValue - value
			if d < minDistance {
				// fmt.Printf("Lease %d: Picked %d(%d) for %d with distance %d\n", lease, entry.hashValue, entry.node, value, d)
				minDistance = d
				finalOwner = entry.node
			}
		}
		// fmt.Printf("Lease %d: Picked %d with distance %d\n", lease, finalOwner, minDistance)
		nodeCounts[finalOwner]++
	}
	// fmt.Printf("NODES: %+v\n", nodes)
	// fmt.Printf("NODE COUNTS: %+v\n", nodeCounts)
	return nodeCounts
}
