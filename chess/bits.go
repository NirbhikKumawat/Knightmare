package chess

import "math/bits"

func SetBit(bb *uint64, square uint8) {
	var k uint64
	k = 1 << square
	*bb |= k
}
func ClearBit(bb *uint64, square uint8) {
	var k uint64
	k = 1 << square
	*bb &= ^k
}
func GetBit(bb uint64, square uint8) (k uint64) {
	k = 1 << square
	k &= bb
	return
}
func PopBit(bb *uint64) uint8 {
	sq := uint8(bits.TrailingZeros64(*bb))
	*bb &= *bb - 1
	return sq
}
