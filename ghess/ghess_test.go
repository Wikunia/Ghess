package ghess

import (
	"testing"
)

func TestNumMoves(t *testing.T) {
	var moves []JSONMove
	startFEN := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
	for _, test := range numMovesTests {
		board := getBoardFromFen(startFEN)
		for _, moveStr := range test.moves {
			err := board.moveLongAlgebraic(moveStr)
			if err != nil {
				t.Errorf(err.Error())
			}
		}
		fen := board.getFen()
		n := getNumberOfMoves(fen, test.ply, &moves)
		if n != test.expected {
			t.Errorf("Moves(%v) with ply: %d expected %d, Actual %d", test.moves, test.ply, test.expected, n)
		}
	}
}

func TestFen(t *testing.T) {
	// not an actual FEN
	expected := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq c2 0 1"
	board := getBoardFromFen("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq c2 0 1")
	actual := board.getFen()
	if actual != expected {
		t.Errorf("FEN expected: %s, actual: %s", expected, actual)
	}
}
