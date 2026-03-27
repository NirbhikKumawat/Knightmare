package main

import (
	"github.com/NirbhikKumawat/GoChess/chess"
)

func main() {
	chess.InitZobrist()
	chess.UCILoop()
}
