package algos

import (
	"hash"
	"hash/fnv"
)

var NumNodes = 10
var NumLeases = 1000

var Hasher hash.Hash64 = fnv.New64a()
