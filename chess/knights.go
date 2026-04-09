package chess

var KnightAttacks [64]uint64 // KnightAttacks stores the knight attack bitmasks for each square

// maskKnightAttacks generates knight attacks for a square
func maskKnightAttacks(sq uint8) uint64 {
	var bitboard uint64 = 0
	var attacks uint64 = 0
	SetBit(&bitboard, sq)
	//NNE
	if (bitboard & NotHFile) != 0 {
		attacks |= bitboard << 17
	}
	//NEE
	if (bitboard & NotGHFile) != 0 {
		attacks |= bitboard << 10
	}
	//NNW
	if (bitboard & NotAFile) != 0 {
		attacks |= bitboard << 15
	}
	//NWW
	if (bitboard & NotABFile) != 0 {
		attacks |= bitboard << 6
	}
	//SSE
	if (bitboard & NotHFile) != 0 {
		attacks |= bitboard >> 15
	}
	//SEE
	if (bitboard & NotGHFile) != 0 {
		attacks |= bitboard >> 6
	}
	//SSW
	if (bitboard & NotAFile) != 0 {
		attacks |= bitboard >> 17
	}
	//SWW
	if (bitboard & NotABFile) != 0 {
		attacks |= bitboard >> 10
	}
	return attacks
}

// generateKnightMoves generates knight attack moves
func (board *Board) generateKnightMoves(ml *MoveList) {
	color := board.SideToMove
	enemyColor := color ^ 1
	enemyPieces := board.Colors[enemyColor]
	knights := board.Colors[color] & board.Pieces[Knight]
	for knights != 0 {
		currSq := PopBit(&knights)
		attacks := KnightAttacks[currSq] &^ board.Colors[color]
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
}

// generateSliderMoves generates attack moves for queen,bishop and rook
func (board *Board) generateSliderMoves(ml *MoveList, piece uint8) {
	occupied := board.Colors[White] | board.Colors[Black]
	color := board.SideToMove
	enemyColor := color ^ 1
	enemyPieces := board.Colors[enemyColor]
	if piece == Queen {
		queens := board.Colors[color] & board.Pieces[Queen]
		for queens != 0 {
			currSq := PopBit(&queens)
			attacks := (GetRookAttacks(currSq, occupied) | GetBishopAttacks(currSq, occupied)) &^ board.Colors[color]
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
	}
	if piece == Bishop {
		bishops := board.Colors[color] & board.Pieces[Bishop]
		for bishops != 0 {
			currSq := PopBit(&bishops)
			attacks := GetBishopAttacks(currSq, occupied) &^ board.Colors[color]
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
	}
	if piece == Rook {
		rooks := board.Colors[color] & board.Pieces[Rook]
		for rooks != 0 {
			currSq := PopBit(&rooks)
			attacks := GetRookAttacks(currSq, occupied) &^ board.Colors[color]
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
	}
}
