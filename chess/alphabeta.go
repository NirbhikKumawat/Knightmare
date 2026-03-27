package chess

import "math/bits"

func (board *Board) AlphaBeta(depth, alpha, beta int, isMax bool) int {
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
			score := boardCopy.AlphaBeta(depth-1, alpha, beta, false)
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
			score := boardCopy.AlphaBeta(depth-1, alpha, beta, true)
			if score < bestScore {
				bestScore = score
				bestMove = moves.Moves[i]
			}
			if score < beta {
				flag = TTExact
				beta = score
			}
			if beta <= alpha {
				flag = TTAlpha
				break
			}
		}
	}
	board.StoreTT(depth, bestScore, flag, bestMove)
	return bestScore
}
func (board *Board) QuiescenceSearch(alpha, beta int, isMax bool) int {
	SearchNodes++
	CheckTime()
	if StopSearch {
		return 0
	}
	standPat := board.Evaluate()
	color := board.SideToMove
	kingsq := uint8(bits.TrailingZeros64(board.Colors[color] & board.Pieces[King]))
	inCheck := board.IsSquareAttacked(uint8(kingsq), color^1)
	if !inCheck {
		if isMax {
			if standPat >= beta {
				return beta
			}
			if standPat > alpha {
				alpha = standPat
			}
		} else {
			if standPat <= alpha {
				return alpha
			}
			if standPat < beta {
				beta = standPat
			}
		}
	}
	bestScore := standPat
	if inCheck {
		if isMax {
			bestScore = ColorScores[White]
		} else {
			bestScore = ColorScores[Black]
		}
	}
	moves := board.GenerateLegalMoves()
	board.SortMoves(&moves, 0)
	for i := 0; i < moves.Count; i++ {
		flags := moves.Moves[i].Flags()
		isCapture := (flags == 4 || flags == 5) || (flags >= 11)
		if !isCapture {
			continue
		} else {
			boardCopy := *board
			boardCopy.MakeMove(moves.Moves[i])
			if isMax {
				score := boardCopy.QuiescenceSearch(alpha, beta, false)
				if score > bestScore {
					bestScore = score
				}
				if score > alpha {
					alpha = score
				}
				if alpha >= beta {
					break
				}
			} else {
				score := boardCopy.QuiescenceSearch(alpha, beta, true)
				if score < bestScore {
					bestScore = score
				}
				if score < beta {
					beta = score
				}
				if alpha >= beta {
					break
				}
			}
		}
	}
	return bestScore
}
