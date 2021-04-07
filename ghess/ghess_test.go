package ghess

import (
	"testing"
)

/*
func TestNumMoves(t *testing.T) {
	startFEN := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
	for _, test := range numMovesTests {
		board := GetBoardFromFen(startFEN)
		for _, moveStr := range test.moves {
			err := board.MoveLongAlgebraic(moveStr)
			if err != nil {
				t.Errorf(err.Error())
			}
		}
		n := board.GetNumberOfMoves(test.ply, false)
		if n != test.expected {
			t.Errorf("Moves(%v) with ply: %d expected %d, Actual %d", test.moves, test.ply, test.expected, n)
		}
	}
}
*/

/*
func TestNumMovesFromFEN(t *testing.T) {
	for _, test := range numMovesFromFENTests {
		board := GetBoardFromFen(test.fen)
		for _, moveStr := range test.moves {
			err := board.MoveLongAlgebraic(moveStr)
			if err != nil {
				t.Errorf(err.Error())
			}
		}
		n := board.GetNumberOfMoves(test.ply, true)
		if n != test.expected {
			t.Errorf("Fen(%s) + moves: %v with ply: %d expected %d, Actual %d", test.fen, test.moves, test.ply, test.expected, n)
		}
	}
}
*/

/*
func TestIsLegal(t *testing.T) {
	for _, test := range legalMovesTests {
		board := GetBoardFromFen(test.fen)
		_, err := board.getMoveFromLongAlgebraic(test.moveStr)
		if err != nil && test.legal {
			t.Errorf("should be legal but has algebraic error for %s with error %s", test.moveStr, err)
		} else if !test.legal && err == nil {
			t.Errorf("should be illegal but has no algebraic error for %s", test.moveStr)
		}
	}
}

func TestFen(t *testing.T) {
	// not an actual FEN
	expected := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq c2 0 1"
	board := GetBoardFromFen("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq c2 0 1")
	actual := board.getFen()
	if actual != expected {
		t.Errorf("FEN expected: %s, actual: %s", expected, actual)
	}
}
*/

/*
func TestHalfMoves(t *testing.T) {
	startFEN := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
	for _, test := range halfMovesTests {
		board := GetBoardFromFen(startFEN)
		for _, moveStr := range test.moves {
			err := board.MoveLongAlgebraic(moveStr)
			if err != nil {
				t.Errorf(err.Error())
			}
		}
		if board.halfMoves != test.expected {
			t.Errorf("Half moves expected: %d, actual: %d", test.expected, board.halfMoves)
		}
	}
}

func BenchmarkNumMoves(b *testing.B) {
	startFEN := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
	board := GetBoardFromFen(startFEN)
	board.MoveLongAlgebraic("e2-e4")
	for i := 0; i < b.N; i++ {
		board.GetNumberOfMoves(3, false)
	}
}
*/

/*
func BenchmarkNumMoves(b *testing.B) {
	startFEN := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
	board := GetBoardFromFen(startFEN)
	board.MoveLongAlgebraic("e2-e4")
	for i := 0; i < b.N; i++ {
		board.GetNumberOfMoves(3, false)
	}
}
*/

func BenchmarkMove(b *testing.B) {
	startFEN := "8/5r2/8/8/2B5/8/4Q3/8 w - - 0 1"
	board := GetBoardFromFen(startFEN)
	for i := 0; i < b.N; i++ {
		move := board.NewMove(3, 0, 55, false)
		board.Move(&move)
		move = board.NewMove(3, 0, 52, false)
		board.Move(&move)
	}
}
