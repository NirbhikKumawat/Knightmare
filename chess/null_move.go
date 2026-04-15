package chess

import "math/bits"

func (board *Board) HasBigPieces(color uint8) bool {
	return board.Colors[color]&(board.Pieces[Rook]|board.Pieces[Queen]|board.Pieces[Bishop]|board.Pieces[Knight]) != 0
}
func (board *Board) MakeNullMove() {
	board.SideToMove = board.SideToMove ^ 1
	board.Hash ^= ZobristSideToMove
	if board.EnPassantSquare != 255 {
		board.Hash ^= ZobristEnPassant[board.EnPassantSquare%8]
		board.EnPassantSquare = 255
	}
}

func (board *Board) AlphaBetaNull(depth, alpha, beta int, isMax, canNull bool) int {
	SearchNodes++
	CheckTime()
	if StopSearch {
		return 0
	}
	cscore, cflag, cdepth, cbestmove, ok := board.ProbeTT()
	if ok {
		if cdepth >= depth {
			switch cflag {
			case TTExact:
				return cscore
			case TTAlpha:
				if cscore <= alpha {
					return cscore
				}
			case TTBeta:
				if cscore >= beta {
					return cscore
				}
			}
		}
	}
	if depth == 0 {
		return board.QuiescenceSearch(alpha, beta, isMax)
	}
	moves := board.GenerateLegalMoves()
	color := board.SideToMove

	if moves.Count == 0 {
		king := bits.TrailingZeros64(board.Colors[color] & board.Pieces[King])
		if board.IsSquareAttacked(uint8(king), color^1) {
			if isMax {
				return ColorScores[White] + depth
			} else {
				return ColorScores[Black] - depth
			}
		} else {
			return 0
		}
	}
	kingSq := bits.TrailingZeros64(board.Colors[color] & board.Pieces[King])
	inCheck := board.IsSquareAttacked(uint8(kingSq), color^1)

	R := 2
	if canNull && depth >= 3 && !inCheck && board.HasBigPieces(color) {
		boardCopy := *board
		boardCopy.MakeNullMove()
		if isMax {
			score := boardCopy.AlphaBetaNull(depth-1-R, beta-1, beta, false, false)
			if score >= beta {
				return beta
			}
		} else {
			score := boardCopy.AlphaBetaNull(depth-1-R, alpha, alpha+1, true, false)
			if score <= alpha {
				return alpha
			}
		}
	}
	board.SortMoves(&moves, cbestmove)
	var bestScore int
	var bestMove Move
	var flag int
	if isMax {
		flag = TTAlpha
		bestScore = ColorScores[White]
		for i := 0; i < moves.Count; i++ {
			boardCopy := *board
			boardCopy.MakeMove(moves.Moves[i])
			score := boardCopy.AlphaBetaNull(depth-1, alpha, beta, false, true)
			if score > bestScore {
				bestScore = score
				bestMove = moves.Moves[i]
			}
			if score > alpha {
				flag = TTExact
				alpha = score
			}
			if alpha >= beta {
				flag = TTBeta
				break
			}
		}
	} else {
		flag = TTBeta
		bestScore = ColorScores[Black]
		for i := 0; i < moves.Count; i++ {
			boardCopy := *board
			boardCopy.MakeMove(moves.Moves[i])
			score := boardCopy.AlphaBetaNull(depth-1, alpha, beta, true, true)
			if score < bestScore {
				bestScore = score
				bestMove = moves.Moves[i]
			}
			if score < beta {
				flag = TTExact
				beta = score
			}
			if alpha >= beta {
				flag = TTAlpha
				break
			}
		}
	}
	board.StoreTT(depth, bestScore, flag, bestMove)
	return bestScore
}
