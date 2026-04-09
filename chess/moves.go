package chess

// Move stores a chess move, flags(4-bits)|to-square(6-bits)|from-square(6-bits)
type Move uint16

// MoveList stores moves at a given board state
type MoveList struct {
	Moves [256]Move // moves in the MoveList
	Count int       // no of valid moves in the MoveList
}

// Add adds a new Move to the MoveList
func (ml *MoveList) Add(m Move) {
	if ml.Count >= 256 {
		panic("Too many moves")
	}
	ml.Moves[ml.Count] = m
	ml.Count++
}

// NewMove packs starting square,ending square and flags into a 16-bit unsigned integer
func NewMove(from uint8, to uint8, flags uint16) Move {
	return Move(flags<<12 | uint16(to)<<6 | uint16(from))
}

// From extracts the starting square index from a Move
func (m Move) From() uint8 {
	return uint8(m & 63)
}

// To extracts the final square index from a Move
func (m Move) To() uint8 {
	return uint8((m >> 6) & 63)
}

// Flags extracts the flag from a Move
func (m Move) Flags() uint16 {
	return uint16(m >> 12)
}
