package chess

import "math/rand"

// randomUint64 generates a random unsigned integer
func randomUint64() uint64 {
	return rand.Uint64()
}

// randomUint64Sparse generates a random unsigned integer with less no of ones in its binary representation
func randomUint64Sparse() uint64 {
	return randomUint64() & randomUint64() & randomUint64()
}
