package algos

import (
	"crypto/md5"
	"fmt"
	"math/rand"
)

func HighestRandomWeight() []int {
	r := rand.Int63()
	nodeCounts := make([]int, NumNodes)
	for lease := 0; lease < NumLeases; lease++ {
		m := ""
		n := -1
		for node := 0; node < NumNodes; node++ {
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
