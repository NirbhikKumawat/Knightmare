package chess

func (board *Board) ScoreMove(m Move, hashMove Move) int {
	if m == hashMove {
		return 20000000
	}
	score := 0
	from := m.From()
	to := m.To()
	flags := m.Flags()
	attacker := board.GetPieceType(from)

	isCapture := flags == 4 || flags == 5 || (flags >= 12 && flags <= 15)
	if isCapture {
		victim := board.GetPieceType(to)
		if flags == 5 {
			victim = Pawn
		}
		score += 1000000 + (PieceScores[victim] * 10) - PieceScores[attacker]
	}
	if flags >= 8 && flags <= 15 {
		if flags == 11 || flags == 15 {
			score += 900000
		} else {
			score += 300000
		}
	}

	return score
}
func (board *Board) SortMoves(moves *MoveList, hashMove Move) {
	scores := make([]int, moves.Count)
	for i := 0; i < moves.Count; i++ {
		scores[i] = board.ScoreMove(moves.Moves[i], hashMove)
	}
	for i := 1; i < moves.Count; i++ {
		keyMove := moves.Moves[i]
		keyScore := scores[i]
		j := i - 1
		for j >= 0 && scores[j] < keyScore {
			scores[j+1] = scores[j]
			moves.Moves[j+1] = moves.Moves[j]
			j--
		}
		scores[j+1] = keyScore
		moves.Moves[j+1] = keyMove
	}
}
