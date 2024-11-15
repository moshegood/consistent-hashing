package main

import (
	"crypto/md5"
	"encoding/binary"
	"fmt"
	"math"
	"math/rand"
	"sort"
	"time"
)

const numNodes = 10
const numLeases = 1000
const runs = 200

func main() {
	rand.Seed(time.Now().UnixNano())
	fmt.Println("=============== RingHash ===============")
	for vNodes := 1; vNodes < 1000; vNodes *= 3 {
		runTest(fmt.Sprintf("RingHash(%d)", vNodes), func() []int { return RingHash(vNodes) })
	}

	fmt.Println("========== HighestRandomWeight ==========")
	runTest("Rendezvous", HighestRandomWeight)

	fmt.Println("============== MultiProbe ==============")
	for vNodes := 1; vNodes < 30; vNodes += 5 {
		runTest(fmt.Sprintf("MultiProbe(%d)", vNodes), func() []int { return MultiProbeHashing(vNodes) })
	}
}

func MultiMulti(vNodeMultiplier, leaseMultiplier int) []int {
	r := rand.Int63()
	nodeCounts := make([]int, numNodes)

	var nodes []circularHashEntry
	for node := 0; node < numNodes; node++ {
		for v := 0; v < vNodeMultiplier; v++ {
			value := md5.Sum([]byte(fmt.Sprintf("node %d-%d-%d", r, node, v)))
			str := string(value[:])
			nodes = append(nodes, circularHashEntry{str, node})
		}
	}
	sort.Slice(nodes, func(i, j int) bool {
		return nodes[i].hashValue < nodes[j].hashValue
	})
	for lease := 0; lease < numLeases; lease++ {
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

func MultiProbeHashing(vNodeMultiplier int) []int {
	r := rand.Int63()
	nodeCounts := make([]int, numNodes)

	var nodes []circularHashEntry
	for node := 0; node < numNodes; node++ {
		value := md5.Sum([]byte(fmt.Sprintf("node %d-%d", r, node)))
		str := string(value[:])
		nodes = append(nodes, circularHashEntry{str, node})
	}
	sort.Slice(nodes, func(i, j int) bool {
		return nodes[i].hashValue < nodes[j].hashValue
	})
	for lease := 0; lease < numLeases; lease++ {
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

func HighestRandomWeight() []int {
	r := rand.Int63()
	nodeCounts := make([]int, numNodes)
	for lease := 0; lease < numLeases; lease++ {
		m := ""
		n := -1
		for node := 0; node < numNodes; node++ {
			value := md5.Sum([]byte(fmt.Sprintf("%d-%d-%d", r, lease, node)))
			str := string(value[:])
			if str > m {
				m = str
				n = node
			}
		}
		nodeCounts[n]++
	}
	return nodeCounts
}

func RingHash(vNodeMultiplier int) []int {
	r := rand.Int63()
	var nodes []circularHashEntry
	for node := 0; node < numNodes; node++ {
		for v := 0; v < vNodeMultiplier; v++ {
			value := md5.Sum([]byte(fmt.Sprintf("node %d-%d-%d", r, node, v)))
			str := string(value[:])
			nodes = append(nodes, circularHashEntry{str, node})
		}
	}
	sort.Slice(nodes, func(i, j int) bool {
		return nodes[i].hashValue < nodes[j].hashValue
	})
	nodeCounts := [numNodes]int{}
	leases := []string{}
	for lease := 0; lease < numLeases; lease++ {
		value := md5.Sum([]byte(fmt.Sprintf("lease %d-%d", r, lease)))
		str := string(value[:])
		nodeCounts[getCircularHashOwner(nodes, str).node]++
		leases = append(leases, fmt.Sprintf("%x", str))
	}
	return nodeCounts[:]
}

func runTest(name string, f func() []int) {
	maxV := 0.0
	totalV := 0.0
	totalSpread := 0
	maxSpread := 0
	lowestNodeCount := numLeases
	highestNodeCount := 0
	for run := 0; run < runs; run++ {
		nodeCounts := f()
		vrnce := variance(nodeCounts[:])
		if vrnce > maxV {
			maxV = vrnce
		}
		totalV += vrnce
		thisRunSpread := spread(nodeCounts)
		totalSpread += thisRunSpread
		if maxSpread < thisRunSpread {
			maxSpread = thisRunSpread
		}
		for _, count := range nodeCounts {
			if count < lowestNodeCount {
				lowestNodeCount = count
			}
			if count > highestNodeCount {
				highestNodeCount = count
			}
		}
	}
	fmt.Printf("Test: %s\t- Variance: max - %0.2f avg - %0.2f. Spread: max - %v avg %v\n", name, maxV, totalV/runs, maxSpread, totalSpread/runs)
	fmt.Printf("\tLowest node count: %v. Higest node count: %v. Diff: %v\n", lowestNodeCount, highestNodeCount, highestNodeCount-lowestNodeCount)
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

func variance(data []int) float64 {
	sum := 0
	for _, d := range data {
		sum += d
	}
	average := float64(sum) / float64(len(data))
	variance := 0.0
	for _, d := range data {
		variance += (float64(d) - average) * (float64(d) - average)
	}
	variance /= float64(len(data))
	return variance
}

func spread(data []int) int {
	min := data[0]
	max := data[0]
	for _, d := range data {
		if d < min {
			min = d
		}
		if d > max {
			max = d
		}
	}
	return max - min
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
