package ghess

type halfMoves struct {
	moves    []string
	expected int
}

var halfMovesTests = []halfMoves{
	{[]string{"e2-e4"}, 0},
	{[]string{"e2-e4", "b7-b6", "g1-f3", "c8-a6", "g2-g3", "d7-d5", "f1-g2", "d5-e4", "f3-d4", "e7-e5", "d2-d3", "f8-b4"}, 1},
}

type reverseMoves struct {
	fen string
}

var reverseMovesTests = []reverseMoves{
	{"r3k2r/p2pqpb1/bn2pnp1/1BpPN3/1p2P3/2N2Q1p/PPPB1PPP/R3K2R w KQkq c6 0 2"},
	// {"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"},
}
