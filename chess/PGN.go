package chess

import (
	"errors"
	"strings"
)

// Game stores a game in PGN format
type Game struct {
	Headers map[string]string // Headers store metadata of the game
	Moves   []Move            // Moves stores the moves in the game
}

// ParsePGN parses a PGN file
func ParsePGN(pgnText string) (*Game, error) {
	game := &Game{
		Headers: make(map[string]string),
		Moves:   []Move{},
	}
	lines := strings.Split(pgnText, "\n")
	var moveTextBuilder strings.Builder
	for _, line := range lines {
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			content := line[1 : len(line)-1]
			parts := strings.SplitN(content, " ", 2)
			if len(parts) == 2 {
				key := parts[0]
				value := strings.Trim(parts[1], "\"")
				game.Headers[key] = value
			}
		} else {
			moveTextBuilder.WriteString(line + " ")
		}
	}
	moveText := moveTextBuilder.String()
	tokens := strings.Fields(moveText)

	board, err := ParseFEN("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
	if err != nil {
		return nil, err
	}
	for _, token := range tokens {
		if strings.HasSuffix(token, ".") {
			continue
		}
		if token == "1-0" || token == "0-1" || token == "1/2-1/2" || token == "*" {
			continue
		}
		move, err := board.ParseSAN(token)
		if err != nil {
			return nil, errors.New("failed to parse move " + token + ": " + err.Error())
		}
		game.Moves = append(game.Moves, move)
		if !board.MakeMove(move) {
			return nil, errors.New("parsed an illegal move: " + token)
		}
	}
	return game, nil
}
