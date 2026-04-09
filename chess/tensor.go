package chess

// ToTensor board struct into tensor object
func (board *Board) ToTensor() []float32 {
	tensor := make([]float32, 896) // 14 8x8 bitboards
	// (0-5) for white pieces and (6-11) for black pieces
	for color := uint8(0); color <= 1; color++ {
		for piece := uint8(0); piece <= 5; piece++ {
			plane := int(6*color + piece)
			planeOffset := plane * 64
			bb := board.Colors[color] & board.Pieces[piece]
			for bb != 0 {
				sq := PopBit(&bb)
				tensor[planeOffset+int(sq)] = 1.0
			}
		}
	}
	// (12) to indicate side to move (all 1.0 for white,all 0.0 for black)
	if board.SideToMove == White {
		for i := 12 * 64; i < 13*64; i++ {
			tensor[i] = 1.0
		}
	}
	// (13) indicate special moves(en passant,castling)
	specialOffset := 13 * 64
	if board.CastlingRights&WhiteKingside != 0 {
		tensor[specialOffset+0] = 1.0
	}
	if board.CastlingRights&WhiteQueenside != 0 {
		tensor[specialOffset+1] = 1.0
	}
	if board.CastlingRights&BlackKingside != 0 {
		tensor[specialOffset+2] = 1.0
	}
	if board.CastlingRights&BlackQueenside != 0 {
		tensor[specialOffset+3] = 1.0
	}

	if board.EnPassantSquare != 255 {
		tensor[specialOffset+int(board.EnPassantSquare)] = 1.0
	}

	return tensor
}

// MoveToIndex converts move into integer, modulo 4096 and more than 4096 to indicate promotions
func MoveToIndex(m Move) int {
	from := int(m.From())
	to := int(m.To())
	flags := int(m.Flags())
	promo := (flags >= 8 && flags <= 10) || (flags >= 12 && flags <= 14)
	move := (from * 64) + to
	if promo {
		switch flags {
		case 8, 12:
			return move + 4096
		case 9, 13:
			return move + 4096 + 4096
		case 10, 14:
			return move + 4096 + 4096 + 4096
		}
	}
	return (from * 64) + to
}

// IndexToMove converts index to move
func IndexToMove(index int, board *Board) Move {
	baseIndex := index % 4096
	from := uint8(baseIndex / 64)
	to := uint8(baseIndex % 64)

	movingPiece := board.GetPieceType(from)
	targetPiece := board.GetPieceType(to)
	isCapture := targetPiece != Empty

	var flags uint16 = 0

	if movingPiece == Pawn {
		isPromo := (to >= 56) || (to <= 7)

		if isPromo {
			baseFlag := uint16(11)
			if index >= 12288 {
				baseFlag = 10 // Rook
			} else if index >= 8192 {
				baseFlag = 9 // Bishop
			} else if index >= 4096 {
				baseFlag = 8 // Knight
			}

			if isCapture {
				baseFlag += 4
			}
			flags = baseFlag
		} else if isCapture {
			flags = 4
		} else if from%8 != to%8 {
			flags = 5
		} else if int(to)-int(from) == 16 || int(from)-int(to) == 16 {
			flags = 1
		} else {
			flags = 0
		}

	} else if movingPiece == King {
		if from == 4 && to == 6 {
			flags = 2
		} else if from == 4 && to == 2 {
			flags = 3
		} else if from == 60 && to == 62 {
			flags = 2
		} else if from == 60 && to == 58 {
			flags = 3
		} else if isCapture {
			flags = 4
		}

	} else {

		if isCapture {
			flags = 4
		}
	}

	return NewMove(from, to, flags)
}
