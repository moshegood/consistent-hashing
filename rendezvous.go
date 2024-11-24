package algos

import (
	"fmt"
	"math/rand"
)

func HighestRandomWeight() []int {
	r := rand.Int63()
	nodeCounts := make([]int, NumNodes)
	for lease := 0; lease < NumLeases; lease++ {
		var max uint64 = 0
		n := -1
		for node := 0; node < NumNodes; node++ {
			value := HashingFunction(fmt.Sprintf("rendezvous %d-%d-%d", r, lease, node))

			if value > max {
				max = value
				n = node
			}
		}
		nodeCounts[n]++
	}
	return nodeCounts
}
