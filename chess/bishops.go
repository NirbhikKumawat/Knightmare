package chess

import "math/bits"

var BishopMagics = [64]uint64{
	0x8210828208020030,
	0x4880204002000,
	0x44011612000048,
	0xc4042088000400,
	0x182061040000000,
	0x62042a2010200400,
	0x802840928400040,
	0x501304084c020850,
	0x100008100400840d,
	0xc00101002088021,
	0x8021040408820400,
	0x8040d02080901,
	0x20031040040814,
	0x818030420048228,
	0x201128410881400,
	0x2240c04100882002,
	0x84a00004a0220240,
	0x14460204840403,
	0x1002202020a00,
	0x40803802094004,
	0x8000802400a00082,
	0x1004090201040220,
	0x1000848402013000,
	0x1002200500880402,
	0x8282400050108625,
	0x801088404080800,
	0x5000208010018480,
	0x400c01102401000a,
	0x401010010104002,
	0x130400201100a011,
	0x1041408441400,
	0x951c008000260100,
	0x802203000061081,
	0x202824820101000,
	0x10102808040b0204,
	0x106405800408200,
	0xc00a004200040208,
	0x4e20040420010090,
	0x102081101046402,
	0x2040c1180846082,
	0xa488140a88402000,
	0x1441004008948,
	0x8145011806000404,
	0x40104a0a00,
	0x1002084104000040,
	0x1013101000200,
	0x40200401020b10d2,
	0x4808010242000680,
	0x4004848430400085,
	0xc404840088044800,
	0x220211180900080,
	0x3080042120120,
	0x4001044098a20001,
	0x80010c4408420000,
	0x2140c81808ac8000,
	0x8082084004020,
	0x10148018c5004,
	0x8084056213101804,
	0x400560e080480820,
	0x180001a8ca08820,
	0x200214090220200,
	0x420104002040100,
	0x28a0710c4501,
	0x1006060202140500,
}

var BishopMasks [64]uint64
var BishopAttacks [64][512]uint64

func maskBishopOccupancy(sq uint8) uint64 {
	var mask uint64 = 0
	targetRank := int(sq / 8)
	targetFile := int(sq % 8)
	//SE
	for r, f := targetRank-1, targetFile+1; r >= 1 && f <= 6; r, f = r-1, f+1 {
		SetBit(&mask, uint8(r*8+f))
	}
	//NE
	for r, f := targetRank+1, targetFile+1; r <= 6 && f <= 6; r, f = r+1, f+1 {
		SetBit(&mask, uint8(r*8+f))
	}
	//NW
	for r, f := targetRank+1, targetFile-1; r <= 6 && f >= 1; r, f = r+1, f-1 {
		SetBit(&mask, uint8(r*8+f))
	}
	//SW
	for r, f := targetRank-1, targetFile-1; r >= 1 && f >= 1; r, f = r-1, f-1 {
		SetBit(&mask, uint8(r*8+f))
	}
	return mask
}

func bishopAttacksOnTheFly(sq uint8, block uint64) uint64 {
	var attacks uint64 = 0
	targetRank := int(sq / 8)
	targetFile := int(sq % 8)
	//SE
	for r, f := targetRank-1, targetFile+1; r >= 0 && f <= 7; r, f = r-1, f+1 {
		square := uint8(r*8 + f)
		SetBit(&attacks, square)
		if (block & (1 << square)) != 0 {
			break
		}
	}
	//NE
	for r, f := targetRank+1, targetFile+1; r <= 7 && f <= 7; r, f = r+1, f+1 {
		square := uint8(r*8 + f)
		SetBit(&attacks, square)
		if (block & (1 << square)) != 0 {
			break
		}
	}
	//NW
	for r, f := targetRank+1, targetFile-1; r <= 7 && f >= 0; r, f = r+1, f-1 {
		square := uint8(r*8 + f)
		SetBit(&attacks, square)
		if (block & (1 << square)) != 0 {
			break
		}
	}
	//SW
	for r, f := targetRank-1, targetFile-1; r >= 0 && f >= 0; r, f = r-1, f-1 {
		square := uint8(r*8 + f)
		SetBit(&attacks, square)
		if (block & (1 << square)) != 0 {
			break
		}
	}

	return attacks
}

func GetBishopAttacks(sq uint8, occupancy uint64) uint64 {
	blockers := occupancy & BishopMasks[sq]
	magicIndex := (blockers * BishopMagics[sq]) >> (64 - bits.OnesCount64(BishopMasks[sq]))
	return BishopAttacks[sq][magicIndex]
}
