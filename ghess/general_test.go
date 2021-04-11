package ghess

type halfMoves struct {
	moves    []string
	expected int
}

var halfMovesTests = []halfMoves{
	{[]string{"e2e4"}, 0},
	{[]string{"e2e4", "b7b6", "g1f3", "c8a6", "g2g3", "d7d5", "f1g2", "d5e4", "f3d4", "e7e5", "d2d3", "f8b4"}, 1},
}

type legalMoves struct {
	fen     string
	moveStr string
	legal   bool
}

var legalMovesTests = []legalMoves{
	{"r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 0", "a2a4", true},
	{"r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 0", "e1g1", true},
	{"r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 0", "a6e2", false},
}
