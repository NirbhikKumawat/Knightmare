package chess

import (
	"math/bits"
	"time"
)

// PieceScores stores points for each piece
var PieceScores = [6]int{
	100, 300, 300, 500, 900, 10000,
}

// ColorScores stores initial points of each color,White wants to maximize score while black wants to minimize score
var ColorScores = [2]int{
	-4000, 4000,
}

// SearchNodes indicates the number of nodes searched till now
var SearchNodes uint64

// EndTime tells till when to search
var EndTime int64

// StopSearch is a flag indicating whether to continue searching
var StopSearch bool

// CheckTime checks whether it exceeded time after every 2048 node searches
func CheckTime() {
	if SearchNodes%2048 == 0 {
		if time.Now().UnixMilli() >= EndTime {
			StopSearch = true
		}
	}
}

// Evaluate evaluates score at current game state
func (board *Board) Evaluate() int {
	score := 0
	var bb uint64
	for piece := Pawn; piece < King; piece++ {
		bb = board.Colors[White] & board.Pieces[piece]
		for bb != 0 {
			sq := PopBit(&bb)
			score += PSTs[piece][sq] + PieceScores[piece]
		}
	}
	for piece := Pawn; piece < King; piece++ {
		bb = board.Colors[Black] & board.Pieces[piece]
		for bb != 0 {
			sq := PopBit(&bb)
			score -= PSTs[piece][sq^56] + PieceScores[piece]
		}
	}
	return score
}

// Minimax returns the best possible score after a certain depth
func (board *Board) Minimax(depth int, isMax bool) int {
	if depth == 0 {
		return board.Evaluate()
	}
	moves := board.GenerateLegalMoves()
	color := board.SideToMove

	// Checkmate/Stalemate
	if moves.Count == 0 {
		king := bits.TrailingZeros64(board.Colors[color] & board.Pieces[King])
		if board.IsSquareAttacked(uint8(king), color^1) {
			// Losing one wants to delay checkmate
			if isMax {
				return ColorScores[White] + depth
			}
			return ColorScores[Black] - depth
		} else {
			return 0
		}
	}

	// Other side will try to win
	if isMax {
		bestScore := ColorScores[White]
		for i := 0; i < moves.Count; i++ {
			boardCopy := *board
			boardCopy.MakeMove(moves.Moves[i])
			score := boardCopy.Minimax(depth-1, false)
			if score > bestScore {
				bestScore = score
			}
		}
		return bestScore
	} else {
		bestScore := ColorScores[Black]
		for i := 0; i < moves.Count; i++ {
			boardCopy := *board
			boardCopy.MakeMove(moves.Moves[i])
			score := boardCopy.Minimax(depth-1, true)
			if score < bestScore {
				bestScore = score
			}
		}
		return bestScore
	}
}

// SearchBestMove returns the best possible move found
func (board *Board) SearchBestMove(depth int) Move {
	moves := board.GenerateLegalMoves()
	_, _, _, cbestmove, ok := board.ProbeTT()
	if !ok {
		cbestmove = Move(0)
	}
	board.SortMoves(&moves, cbestmove)
	if moves.Count == 0 {
		return Move(0)
	}
	bestMove := moves.Moves[0]
	isMax := board.SideToMove == White
	alpha := ColorScores[White]
	beta := ColorScores[Black]
	if isMax {
		bestScore := ColorScores[White]
		for i := 0; i < moves.Count; i++ {
			boardCopy := *board
			boardCopy.MakeMove(moves.Moves[i])
			score := boardCopy.AlphaBetaNull(depth-1, alpha, beta, false, true)
			if StopSearch {
				return Move(0)
			}
			if score > bestScore {
				bestScore = score
				bestMove = moves.Moves[i]
			}
			if score > alpha {
				alpha = score
			}
		}
	} else {
		bestScore := ColorScores[Black]
		for i := 0; i < moves.Count; i++ {
			boardCopy := *board
			boardCopy.MakeMove(moves.Moves[i])
			score := boardCopy.AlphaBetaNull(depth-1, alpha, beta, true, true)
			if StopSearch {
				return Move(0)
			}
			if score < bestScore {
				bestScore = score
				bestMove = moves.Moves[i]
			}
			if score < beta {
				beta = score
			}
		}
	}
	return bestMove
}

// SearchWithTime searches for a move with time,instead of depth
func (board *Board) SearchWithTime(timeLimitMs int64) Move {
	SearchNodes = 0
	StopSearch = false
	EndTime = time.Now().UnixMilli() + timeLimitMs
	move := board.GenerateLegalMoves()
	if move.Count == 0 {
		return Move(0)
	}
	bestMove := move.Moves[0]
	for depth := 1; depth < 100; depth++ {
		currBestMove := board.SearchBestMove(depth)
		if StopSearch {
			break
		}
		if currBestMove != Move(0) {
			bestMove = currBestMove
		}
	}
	return bestMove
}
