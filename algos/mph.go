package algos

import (
	"crypto/md5"
	"encoding/binary"
	"fmt"
	"math"
	"math/rand"
	"sort"
)

func MultiProbeHashing(vNodeMultiplier int) []int {
	r := rand.Int63()
	nodeCounts := make([]int, NumNodes)

	var nodes []circularHashEntry
	for node := 0; node < NumNodes; node++ {
		value := md5.Sum([]byte(fmt.Sprintf("node %d-%d", r, node)))
		str := string(value[:])
		nodes = append(nodes, circularHashEntry{str, node})
	}
	sort.Slice(nodes, func(i, j int) bool {
		return nodes[i].hashValue < nodes[j].hashValue
	})
	for lease := 0; lease < NumLeases; lease++ {
		finalOwner := 0
		minDistance := math.MaxInt64
		for v := 0; v < vNodeMultiplier; v++ {
			value := md5.Sum([]byte(fmt.Sprintf("lease %d-%d-%d", r, lease, v)))
			str := string(value[:])
			entry := getCircularHashOwner(nodes, str)
			d := distance(entry.hashValue, str)
			if d < minDistance {
				// fmt.Printf("Lease %d: Picked %x(%d) for %x with distance %d\n", lease, entry.hashValue, entry.node, str, d)
				minDistance = d
				finalOwner = entry.node
			}
		}
		// fmt.Printf("Lease %d: Picked %d with distance %d\n", lease, finalOwner, minDistance)
		nodeCounts[finalOwner]++
	}
	// fmt.Printf("NODES: %+v\n", nodes)
	return nodeCounts
}

func distance(a, b string) int {
	buffer := make([]byte, 8)
	copy(buffer, []byte(a))
	aValue := binary.BigEndian.Uint64(buffer)
	copy(buffer, []byte(b))
	bValue := binary.BigEndian.Uint64(buffer)
	d := int(aValue - bValue)
	return d
}
