package algos

func MultiProbeHashing(vNodeMultiplier int) []int {
	ring := MakeRing(NumNodes, 1)
	nodeCounts := make([]int, NumNodes)
	for lease := 0; lease < NumLeases; lease++ {
		owner := FindOwnerOfLeaseInRing(ring, lease, vNodeMultiplier)
		nodeCounts[owner]++
	}
	return nodeCounts
}
