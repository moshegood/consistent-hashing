package algos

import (
	"crypto/md5"
	"encoding/binary"
	"unsafe"
)

var NumNodes = 10
var NumLeases = 1000

var HashingFunction func(string) uint64

func md5HashingFunction(s string) uint64 {
	hasher := md5.New()
	hasher.Write([]byte(s))
	bytes := hasher.Sum(nil)
	return binary.BigEndian.Uint64(bytes[:8])
}

const is64Bit = uint64(^uintptr(0)) == ^uint64(0)

//go:noescape
//go:linkname strhash runtime.strhash
func strhash(a unsafe.Pointer, h uintptr) uintptr

// This is the wyhash function used internally by the go runtime
func internalHashingFunction(s string) uint64 {
	return uint64(strhash(unsafe.Pointer(&s), 0))
}

func init() {
	if is64Bit && false {
		HashingFunction = internalHashingFunction
	} else {
		HashingFunction = md5HashingFunction
	}
}
