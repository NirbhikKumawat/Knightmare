package chess

import (
	"testing"
)

func TestBitBoardOperations(t *testing.T) {
	var bb uint64 = 0
	SetBit(&bb, 0)
	SetBit(&bb, 63)
	if GetBit(bb, 0) == 0 {
		t.Errorf("Expected bit 0 (A1) to be set")
	}
	if GetBit(bb, 63) == 0 {
		t.Errorf("Expected bit 63 (H8) to be set")
	}
	if GetBit(bb, 1) != 0 {
		t.Errorf("Expected bit 1 (B1) to be clear")
	}
	sq := PopBit(&bb)
	if sq != 0 {
		t.Errorf("Expected PopBit to return 0, got %d", sq)
	}
	if GetBit(bb, 0) != 0 {
		t.Errorf("Expected bit 0 to be cleared by PopBit")
	}
	ClearBit(&bb, 63)
	if bb != 0 {
		t.Errorf("Expected bitboard to be empty (0) after clearing last bit, got %d", bb)
	}
}
func TestSquareParsing(t *testing.T) {
	tests := []struct {
		str       string
		idx       uint8
		expectErr bool
	}{
		{"a1", 0, false},
		{"h1", 7, false},
		{"a8", 56, false},
		{"h8", 63, false},
		{"e4", 28, false},
		{"i9", 0, true},
		{"a0", 0, true},
		{"", 0, true},
	}

	for _, tt := range tests {
		idx, err := ParseSquareS2I(tt.str)
		if (err != nil) != tt.expectErr {
			t.Errorf("ParseSquareS2I(%s) unexpected error status: %v", tt.str, err)
		}
		if !tt.expectErr && idx != tt.idx {
			t.Errorf("ParseSquareS2I(%s) expected %d, got %d", tt.str, tt.idx, idx)
		}
		if !tt.expectErr {
			str, err := ParseSquareI2S(tt.idx)
			if err != nil {
				t.Errorf("ParseSquareI2S(%d) unexpected error: %v", tt.idx, err)
			}
			if str != tt.str {
				t.Errorf("ParseSquareI2S(%d) expected %s, got %s", tt.idx, tt.str, str)
			}
		}
	}
}
func TestParseFENStartingPosition(t *testing.T) {
	fen := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
	board, err := ParseFEN(fen)
	if err != nil {
		t.Fatalf("Failed to parse valid FEN: %v", err)
	}
	if board.SideToMove != White {
		t.Errorf("Expected White to move, got %d", board.SideToMove)
	}
	if GetBit(board.Colors[White], 12) == 0 || GetBit(board.Pieces[Pawn], 12) == 0 {
		t.Errorf("Expected White Pawn on e2 (Square 12)")
	}
	if GetBit(board.Colors[Black], 60) == 0 || GetBit(board.Pieces[King], 60) == 0 {
		t.Errorf("Expected Black King on e8 (Square 60)")
	}
	if GetBit(board.Colors[White], 28) != 0 || GetBit(board.Colors[Black], 28) != 0 {
		t.Errorf("Expected e4 to be empty")
	}
}
