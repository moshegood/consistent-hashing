package algos

import (
	"encoding/binary"
	"hash"
	"hash/fnv"
)

var NumNodes = 10
var NumLeases = 1000

var Hasher hash.Hash64 = fnv.New64a()

type md5Hasher struct {
	hash.Hash
}

func (h *md5Hasher) Sum64() uint64 {
	bytes := h.Sum(nil)
	return binary.BigEndian.Uint64(bytes[:8])
}
