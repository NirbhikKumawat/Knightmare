package chess

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func ParseFEN(fen string) (*Board, error) {
	board := &Board{}
	fenParts := strings.Split(fen, " ")
	if len(fenParts) != 6 {
		return nil, errors.New("invalid FEN format")
	}
	pieceRows := strings.Split(fenParts[0], "/")
	ParsePieces(board, pieceRows)
	activeColor := fenParts[1]
	if activeColor == "w" || activeColor == "W" {
		board.SideToMove = White
	} else {
		board.SideToMove = Black
	}
	castlingRights := fenParts[2]
	cr := 0
	for _, rights := range castlingRights {
		if rights == 'Q' {
			cr |= WhiteQueenside
		} else if rights == 'K' {
			cr |= WhiteKingside
		} else if rights == 'q' {
			cr |= BlackQueenside
		} else if rights == 'k' {
			cr |= BlackKingside
		}
	}
	board.CastlingRights = uint8(cr)
	enPassant := fenParts[3]
	val, err := ParseSquareS2I(enPassant)
	if err != nil {
		board.EnPassantSquare = uint8(255)
	} else {
		board.EnPassantSquare = val
	}
	halfMoveStr := fenParts[4]
	fullMoveStr := fenParts[5]
	halfMove64, err := strconv.ParseUint(halfMoveStr, 10, 8)
	if err != nil {
		return nil, err
	}
	board.HalfMoveClock = uint8(halfMove64)
	fullMove64, err := strconv.ParseUint(fullMoveStr, 10, 16)
	if err != nil {
		return nil, err
	}
	board.FullMoveNumber = uint16(fullMove64)
	return board, nil
}

func ParsePieces(board *Board, piecesRows []string) {
	for i, row := range piecesRows {
		rank := 7 - i
		file := 0
		for _, col := range row {
			square := rank*8 + file
			info, ok := pieceInfo[col]
			if !ok {
				file += int(col - '0')
			} else {
				color := info[0]
				piece := info[1]
				SetBit(&board.Colors[color], uint8(square))
				SetBit(&board.Pieces[piece], uint8(square))
				file++
			}

		}
	}
}
func ParseSquareS2I(s string) (uint8, error) {
	n := len(s)
	if n != 2 {
		return 0, errors.New("invalid Square")
	}
	file := s[0] - 'a'
	rank := s[1] - '1'
	if file > 7 || rank > 7 {
		return 0, errors.New("invalid Square")
	}
	return rank*8 + file, nil
}
func ParseSquareI2S(sq uint8) (string, error) {
	if sq > 63 {
		return "", errors.New("invalid Square")
	}
	rank, file := sq/8, sq%8
	return fmt.Sprintf("%c%c", file+'a', rank+'1'), nil
}
