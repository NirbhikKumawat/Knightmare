package main

import "gochess/chess"

func main() {
	/*fen := "rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1"
	board, err := chess.ParseFEN(fen)
	if err != nil {
		fmt.Println(err)
	}
	board.Print()
	board.PerftDivide(2)*/
	chess.GenerateAllMagics()
}
