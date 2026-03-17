package chess

import "math/rand"

func randomUint64() uint64 {
	return rand.Uint64()
}
func randomUint64Sparse() uint64 {
	return randomUint64() & randomUint64() & randomUint64()
}
