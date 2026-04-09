package chess

var PawnAttacks [2][64]uint64 // PawnAttacks stores the pawn attack bitmasks for each square

// maskPawnAttacks generates pawn attacks for a square
func maskPawnAttacks(color uint8, sq uint8) uint64 {
	var bitboard uint64 = 0
	var attacks uint64 = 0
	SetBit(&bitboard, sq)
	if color == White {
		//NE
		if (bitboard & NotHFile) != 0 {
			attacks |= bitboard << 9
		}
		//NW
		if (bitboard & NotAFile) != 0 {
			attacks |= bitboard << 7
		}
	} else {
		//SE
		if (bitboard & NotHFile) != 0 {
			attacks |= bitboard >> 7
		}
		//SW
		if (bitboard & NotAFile) != 0 {
			attacks |= bitboard >> 9
		}
	}
	return attacks
}

// generatePawnMoves generates pawn attack moves
func (board *Board) generatePawnMoves(ml *MoveList) {
	color := board.SideToMove
	pawn := board.Colors[color] & board.Pieces[Pawn]
	occupiedMask := ^(board.Colors[White] | board.Colors[Black])

	if color == White {
		// Handles en passant
		if board.EnPassantSquare != 255 {
			var enPassant uint64 = 1 << board.EnPassantSquare
			leftAttackers := (enPassant >> 9) & pawn & NotHFile
			rightAttackers := (enPassant >> 7) & pawn & NotAFile
			if leftAttackers != 0 {
				ml.Add(NewMove(board.EnPassantSquare-9, board.EnPassantSquare, 5))
			}
			if rightAttackers != 0 {
				ml.Add(NewMove(board.EnPassantSquare-7, board.EnPassantSquare, 5))
			}
		}
		attack1 := (pawn << 8) & occupiedMask
		attack2 := ((attack1 & 0x0000000000FF0000) << 8) & occupiedMask
		// One move
		for attack1 != 0 {
			nextSq := PopBit(&attack1)
			if nextSq >= 56 {
				ml.Add(NewMove(nextSq-8, nextSq, 8))
				ml.Add(NewMove(nextSq-8, nextSq, 9))
				ml.Add(NewMove(nextSq-8, nextSq, 10))
				ml.Add(NewMove(nextSq-8, nextSq, 11))
			} else {
				ml.Add(NewMove(nextSq-8, nextSq, 0))
			}
		}
		// Two moves
		for attack2 != 0 {
			nextSq := PopBit(&attack2)
			ml.Add(NewMove(nextSq-16, nextSq, 1))
		}
		//NE capture
		attacks := (pawn << 9) & NotAFile & board.Colors[Black]
		for attacks != 0 {
			nextSq := PopBit(&attacks)
			if nextSq >= 56 {
				ml.Add(NewMove(nextSq-9, nextSq, 12))
				ml.Add(NewMove(nextSq-9, nextSq, 13))
				ml.Add(NewMove(nextSq-9, nextSq, 14))
				ml.Add(NewMove(nextSq-9, nextSq, 15))
			} else {
				ml.Add(NewMove(nextSq-9, nextSq, 4))
			}

		}
		//NW capture
		attacks = (pawn << 7) & NotHFile & board.Colors[Black]
		for attacks != 0 {
			nextSq := PopBit(&attacks)
			if nextSq >= 56 {
				ml.Add(NewMove(nextSq-7, nextSq, 12))
				ml.Add(NewMove(nextSq-7, nextSq, 13))
				ml.Add(NewMove(nextSq-7, nextSq, 14))
				ml.Add(NewMove(nextSq-7, nextSq, 15))
			} else {
				ml.Add(NewMove(nextSq-7, nextSq, 4))
			}
		}
	} else if color == Black {
		// Handles en passant
		if board.EnPassantSquare != 255 {
			var enPassant uint64 = 1 << board.EnPassantSquare
			leftAttackers := (enPassant << 7) & pawn & NotHFile
			rightAttackers := (enPassant << 9) & pawn & NotAFile
			if leftAttackers != 0 {
				ml.Add(NewMove(board.EnPassantSquare+7, board.EnPassantSquare, 5))
			}
			if rightAttackers != 0 {
				ml.Add(NewMove(board.EnPassantSquare+9, board.EnPassantSquare, 5))
			}
		}
		attack1 := (pawn >> 8) & occupiedMask
		attack2 := ((attack1 & 0x0000FF0000000000) >> 8) & occupiedMask
		// One move
		for attack1 != 0 {
			nextSq := PopBit(&attack1)
			if nextSq <= 7 {
				ml.Add(NewMove(nextSq+8, nextSq, 8))
				ml.Add(NewMove(nextSq+8, nextSq, 9))
				ml.Add(NewMove(nextSq+8, nextSq, 10))
				ml.Add(NewMove(nextSq+8, nextSq, 11))
			} else {
				ml.Add(NewMove(nextSq+8, nextSq, 0))
			}
		}
		// Two moves
		for attack2 != 0 {
			nextSq := PopBit(&attack2)
			ml.Add(NewMove(nextSq+16, nextSq, 1))
		}
		//SW capture
		attacks := (pawn >> 9) & NotHFile & board.Colors[White]
		for attacks != 0 {
			nextSq := PopBit(&attacks)
			if nextSq <= 7 {
				ml.Add(NewMove(nextSq+9, nextSq, 12))
				ml.Add(NewMove(nextSq+9, nextSq, 13))
				ml.Add(NewMove(nextSq+9, nextSq, 14))
				ml.Add(NewMove(nextSq+9, nextSq, 15))
			} else {
				ml.Add(NewMove(nextSq+9, nextSq, 4))
			}
		}
		//SE capture
		attacks = (pawn >> 7) & NotAFile & board.Colors[White]
		for attacks != 0 {
			nextSq := PopBit(&attacks)
			if nextSq <= 7 {
				ml.Add(NewMove(nextSq+7, nextSq, 12))
				ml.Add(NewMove(nextSq+7, nextSq, 13))
				ml.Add(NewMove(nextSq+7, nextSq, 14))
				ml.Add(NewMove(nextSq+7, nextSq, 15))
			} else {
				ml.Add(NewMove(nextSq+7, nextSq, 4))
			}
		}
	}
}
