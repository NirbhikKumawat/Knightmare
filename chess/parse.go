package chess

import (
	"errors"
	"fmt"
	"math/bits"
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
	cr := uint8(0)
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
func (board *Board) GenerateFEN() string {
	occupied := board.Colors[White] | board.Colors[Black]
	var b strings.Builder
	for rank := 7; rank >= 0; rank-- {
		counter := 0
		for file := 0; file <= 7; file++ {
			bit := GetBit(occupied, uint8(rank*8+file))
			if bit == 0 {
				counter++
			} else {
				if counter != 0 {
					b.WriteString(fmt.Sprintf("%d", counter))
				}
				counter = 0
				color := Black
				if (bit & board.Colors[White]) != 0 {
					color = White
				}
				piece := board.GetPieceType(uint8(rank*8 + file))
				switch color {
				case White:
					switch piece {
					case Pawn:
						b.WriteString("P")
					case Knight:
						b.WriteString("N")
					case Bishop:
						b.WriteString("B")
					case Rook:
						b.WriteString("R")
					case Queen:
						b.WriteString("Q")
					case King:
						b.WriteString("K")
					default:
					}
				case Black:
					switch piece {
					case Pawn:
						b.WriteString("p")
					case Knight:
						b.WriteString("n")
					case Bishop:
						b.WriteString("b")
					case Rook:
						b.WriteString("r")
					case Queen:
						b.WriteString("q")
					case King:
						b.WriteString("k")
					default:
					}
				}
			}
		}
		if counter != 0 {
			b.WriteString(fmt.Sprintf("%d", counter))
		}
		if rank != 0 {
			b.WriteString("/")
		}
	}
	toMove := "w"
	if board.SideToMove == Black {
		toMove = "b"
	}
	b.WriteString(fmt.Sprintf(" %s", toMove))
	castling := board.CastlingRights
	if castling == 0 {
		b.WriteString(" -")
	} else {
		b.WriteString(" ")
		if castling&WhiteKingside != 0 {
			b.WriteString("K")
		}
		if castling&WhiteQueenside != 0 {
			b.WriteString("Q")
		}
		if castling&BlackKingside != 0 {
			b.WriteString("k")
		}
		if castling&BlackQueenside != 0 {
			b.WriteString("q")
		}
	}

	if board.EnPassantSquare == 255 {
		b.WriteString(" -")
	} else {
		square, err := ParseSquareI2S(board.EnPassantSquare)
		if err != nil {
			panic(err)
		}
		b.WriteString(fmt.Sprintf(" %s", square))
	}
	b.WriteString(fmt.Sprintf(" %d %d", board.HalfMoveClock, board.FullMoveNumber))
	return b.String()
}
func (board *Board) MoveToSAN(m Move) string {
	from := m.From()
	to := m.To()
	flags := m.Flags()
	piece := board.GetPieceType(from)
	if piece == King {
		if flags == 2 {
			return "O-O" + board.getCheckSuffix(m)
		} else if flags == 3 {
			return "O-O-O" + board.getCheckSuffix(m)
		}
	}
	toStr, _ := ParseSquareI2S(to)
	fromStr, _ := ParseSquareI2S(from)
	var san strings.Builder
	pieceChars := []string{"", "N", "B", "R", "Q", "K"}
	san.WriteString(pieceChars[piece])
	if piece != Pawn && piece != King {
		moves := board.GeneratePseudoLegalMoves()
		var conflicts []Move
		for i := 0; i < moves.Count; i++ {
			testMove := moves.Moves[i]
			boardCopy := *board

			if !boardCopy.MakeMove(testMove) {
				continue
			}

			if testMove.To() == to && testMove.From() != from && board.GetPieceType(testMove.From()) == piece {
				conflicts = append(conflicts, testMove)
			}
		}

		if len(conflicts) > 0 {
			sameFile := false
			sameRank := false
			for _, conflict := range conflicts {
				if conflict.From()%8 == from%8 {
					sameFile = true
				}
				if conflict.From()/8 == from/8 {
					sameRank = true
				}
			}
			if !sameFile {
				san.WriteByte(fromStr[0])
			} else if !sameRank {
				san.WriteByte(fromStr[1])
			} else {
				san.WriteString(fromStr)
			}
		}
	}

	isCapture := flags == 4 || flags == 5 || (flags >= 12 && flags <= 15)
	if piece == Pawn && isCapture {
		san.WriteByte(fromStr[0])
		san.WriteString("x")
	} else if isCapture {
		san.WriteString("x")
	}

	san.WriteString(toStr)

	if flags >= 8 && flags <= 15 {
		san.WriteString("=")
		promoFlags := flags & 3
		promoChars := []string{"N", "B", "R", "Q"}
		san.WriteString(promoChars[promoFlags])
	}

	san.WriteString(board.getCheckSuffix(m))

	return san.String()
}

func (board *Board) getCheckSuffix(m Move) string {
	boardCopy := *board
	boardCopy.MakeMove(m)

	enemyColor := boardCopy.SideToMove
	kingBits := boardCopy.Colors[enemyColor] & boardCopy.Pieces[King]

	if kingBits == 0 {
		return ""
	}
	kingSq := uint8(bits.TrailingZeros64(kingBits))
	isCheck := boardCopy.IsSquareAttacked(kingSq, enemyColor^1)
	if !isCheck {
		return ""
	}

	moves := boardCopy.GeneratePseudoLegalMoves()
	hasLegalMove := false
	for i := 0; i < moves.Count; i++ {
		testBoard := boardCopy
		if testBoard.MakeMove(moves.Moves[i]) {
			hasLegalMove = true
			break
		}
	}
	if !hasLegalMove {
		return "#"
	}
	return "+"
}
func (board *Board) ParseSAN(san string) (Move, error) {
	san = strings.TrimSpace(san)
	moves := board.GeneratePseudoLegalMoves()
	for i := 0; i < moves.Count; i++ {
		testMove := moves.Moves[i]

		boardCopy := *board
		if !boardCopy.MakeMove(testMove) {
			continue
		}
		generatedSAN := board.MoveToSAN(testMove)
		if generatedSAN == san {
			return testMove, nil
		}
		cleanInput := strings.TrimRight(san, "+#")
		cleanGenerated := strings.TrimRight(generatedSAN, "+#")
		if cleanInput != cleanGenerated {
			return testMove, nil
		}
	}
	return 0, errors.New("illegal or unrecognized SAN" + san)
}
