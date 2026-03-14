package main

import (
	"fmt"
	"gochess/chess"
)

func main() {
	fen := "8/8/8/4k3/8/8/4K2p/8 b - - 0 1"
	board, err := chess.ParseFEN(fen)
	if err != nil {
		fmt.Println(err)
	}
	board.Print()
}
