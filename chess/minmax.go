package chess

import (
	"math/bits"
	"time"
)

var PieceScores = [6]int{
	100, 300, 300, 500, 900, 10000,
}
var ColorScores = [2]int{
	-4000, 4000,
}
var SearchNodes uint64
var EndTime int64
var StopSearch bool

func CheckTime() {
	if SearchNodes%2048 == 0 {
		if time.Now().Unix() > EndTime {
			StopSearch = true
		}
	}
}

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
func (board *Board) Minimax(depth int, isMax bool) int {
	if depth == 0 {
		return board.Evaluate()
	}
	moves := board.GenerateLegalMoves()
	color := board.SideToMove
	if moves.Count == 0 {
		king := bits.TrailingZeros64(board.Colors[color] & board.Pieces[King])
		if board.IsSquareAttacked(uint8(king), color^1) {
			if isMax {
				return ColorScores[White] + depth
			}
			return ColorScores[Black] - depth
		} else {
			return 0
		}
	}
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
			score := boardCopy.AlphaBeta(depth-1, alpha, beta, false)
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
			score := boardCopy.AlphaBeta(depth-1, alpha, beta, true)
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
func (board *Board) SearchWithTime(timeLimitMs int64) Move {
	SearchNodes = 0
	StopSearch = false
	EndTime = time.Now().Unix() + timeLimitMs
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
