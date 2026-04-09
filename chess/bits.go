package chess

import "math/bits"

// SetBit sets the bit at specified square to 1
func SetBit(bb *uint64, square uint8) {
	var k uint64
	k = 1 << square
	*bb |= k
}

// ClearBit sets the bit at specified square to 0
func ClearBit(bb *uint64, square uint8) {
	var k uint64
	k = 1 << square
	*bb &= ^k
}

// GetBit returns the bit at given square
func GetBit(bb uint64, square uint8) (k uint64) {
	k = 1 << square
	k &= bb
	return
}

// PopBit clears the least significant bit and returns its position,used for iterating over a bitboard
func PopBit(bb *uint64) uint8 {
	sq := uint8(bits.TrailingZeros64(*bb))
	*bb &= *bb - 1
	return sq
}
