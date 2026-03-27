package chess

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func UCILoop() {
	scanner := bufio.NewScanner(os.Stdin)
	board, _ := ParseFEN("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		tokens := strings.Fields(line)
		command := tokens[0]
		switch command {
		case "uci":
			fmt.Println("id name Knightmare")
			fmt.Println("id author NirbhikTheNice")
			fmt.Println("uciok")
		case "isready":
			fmt.Println("readyok")
		case "ucinewgame":
			TranspositionTable = [TTSize]TTEntry{}
		case "position":
			ParsePosition(board, tokens)
		case "go":
			ParseGo(board, tokens)
		case "stop":
			StopSearch = true
		case "quit":
			return
		}
	}
}
func (m Move) ToLAN() string {
	fromStr, _ := ParseSquareI2S(m.From())
	toStr, _ := ParseSquareI2S(m.To())
	lan := fromStr + toStr
	flag := m.Flags()
	switch flag {
	case 8, 12:
		lan += "n"
	case 9, 13:
		lan += "b"
	case 10, 14:
		lan += "r"
	case 11, 15:
		lan += "q"
	}
	return lan
}
func (board *Board) ParseMoveLAN(lan string) Move {
	moves := board.GenerateLegalMoves()
	var move string
	for i := 0; i < moves.Count; i++ {
		move = moves.Moves[i].ToLAN()
		if move == lan {
			return moves.Moves[i]
		}
	}
	return Move(0)
}
func ParsePosition(board *Board, tokens []string) {
	movesIndex := len(tokens)
	for i, token := range tokens {
		if token == "moves" {
			movesIndex = i
			break
		}
	}
	switch tokens[1] {
	case "startpos":
		newBoard, _ := ParseFEN("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
		*board = *newBoard
	case "fen":
		fen := strings.Join(tokens[2:movesIndex], " ")
		newBoard, _ := ParseFEN(fen)
		*board = *newBoard
	}
	for i := movesIndex + 1; i < len(tokens); i++ {
		move := board.ParseMoveLAN(tokens[i])
		if ok := board.MakeMove(move); !ok {
			return
		}
	}
}
func ParseGo(board *Board, tokens []string) {
	side := board.SideToMove
	var btime, winc, binc, wtime, moveTime int
	for i := 1; i < len(tokens); i++ {

		if i+1 < len(tokens) {
			switch tokens[i] {
			case "btime":
				btime, _ = strconv.Atoi(tokens[i+1])
			case "wtime":
				wtime, _ = strconv.Atoi(tokens[i+1])
			case "binc":
				binc, _ = strconv.Atoi(tokens[i+1])
			case "winc":
				winc, _ = strconv.Atoi(tokens[i+1])
			case "movetime":
				moveTime, _ = strconv.Atoi(tokens[i+1])
			}
		}
	}
	if moveTime == 0 {
		if side == White && wtime > 0 {
			moveTime = wtime/30 + winc/2
		} else if side == Black && btime > 0 {
			moveTime = btime/30 + binc/2
		} else {
			moveTime = 2000
		}
	}
	StopSearch = false
	searchBoard := *board
	go func(b *Board) {
		bestMove := b.SearchWithTime(int64(moveTime))
		fmt.Printf("bestmove %s\n", bestMove.ToLAN())
	}(&searchBoard)
}
