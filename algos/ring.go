package algos

import (
	"crypto/md5"
	"encoding/binary"
	"fmt"
	"math"
	"math/rand"
	"sort"
)

func RingHash(vNodeMultiplier int) []int {
	ring := MakeRing(NumNodes, vNodeMultiplier)
	r := rand.Int63()
	nodeCounts := make([]int, NumNodes)
	leases := []string{}
	for lease := 0; lease < NumLeases; lease++ {
		value := md5.Sum([]byte(fmt.Sprintf("lease %d-%d", r, lease)))
		str := string(value[:])
		nodeCounts[getCircularHashOwner(ring, str).node]++
		leases = append(leases, fmt.Sprintf("%x", str))
	}
	return nodeCounts
}

// Helper functionality - used for both RingHashes and MPH
type circularHashEntry struct {
	hashValue string
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
			value := md5.Sum([]byte(fmt.Sprintf("node %d-%d-%d", r, node, v)))
			str := string(value[:])
			nodes = append(nodes, circularHashEntry{str, node})
		}
	}
	sort.Slice(nodes, func(i, j int) bool {
		return nodes[i].hashValue < nodes[j].hashValue
	})
	return nodes
}

func FindOwnerOfLeaseInRing(ring []circularHashEntry, lease int, vNodeMultiplier int) {
	r := rand.Int63()
	finalOwner := 0
	var minDistance uint64 = math.MaxUint64
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
	return finalOwner
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

func distance(a, b string) uint64 {
	buffer := make([]byte, 8)
	copy(buffer, []byte(a))
	aValue := binary.BigEndian.Uint64(buffer)
	copy(buffer, []byte(b))
	bValue := binary.BigEndian.Uint64(buffer)
	d := aValue - bValue
	return d
}
