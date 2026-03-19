package chess

var KingAttacks [64]uint64

func maskKingAttacks(sq uint8) uint64 {
	var bitboard uint64 = 0
	var attacks uint64 = 0
	SetBit(&bitboard, sq)
	//N
	attacks |= bitboard << 8
	//S
	attacks |= bitboard >> 8
	//E
	if (bitboard & NotHFile) != 0 {
		attacks |= bitboard << 1
	}
	//W
	if (bitboard & NotAFile) != 0 {
		attacks |= bitboard >> 1
	}
	//NW
	if (bitboard & NotAFile) != 0 {
		attacks |= bitboard << 7
	}
	//NE
	if (bitboard & NotHFile) != 0 {
		attacks |= bitboard << 9
	}
	//SE
	if (bitboard & NotHFile) != 0 {
		attacks |= bitboard >> 7
	}
	//SW
	if (bitboard & NotAFile) != 0 {
		attacks |= bitboard >> 9
	}
	return attacks
}
func (board *Board) generateKingMoves(ml *MoveList) {
	color := board.SideToMove
	king := board.Colors[color] & board.Pieces[King]
	castling := board.CastlingRights
	occupied := board.Colors[White] | board.Colors[Black]
	enemyColor := color ^ 1
	enemyPieces := board.Colors[enemyColor]
	for king != 0 {
		currSq := PopBit(&king)
		attacks := KingAttacks[currSq] &^ board.Colors[color]
		captures := attacks & enemyPieces
		quiets := attacks &^ enemyPieces
		for captures != 0 {
			nextSq := PopBit(&captures)
			ml.Add(NewMove(currSq, nextSq, 4))
		}
		for quiets != 0 {
			nextSq := PopBit(&quiets)
			ml.Add(NewMove(currSq, nextSq, 0))
		}
	}
	if color == White {
		if castling&WhiteKingside != 0 {
			if occupied&((1<<5)|(1<<6)) == 0 {
				ml.Add(NewMove(uint8(4), uint8(6), 2))
			}
		}
		if castling&WhiteQueenside != 0 {
			if occupied&((1<<1)|(1<<2)|(1<<3)) == 0 {
				ml.Add(NewMove(uint8(4), uint8(2), 3))
			}
		}
	} else if color == Black {
		if castling&BlackKingside != 0 {
			if occupied&((1<<61)|(1<<62)) == 0 {
				ml.Add(NewMove(uint8(60), uint8(62), 2))
			}
		}
		if castling&BlackQueenside != 0 {
			if occupied&((1<<57)|(1<<58)|(1<<59)) == 0 {
				ml.Add(NewMove(uint8(60), uint8(58), 3))
			}
		}
	}
}
