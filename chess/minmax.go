package chess

import "math/bits"

var PieceScores = [6]int{
	100, 300, 300, 500, 900, 10000,
}
var ColorScores = [2]int{
	-4000, 4000,
}

func (board *Board) Evaluate() int {
	score := 0
	for piece := Pawn; piece < King; piece++ {
		score += PieceScores[piece] * bits.OnesCount64(board.Colors[White]&board.Pieces[piece])
	}
	for piece := Pawn; piece < King; piece++ {
		score -= PieceScores[piece] * bits.OnesCount64(board.Colors[Black]&board.Pieces[piece])
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
	board.SortMoves(&moves)
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
