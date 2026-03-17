package chess

import (
	"testing"
)

func TestPerftStartingPosition(t *testing.T) {
	fen := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
	board, err := ParseFEN(fen)
	if err != nil {
		t.Fatalf("Failed to parse FEN: %v", err)
	}

	tests := []struct {
		depth    int
		expected uint64
	}{
		{1, 20},
		{2, 400},
		{3, 8902},
		{4, 197281},
		{5, 4865609},
		{6, 119060324},
		//{7, 3195901860},
		//{8, 84998978956},
		//{9, 2439530234167},
		//{10, 69352859712417},
		//{11, 2097651003696806},
		//{12, 62854969236701747},
		//{13, 1981066775000396239},
		//{14, 61885021521585529237},
		//{15, 2015099950053364471960},
	}

	for _, tt := range tests {
		nodes := board.Perft(tt.depth)
		if nodes != tt.expected {
			t.Errorf("Perft(%d) failed: expected %d, got %d", tt.depth, tt.expected, nodes)
		} else {
			t.Logf("Perft(%d) passed: %d nodes", tt.depth, nodes)
		}
	}
}
