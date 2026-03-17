package chess

import (
	"fmt"
	"math/bits"
)

func findMagicBishop(sq uint8, relevantBits int) uint64 {
	perms := 1 << relevantBits
	occupancies := make([]uint64, perms)
	attacks := make([]uint64, perms)
	mask := maskBishopOccupancy(sq)
	for i := 0; i < perms; i++ {
		occupancies[i] = setOccupancy(i, relevantBits, mask)
		attacks[i] = bishopAttacksOnTheFly(sq, occupancies[i])
	}
	usedAttacks := make([]uint64, perms)
	for {
		magic := randomUint64Sparse()

		if bits.OnesCount64(mask*magic) < 6 {
			continue
		}

		for i := range usedAttacks {
			usedAttacks[i] = 0
		}
		fail := false
		for i := 0; i < perms; i++ {
			magicIndex := (occupancies[i] * magic) >> (64 - relevantBits)
			if usedAttacks[magicIndex] == 0 {
				usedAttacks[magicIndex] = attacks[i]
			} else if usedAttacks[magicIndex] != attacks[i] {
				fail = true
				break
			}
		}
		if !fail {
			return magic
		}
	}
}
func findMagicRook(sq uint8, relevantBits int) uint64 {
	perms := 1 << relevantBits
	occupancies := make([]uint64, perms)
	attacks := make([]uint64, perms)
	mask := maskRookOccupancy(sq)
	for i := 0; i < perms; i++ {
		occupancies[i] = setOccupancy(i, relevantBits, mask)
		attacks[i] = rookAttacksOnTheFly(sq, occupancies[i])
	}
	usedAttacks := make([]uint64, perms)
	for {
		magic := randomUint64Sparse()

		if bits.OnesCount64(mask*magic) < 6 {
			continue
		}

		for i := range usedAttacks {
			usedAttacks[i] = 0
		}

		fail := false
		for i := 0; i < perms; i++ {
			magicIndex := (occupancies[i] * magic) >> (64 - relevantBits)
			if usedAttacks[magicIndex] == 0 {
				usedAttacks[magicIndex] = attacks[i]
			} else if usedAttacks[magicIndex] != attacks[i] {
				fail = true
				break
			}
		}
		if !fail {
			return magic
		}
	}
}
func GenerateAllMagics() {
	fmt.Println("var BishopMagics = [64]uint64{")
	for sq := 0; sq < 64; sq++ {
		relevantBits := bits.OnesCount64(maskBishopOccupancy(uint8(sq)))
		magic := findMagicBishop(uint8(sq), relevantBits)
		fmt.Printf("\t0x%x,\n", magic)
	}
	fmt.Println("}")

	fmt.Println("\nvar RookMagics = [64]uint64{")
	for sq := 0; sq < 64; sq++ {
		relevantBits := bits.OnesCount64(maskRookOccupancy(uint8(sq)))
		magic := findMagicRook(uint8(sq), relevantBits)
		fmt.Printf("\t0x%x,\n", magic)
	}
	fmt.Println("}")
}
