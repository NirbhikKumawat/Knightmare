package chess

import (
	"fmt"
	"math/bits"
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
	Empty
)

var CastlingMasks = [64]uint8{
	13, 15, 15, 15, 12, 15, 15, 14,
	15, 15, 15, 15, 15, 15, 15, 15,
	15, 15, 15, 15, 15, 15, 15, 15,
	15, 15, 15, 15, 15, 15, 15, 15,
	15, 15, 15, 15, 15, 15, 15, 15,
	15, 15, 15, 15, 15, 15, 15, 15,
	15, 15, 15, 15, 15, 15, 15, 15,
	7, 15, 15, 15, 3, 15, 15, 11,
}

var (
	WhiteKingside  uint8 = 1
	WhiteQueenside uint8 = 2
	BlackKingside  uint8 = 4
	BlackQueenside uint8 = 8
)

const (
	NotAFile  uint64 = 0xfefefefefefefefe
	NotHFile  uint64 = 0x7f7f7f7f7f7f7f7f
	NotABFile uint64 = 0xfcfcfcfcfcfcfcfc
	NotGHFile uint64 = 0x3f3f3f3f3f3f3f3f
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
	Hash            uint64
}

func init() {
	for sq := 0; sq < 64; sq++ {
		KnightAttacks[sq] = maskKnightAttacks(uint8(sq))
		KingAttacks[sq] = maskKingAttacks(uint8(sq))
		PawnAttacks[White][sq] = maskPawnAttacks(White, uint8(sq))
		PawnAttacks[Black][sq] = maskPawnAttacks(Black, uint8(sq))
	}
	initSliders()
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
func (board *Board) GeneratePseudoLegalMoves() MoveList {
	ml := MoveList{}
	board.generatePawnMoves(&ml)
	board.generateKnightMoves(&ml)
	board.generateKingMoves(&ml)
	board.generateSliderMoves(&ml, Bishop)
	board.generateSliderMoves(&ml, Rook)
	board.generateSliderMoves(&ml, Queen)

	return ml
}
func (board *Board) IsSquareAttacked(sq uint8, attackerColor uint8) bool {
	var pawn uint64
	if attackerColor == White {
		pawn = PawnAttacks[Black][sq] & board.Colors[attackerColor] & board.Pieces[Pawn]
	}
	if attackerColor == Black {
		pawn = PawnAttacks[White][sq] & board.Colors[attackerColor] & board.Pieces[Pawn]
	}
	if pawn != 0 {
		return true
	}
	knight := KnightAttacks[sq] & board.Colors[attackerColor] & board.Pieces[Knight]
	if knight != 0 {
		return true
	}
	king := KingAttacks[sq] & board.Colors[attackerColor] & board.Pieces[King]
	if king != 0 {
		return true
	}
	occupied := board.Colors[White] | board.Colors[Black]
	bishop := GetBishopAttacks(sq, occupied) & (board.Pieces[Bishop] | board.Pieces[Queen]) & board.Colors[attackerColor]
	if bishop != 0 {
		return true
	}
	rook := GetRookAttacks(sq, occupied) & (board.Pieces[Rook] | board.Pieces[Queen]) & board.Colors[attackerColor]
	if rook != 0 {
		return true
	}
	return false
}

func (board *Board) MakeMove(m Move) bool {
	from := m.From()
	to := m.To()
	flags := m.Flags()
	color := board.SideToMove
	piece := board.GetPieceType(from)
	if board.EnPassantSquare != 255 {
		board.Hash ^= ZobristEnPassant[board.EnPassantSquare%8]
	}
	board.EnPassantSquare = 255
	board.HalfMoveClock++
	ClearBit(&board.Colors[color], from)
	ClearBit(&board.Pieces[piece], from)
	board.Hash ^= ZobristPieces[color][piece][from]
	board.Hash ^= ZobristPieces[color][piece][to]
	SetBit(&board.Colors[color], to)
	epiece := board.GetPieceType(to)
	board.Hash ^= ZobristCastling[board.CastlingRights]
	board.CastlingRights &= CastlingMasks[from] & CastlingMasks[to]
	board.Hash ^= ZobristCastling[board.CastlingRights]
	if color == Black {
		board.FullMoveNumber++
	}
	if flags == 4 {
		board.HalfMoveClock = 0
		ClearBit(&board.Colors[1^color], to)
		ClearBit(&board.Pieces[epiece], to)
		board.Hash ^= ZobristPieces[1^color][epiece][to]
	}
	if piece == Pawn {
		board.HalfMoveClock = 0
		switch flags {
		case 1:
			if color == White {
				board.EnPassantSquare = to - 8
			} else {
				board.EnPassantSquare = to + 8
			}
			board.Hash ^= ZobristEnPassant[board.EnPassantSquare%8]
			SetBit(&board.Pieces[Pawn], to)
		case 5:
			SetBit(&board.Pieces[Pawn], to)
			if color == Black {
				ClearBit(&board.Colors[1^color], to+8)
				ClearBit(&board.Pieces[Pawn], to+8)
				board.Hash ^= ZobristPieces[White][Pawn][to+8]
			} else {
				ClearBit(&board.Colors[1^color], to-8)
				ClearBit(&board.Pieces[Pawn], to-8)
				board.Hash ^= ZobristPieces[Black][Pawn][to-8]
			}
		case 8:
			SetBit(&board.Pieces[Knight], to)
			board.Hash ^= ZobristPieces[color][piece][to]
			board.Hash ^= ZobristPieces[color][Knight][to]
		case 9:
			SetBit(&board.Pieces[Bishop], to)
			board.Hash ^= ZobristPieces[color][piece][to]
			board.Hash ^= ZobristPieces[color][Bishop][to]
		case 10:
			SetBit(&board.Pieces[Rook], to)
			board.Hash ^= ZobristPieces[color][piece][to]
			board.Hash ^= ZobristPieces[color][Rook][to]
		case 11:
			SetBit(&board.Pieces[Queen], to)
			board.Hash ^= ZobristPieces[color][piece][to]
			board.Hash ^= ZobristPieces[color][Queen][to]
		case 12:
			SetBit(&board.Pieces[Knight], to)
			ClearBit(&board.Colors[1^color], to)
			ClearBit(&board.Pieces[epiece], to)
			board.Hash ^= ZobristPieces[color][piece][to]
			board.Hash ^= ZobristPieces[color][Knight][to]
			board.Hash ^= ZobristPieces[1^color][epiece][to]
		case 13:
			SetBit(&board.Pieces[Bishop], to)
			ClearBit(&board.Colors[1^color], to)
			ClearBit(&board.Pieces[epiece], to)
			board.Hash ^= ZobristPieces[color][piece][to]
			board.Hash ^= ZobristPieces[color][Bishop][to]
			board.Hash ^= ZobristPieces[1^color][epiece][to]
		case 14:
			SetBit(&board.Pieces[Rook], to)
			ClearBit(&board.Colors[1^color], to)
			ClearBit(&board.Pieces[epiece], to)
			board.Hash ^= ZobristPieces[color][piece][to]
			board.Hash ^= ZobristPieces[color][Rook][to]
			board.Hash ^= ZobristPieces[1^color][epiece][to]
		case 15:
			SetBit(&board.Pieces[Queen], to)
			ClearBit(&board.Colors[1^color], to)
			ClearBit(&board.Pieces[epiece], to)
			board.Hash ^= ZobristPieces[color][piece][to]
			board.Hash ^= ZobristPieces[color][Queen][to]
			board.Hash ^= ZobristPieces[1^color][epiece][to]
		default:
			SetBit(&board.Pieces[Pawn], to)
		}
	} else {
		SetBit(&board.Pieces[piece], to)
		if piece == King {
			switch flags {
			case 2:
				if color == White {
					ClearBit(&board.Pieces[Rook], 7)
					ClearBit(&board.Colors[White], 7)
					board.Hash ^= ZobristPieces[White][Rook][7]
					SetBit(&board.Pieces[Rook], 5)
					SetBit(&board.Colors[White], 5)
					board.Hash ^= ZobristPieces[White][Rook][5]
					if board.IsSquareAttacked(4, Black) || board.IsSquareAttacked(5, Black) {
						return false
					}
				} else {
					ClearBit(&board.Pieces[Rook], 63)
					ClearBit(&board.Colors[Black], 63)
					board.Hash ^= ZobristPieces[Black][Rook][63]
					SetBit(&board.Pieces[Rook], 61)
					SetBit(&board.Colors[Black], 61)
					board.Hash ^= ZobristPieces[Black][Rook][61]
					if board.IsSquareAttacked(60, White) || board.IsSquareAttacked(61, White) {
						return false
					}
				}
			case 3:
				if color == White {
					ClearBit(&board.Pieces[Rook], 0)
					ClearBit(&board.Colors[White], 0)
					board.Hash ^= ZobristPieces[White][Rook][0]
					SetBit(&board.Pieces[Rook], 3)
					SetBit(&board.Colors[White], 3)
					board.Hash ^= ZobristPieces[White][Rook][3]
					if board.IsSquareAttacked(4, Black) || board.IsSquareAttacked(3, Black) {
						return false
					}
				} else {
					ClearBit(&board.Pieces[Rook], 56)
					ClearBit(&board.Colors[Black], 56)
					board.Hash ^= ZobristPieces[Black][Rook][56]
					SetBit(&board.Pieces[Rook], 59)
					SetBit(&board.Colors[Black], 59)
					board.Hash ^= ZobristPieces[Black][Rook][59]
					if board.IsSquareAttacked(60, White) || board.IsSquareAttacked(59, White) {
						return false
					}
				}
			}
		}
	}
	kingSq := uint8(bits.TrailingZeros64(board.Colors[color] & board.Pieces[King]))
	if board.IsSquareAttacked(kingSq, color^1) {
		return false
	}
	board.SideToMove ^= 1
	board.Hash ^= ZobristSideToMove
	return true

}
func (board *Board) GetPieceType(sq uint8) uint8 {
	for i := Pawn; i <= King; i++ {
		if GetBit(board.Pieces[i], sq) != 0 {
			return uint8(i)
		}
	}
	return Empty
}
func (board *Board) GetColorType(sq uint8) uint8 {
	if GetBit(board.Colors[White], sq) != 0 {
		return uint8(Black)
	}
	return White
}

// 0000	0	Quiet move (Default)
// 0001	1	Double pawn push
// 0010	2	King-side castle
// 0011	3	Queen-side castle
// 0100	4	Standard capture
// 0101	5	En Passant capture
// 1000	8	Knight promotion
// 1001	9	Bishop promotion
// 1010	10	Rook promotion
// 1011	11	Queen promotion
// 1100	12	Knight promotion + capture
// 1101	13	Bishop promotion + capture
// 1110	14	Rook promotion + capture
// 1111	15	Queen promotion + capture
