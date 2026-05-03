package algos

import (
	"fmt"
	"hash/fnv"
	"math/rand"
	"testing"
	"time"

	"github.com/cespare/xxhash/v2"
	"github.com/orisano/wyhash"
	rapidhash "github.com/vkudryk/rapidhash-go"
)

const runs = 200

func init() {
	rand.Seed(time.Now().UnixNano())
}

func TestRingHash(t *testing.T) {
	t.Log("=============== RingHash ===============")
	for vNodes := 1; vNodes < 2000; vNodes *= 4 {
		runTest(t, fmt.Sprintf("RingHash(%d)", vNodes), func() []int { return RingHash(vNodes) })
	}
}

func TestHighestRandomWeight(t *testing.T) {
	t.Log("========== HighestRandomWeight ==========")
	runTest(t, "Rendezvous", HighestRandomWeight)

}

func TestMultiProbe(t *testing.T) {
	t.Log("============== MultiProbe ==============")
	for vNodes := 1; vNodes < 30; vNodes += 5 {
		runTest(t, fmt.Sprintf("MultiProbe(%d)", vNodes), func() []int { return MultiProbeHashing(vNodes) })
	}
}

func BenchmarkHashingFunctions(b *testing.B) {
	input := "benchmark-input-string-for-hashing-performance-comparison"
	b_input := []byte(input)

	b.Run("InternalUnsafe", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = internalHashingFunction(input)
		}
	})

	b.Run("MD5Slow", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = md5HashingFunction(input)
		}
	})

	b.Run("FNVFast", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			h := fnv.New64a()
			h.Write(b_input)
			_ = h.Sum64()
		}
	})
	b.Run("Wyhash", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			h := wyhash.New(0)
			h.Write(b_input)
			_ = h.Sum64()
		}
	})
	b.Run("RapidHash", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = rapidhash.Hash(b_input)
		}
	})

	b.Run("XXHash", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = xxhash.Sum64String(input)
		}
	})
}

func runTest(t *testing.T, name string, f func() []int) {
	maxV := 0.0
	totalV := 0.0
	totalSpread := 0
	maxSpread := 0
	lowestNodeCount := NumLeases
	highestNodeCount := 0
	start := time.Now()
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
	t.Logf("Test: %s", name)
	t.Logf("\tVariance: max = %0.2f avg = %0.2f. Spread: max = %v avg = %v\n", maxV, totalV/runs, maxSpread, totalSpread/runs)
	t.Logf("\tLowest node count: %v. Higest node count: %v. Diff: %v\n", lowestNodeCount, highestNodeCount, highestNodeCount-lowestNodeCount)
	t.Logf("\tTime:%v", time.Since(start))
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
