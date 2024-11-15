package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/moshegood/consistent-hashing/algos"
)

const runs = 200

func main() {
	rand.Seed(time.Now().UnixNano())
	fmt.Println("=============== RingHash ===============")
	for vNodes := 1; vNodes < 1000; vNodes *= 3 {
		runTest(fmt.Sprintf("RingHash(%d)", vNodes), func() []int { return algos.RingHash(vNodes) })
	}

	fmt.Println("========== HighestRandomWeight ==========")
	runTest("Rendezvous", algos.HighestRandomWeight)

	fmt.Println("============== MultiProbe ==============")
	for vNodes := 1; vNodes < 30; vNodes += 5 {
		runTest(fmt.Sprintf("MultiProbe(%d)", vNodes), func() []int { return algos.MultiProbeHashing(vNodes) })
	}
}

func runTest(name string, f func() []int) {
	maxV := 0.0
	totalV := 0.0
	totalSpread := 0
	maxSpread := 0
	lowestNodeCount := algos.NumLeases
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
