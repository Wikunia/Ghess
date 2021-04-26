package ghess

type halfMoves struct {
	moves    []string
	expected int
}

var halfMovesTests = []halfMoves{
	{[]string{"e2e4"}, 0},
	{[]string{"e2e4", "b7b6", "g1f3", "c8a6", "g2g3", "d7d5", "f1g2", "d5e4", "f3d4", "e7e5", "d2d3", "f8b4"}, 1},
}

type nextMoves struct {
	moves    []string
	expected int
}

var nextMovesTests = []nextMoves{
	{[]string{"e2e4"}, 1},
	{[]string{"e2e4", "b7b6", "g1f3", "c8a6", "g2g3", "d7d5", "f1g2", "d5e4", "f3d4", "e7e5", "d2d3", "f8b4"}, 7},
}

type staticEvaluationStruct struct {
	fen      string
	expected float64
}

var staticEvaluationTests = []staticEvaluationStruct{
	{"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1", 0.0},
	{"rnbqkbnr/pppppppp/8/8/3P4/8/PPP1PPPP/RNBQKBNR w KQkq - 0 1", 4.0},
	{"rnbqkb1r/pppppppp/5n2/8/3P4/8/PPP1PPPP/RNBQKBNR w KQkq - 0 1", 2.0},
	{"rnbqkb1r/pppppppp/5n2/8/3P1B2/8/PPP1PPPP/RN1QKBNR w KQkq - 0 1", 5.0},
}
