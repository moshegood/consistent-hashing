package algos

import (
	"fmt"

	"github.com/dchest/siphash"
)

const K0 = uint64(316665572293978160)
const K1 = uint64(8573005253291875333)

func HashUnaryKey(label string, slug int64, key int) uint64 {
	return Hash([]byte(fmt.Sprintf("%s:%d-%d", label, slug, key)))
}

func HashBinaryKey(label string, slug int64, key0 int, key1 int) uint64 {
	return Hash([]byte(fmt.Sprintf("%s:%d-%d-%d", label, slug, key0, key1)))
}

func Hash(bytes []byte) uint64 {
	return Siphash(bytes)
}

func Siphash(bytes []byte) uint64 {
	return siphash.Hash(K0, K1, bytes)
}
