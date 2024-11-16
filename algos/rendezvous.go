package algos

import (
	"math/rand"
)

func HighestRandomWeight() []int {
	r := rand.Int63()
	nodeCounts := make([]int, NumNodes)
	for lease := 0; lease < NumLeases; lease++ {
		m := uint64(0)
		n := -1
		for node := 0; node < NumNodes; node++ {
			value := HashBinaryKey("rendezvous", r, lease, node)
			if value > m {
				m = value
				n = node
			}
		}
		nodeCounts[n]++
	}
	return nodeCounts
}
