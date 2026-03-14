package chess

import (
	"fmt"
)

const (
	White = iota
	Black
)
const (
	Pawn = iota
	Knight
	Bishop
	Rook
	Queen
	King
)
const (
	WhiteKingside  = 1
	WhiteQueenside = 2
	BlackKingside  = 4
	BlackQueenside = 8
)

var pieceChars = [2][6]rune{
	{'P', 'N', 'B', 'R', 'Q', 'K'},
	{'p', 'n', 'b', 'r', 'q', 'k'},
}
var pieceInfo = map[rune][2]int{
	'P': {0, 0},
	'N': {0, 1},
	'B': {0, 2},
	'R': {0, 3},
	'Q': {0, 4},
	'K': {0, 5},
	'p': {1, 0},
	'n': {1, 1},
	'b': {1, 2},
	'r': {1, 3},
	'q': {1, 4},
	'k': {1, 5},
}

type Board struct {
	Colors          [2]uint64
	Pieces          [6]uint64
	SideToMove      uint8
	CastlingRights  uint8
	EnPassantSquare uint8
	HalfMoveClock   uint8
	FullMoveNumber  uint16
}

func (board *Board) Print() {
	pieceChars := [2][6]rune{
		{'P', 'N', 'B', 'R', 'Q', 'K'},
		{'p', 'n', 'b', 'r', 'q', 'k'},
	}
	fmt.Println()
	for rank := 7; rank >= 0; rank-- {
		fmt.Printf("%d ", rank+1)
		for file := 0; file < 8; file++ {
			square := uint8(rank*8 + file)
			char := '.'
			for color := White; color <= Black; color++ {
				if GetBit(board.Colors[color], square) != 0 {
					for piece := Pawn; piece <= King; piece++ {
						if GetBit(board.Pieces[piece], square) != 0 {
							char = pieceChars[color][piece]
							break
						}
					}
					break
				}
			}
			fmt.Printf("%c ", char)
		}
		fmt.Println()
	}
	fmt.Println("\n  a b c d e f g h")
}
