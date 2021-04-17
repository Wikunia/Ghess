package ghess

import (
	"testing"
)

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
		n := board.GetNumberOfMoves(test.ply)
		if n != test.expected {
			t.Errorf("Moves(%v) with ply: %d expected %d, Actual %d", test.moves, test.ply, test.expected, n)
		}
	}
}

func TestNumMovesFromFEN(t *testing.T) {
	for _, test := range numMovesFromFENTests {
		board := GetBoardFromFen(test.fen)
		for _, moveStr := range test.moves {
			err := board.MoveLongAlgebraic(moveStr)
			if err != nil {
				t.Errorf(err.Error())
			}
		}
		n := board.GetNumberOfMoves(test.ply)
		if n != test.expected {
			t.Errorf("Fen(%s) + moves: %v with ply: %d expected %d, Actual %d", test.fen, test.moves, test.ply, test.expected, n)
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

func TestNextMoves(t *testing.T) {
	startFEN := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
	for _, test := range nextMovesTests {
		board := GetBoardFromFen(startFEN)
		for _, moveStr := range test.moves {
			err := board.MoveLongAlgebraic(moveStr)
			if err != nil {
				t.Errorf(err.Error())
			}
		}
		if board.nextMove != test.expected {
			t.Errorf("Next move expected: %d, actual: %d", test.expected, board.nextMove)
		}
	}
}

func TestNotation(t *testing.T) {
	for _, test := range standardAlgebraicTests {
		board := GetBoardFromFen(test.fen)
		move, err := board.GetMoveFromLongAlgebraic(test.move)
		if err != nil {
			t.Errorf(err.Error())
		}
		moveStr := board.getStandardAlgebraicFromMove(&move)
		if moveStr != test.expected {
			t.Errorf("Standard Algebraic Notation expected %s actually %s", test.expected, moveStr)
		}
	}
}

func TestBits2Array(t *testing.T) {
	startFEN := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
	board := GetBoardFromFen(startFEN)
	// should be the ranks 4+5 (0 based from top)
	arr := bits2array(board.whitePieceMovB)
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			if i == 4 || i == 5 {
				if !arr[i][j] {
					t.Errorf("White has no vision on %d,%d but should have", i, j)
				}
			} else {
				if arr[i][j] {
					t.Errorf("White has vision on %d,%d but shouldn't have", i, j)
				}
			}
		}
	}
	printBits(board.whitePieceMovB)
}

func TestEngines(t *testing.T) {
	move := Move{}
	for _, test := range engineMovesTests {
		board := GetBoardFromFen(test.fen)
		switch test.engineName {
		case "random":
			move = board.randomEngineMove()
		case "captureRandom":
			move = board.captureEngineMove()
		case "checkCaptureRandom":
			move = board.checkCaptureEngineMove()
		case "alphaBeta":
			move = board.AlphaBetaEngineMove()
		}
		algebraic := GetAlgebraicFromMove(&move)
		found := false
		for _, moveStr := range test.possible {
			if algebraic == moveStr {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("The move %s is not a possible engine move for engine %s", algebraic, test.engineName)
		}
	}
}

func BenchmarkNumMove(b *testing.B) {
	startFEN := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
	board := GetBoardFromFen(startFEN)

	for i := 0; i < b.N; i++ {
		board.GetNumberOfMoves(5)
	}
}
