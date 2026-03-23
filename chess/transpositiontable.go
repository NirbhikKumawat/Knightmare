package chess

const (
	TTExact = iota
	TTAlpha
	TTBeta
)

type TTEntry struct {
	Hash     uint64
	Depth    int
	Score    int
	Flag     int
	BestMove Move
}

const TTSize = 1000000

var TranspositionTable [TTSize]TTEntry

func (board *Board) ProbeTT() (int, int, int, Move, bool) {
	currEntry := TranspositionTable[board.Hash%TTSize]
	if currEntry.Hash == board.Hash {
		return currEntry.Score, currEntry.Flag, currEntry.Depth, currEntry.BestMove, true
	}
	return 0, 0, 0, 0, false
}
func (board *Board) StoreTT(depth, score, flag int, bestMove Move) {
	entry := TTEntry{
		Hash:     board.Hash,
		Depth:    depth,
		Score:    score,
		Flag:     flag,
		BestMove: bestMove,
	}
	TranspositionTable[board.Hash%TTSize] = entry
}
