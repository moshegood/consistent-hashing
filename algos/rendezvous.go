package algos

import (
	"fmt"
	"math/rand"
)

func HighestRandomWeight() []int {
	r := rand.Int63()
	nodeCounts := make([]int, NumNodes)
	for lease := 0; lease < NumLeases; lease++ {
		max := 0
		n := -1
		for node := 0; node < NumNodes; node++ {
			Hasher.Reset()
			Hasher.Write([]byte(fmt.Sprintf("rendezvous %d-%d-%d", r, lease, node)))
			value := Hasher.Sum64()

			if value > m {
				max = value
				n = node
			}
		}
		nodeCounts[n]++
	}
	return nodeCounts
}
