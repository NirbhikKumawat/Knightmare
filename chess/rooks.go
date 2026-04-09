package chess

import "math/bits"

// RookMagics is an array of magic numbers for getting a rook attack
var RookMagics = [64]uint64{
	0x80006018400080,
	0x40004620001000,
	0x80088010042000,
	0x300089001002065,
	0x2200082014020050,
	0x200431200081004,
	0x4a00049916000448,
	0x80002180094100,
	0x1028000c00080a2,
	0x8000808040002000,
	0x1100110240a000,
	0x2808008003000,
	0x8002800800040080,
	0xba000410080200,
	0x801004405000200,
	0x120e001401008a42,
	0x880800824c001,
	0xe404000201000,
	0x408260010408200,
	0x40210010010900,
	0x2031000d080100,
	0x8088014020100c,
	0x180808002000900,
	0x4059020010408401,
	0x1980018880204000,
	0x3010200240100140,
	0x8210429100200102,
	0x2028b00100090020,
	0x800040080080080,
	0x8042000200183044,
	0x400480c001019a2,
	0x900881220000d084,
	0x4012804004800120,
	0x10600082804000,
	0x1802202004050,
	0x4002500580800800,
	0xa90240080800800,
	0x8c00800200800400,
	0x10820801040002d0,
	0x884082001104,
	0xa0400080388001,
	0x4020002150084000,
	0x4005108142020022,
	0xc003300100a10008,
	0x408010044290010,
	0x840002008080,
	0x21420084010a0008,
	0x1005410820003,
	0x2040800420400080,
	0x40072000500040,
	0x4000420280302200,
	0x2080180090008080,
	0x2020010c8204e00,
	0x2002008004000280,
	0x110228d002010400,
	0x140010082440200,
	0x108000410021,
	0x203012330820046,
	0x6000314100a009,
	0x422000cb8102042,
	0x800200042008100a,
	0x2581000208840011,
	0x8401000400820001,
	0x508a02481004c02,
}

var RookMasks [64]uint64         // RookMasks generate positions at which pieces can block rooks attacks
var RookAttacks [64][4096]uint64 // RookAttacks stores the possible rook attacks from a given position

// initSliders fills up RookMasks, RookAttacks, BishopMasks and BishopAttacks using magic numbers
func initSliders() {
	for sq := 0; sq < 64; sq++ {
		BishopMasks[sq] = maskBishopOccupancy(uint8(sq))
		RookMasks[sq] = maskRookOccupancy(uint8(sq))
		bishopRelevantBits := bits.OnesCount64(BishopMasks[sq])
		rookRelevantBits := bits.OnesCount64(RookMasks[sq])
		bishopOccupancyIndices := 1 << bishopRelevantBits
		rookOccupancyIndices := 1 << rookRelevantBits
		for i := 0; i < bishopOccupancyIndices; i++ {
			occupancy := setOccupancy(i, bishopRelevantBits, BishopMasks[sq])
			attacks := bishopAttacksOnTheFly(uint8(sq), occupancy)
			magicIndex := (occupancy * BishopMagics[sq]) >> (64 - bishopRelevantBits)
			BishopAttacks[sq][magicIndex] = attacks
		}
		for i := 0; i < rookOccupancyIndices; i++ {
			occupancy := setOccupancy(i, rookRelevantBits, RookMasks[sq])
			attacks := rookAttacksOnTheFly(uint8(sq), occupancy)
			magicIndex := (occupancy * RookMagics[sq]) >> (64 - rookRelevantBits)
			RookAttacks[sq][magicIndex] = attacks
		}

	}
}

// maskRookOccupancy generates a bitboard marking all squares which can block rooks attack
func maskRookOccupancy(sq uint8) uint64 {
	var mask uint64 = 0
	targetRank := int(sq / 8)
	targetFile := int(sq % 8)
	//N
	for r := targetRank + 1; r <= 6; r++ {
		SetBit(&mask, uint8(r*8+targetFile))
	}
	//E
	for r := targetFile + 1; r <= 6; r++ {
		SetBit(&mask, uint8(targetRank*8+r))
	}
	//S
	for r := 1; r < targetRank; r++ {
		SetBit(&mask, uint8(r*8+targetFile))
	}
	//W
	for r := 1; r < targetFile; r++ {
		SetBit(&mask, uint8(targetRank*8+r))
	}
	return mask
}

// setOccupancy generates all possible blocker permutations for a particular square
// index: the permutation number
// bitsInMask : numbers of bits set to 1 in attackMask
func setOccupancy(index int, bitsInMask int, attackMask uint64) uint64 {
	var occupancy uint64 = 0
	for count := 0; count < bitsInMask; count++ {
		sq := PopBit(&attackMask)
		if (index & (1 << count)) != 0 {
			occupancy |= 1 << sq
		}
	}
	return occupancy
}

// rookAttacksOnTheFly generates rook attacks stopping at blockers
func rookAttacksOnTheFly(sq uint8, block uint64) uint64 {
	var attacks uint64 = 0
	targetRank := int(sq / 8)
	targetFile := int(sq % 8)

	//N
	for r := targetRank + 1; r <= 7; r++ {
		square := uint8(r*8 + targetFile)
		SetBit(&attacks, square)
		if (block & (1 << square)) != 0 {
			break
		}
	}
	//S
	for r := targetRank - 1; r >= 0; r-- {
		square := uint8(r*8 + targetFile)
		SetBit(&attacks, square)
		if (block & (1 << square)) != 0 {
			break
		}
	}
	//E
	for r := targetFile + 1; r <= 7; r++ {
		square := uint8(targetRank*8 + r)
		SetBit(&attacks, square)
		if (block & (1 << square)) != 0 {
			break
		}
	}
	//W
	for r := targetFile - 1; r >= 0; r-- {
		square := uint8(targetRank*8 + r)
		SetBit(&attacks, square)
		if (block & (1 << square)) != 0 {
			break
		}
	}
	return attacks
}

// GetRookAttacks returns all the possible rook attacks from a square
func GetRookAttacks(sq uint8, occupancy uint64) uint64 {
	blockers := occupancy & RookMasks[sq]
	magicIndex := (blockers * RookMagics[sq]) >> (64 - bits.OnesCount64(RookMasks[sq]))
	return RookAttacks[sq][magicIndex]
}
