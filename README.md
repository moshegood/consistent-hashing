# consistent-hashing
Testing consistent hashing distributions.

The correct way to analyze these things is to do the math. But running simulations is fun too.

To see the results, run `go test -v ./...`


## Example output
```
$ go test -v -count=1 ./...
=== RUN   TestRingHash
    algo_test.go:19: =============== RingHash ===============
    algo_test.go:66: Test: RingHash(1)	- Variance: max - 28629.60 avg - 7651.17. Spread: max - 588 avg 270
    algo_test.go:67: 	Lowest node count: 0. Higest node count: 590. Diff: 590
    algo_test.go:66: Test: RingHash(3)	- Variance: max - 11607.80 avg - 2848.04. Spread: max - 382 avg 170
    algo_test.go:67: 	Lowest node count: 5. Higest node count: 411. Diff: 406
    algo_test.go:66: Test: RingHash(9)	- Variance: max - 3954.00 avg - 1065.15. Spread: max - 225 avg 104
    algo_test.go:67: 	Lowest node count: 26. Higest node count: 282. Diff: 256
    algo_test.go:66: Test: RingHash(27)	- Variance: max - 1152.20 avg - 458.89. Spread: max - 124 avg 69
    algo_test.go:67: 	Lowest node count: 39. Higest node count: 182. Diff: 143
    algo_test.go:66: Test: RingHash(81)	- Variance: max - 557.60 avg - 200.96. Spread: max - 80 avg 46
    algo_test.go:67: 	Lowest node count: 53. Higest node count: 157. Diff: 104
    algo_test.go:66: Test: RingHash(243)	- Variance: max - 334.00 avg - 127.12. Spread: max - 61 avg 36
    algo_test.go:67: 	Lowest node count: 62. Higest node count: 136. Diff: 74
    algo_test.go:66: Test: RingHash(729)	- Variance: max - 253.60 avg - 100.87. Spread: max - 62 avg 32
    algo_test.go:67: 	Lowest node count: 69. Higest node count: 140. Diff: 71
--- PASS: TestRingHash (1.50s)
=== RUN   TestHighestRandomWeight
    algo_test.go:26: ========== HighestRandomWeight ==========
    algo_test.go:66: Test: Rendezvous	- Variance: max - 243.60 avg - 80.80. Spread: max - 50 avg 29
    algo_test.go:67: 	Lowest node count: 71. Higest node count: 130. Diff: 59
--- PASS: TestHighestRandomWeight (0.53s)
=== RUN   TestMultiProbe
    algo_test.go:32: ============== MultiProbe ==============
    algo_test.go:66: Test: MultiProbe(1)	- Variance: max - 30690.40 avg - 8378.42. Spread: max - 615 avg 283
    algo_test.go:67: 	Lowest node count: 0. Higest node count: 617. Diff: 617
    algo_test.go:66: Test: MultiProbe(6)	- Variance: max - 4660.20 avg - 1004.65. Spread: max - 170 avg 92
    algo_test.go:67: 	Lowest node count: 0. Higest node count: 180. Diff: 180
    algo_test.go:66: Test: MultiProbe(11)	- Variance: max - 2799.60 avg - 541.56. Spread: max - 146 avg 70
    algo_test.go:67: 	Lowest node count: 1. Higest node count: 167. Diff: 166
    algo_test.go:66: Test: MultiProbe(16)	- Variance: max - 2087.80 avg - 411.50. Spread: max - 139 avg 61
    algo_test.go:67: 	Lowest node count: 0. Higest node count: 155. Diff: 155
    algo_test.go:66: Test: MultiProbe(21)	- Variance: max - 2231.60 avg - 362.56. Spread: max - 139 avg 57
    algo_test.go:67: 	Lowest node count: 1. Higest node count: 146. Diff: 145
    algo_test.go:66: Test: MultiProbe(26)	- Variance: max - 1637.80 avg - 257.01. Spread: max - 134 avg 49
    algo_test.go:67: 	Lowest node count: 1. Higest node count: 144. Diff: 143
--- PASS: TestMultiProbe (4.42s)
PASS
```
