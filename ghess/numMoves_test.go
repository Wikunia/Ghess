package ghess

type numMoves struct {
	moves    []string
	ply      int
	expected int
}

var numMovesTests = []numMoves{
	{[]string{}, 1, 20},
	{[]string{}, 2, 400},
	{[]string{}, 3, 8902},
	{[]string{}, 4, 197281},
	{[]string{}, 5, 4865609},

	// no castle through check
	{[]string{"e2-e4", "b7-b6", "g1-f3", "c8-a6", "g2-g3", "d7-d5", "f1-g2", "d5-e4"}, 1, 23},
	// no castle if in check
	{[]string{"e2-e4", "b7-b6", "g1-f3", "c8-a6", "g2-g3", "d7-d5", "f1-g2", "d5-e4", "f3-d4", "e7-e5", "d2-d3", "f8-b4"}, 1, 7},
}

type numMovesFromFEN struct {
	fen      string
	moves    []string
	ply      int
	expected int
}

var numMovesFromFENTests = []numMovesFromFEN{
	// {"r1bqr1k1/ppp2ppp/2p2n2/2b5/4P3/2NP3P/PPP1BPP1/R1BQK2R b KQ - 2 8", []string{}, 2, 1270},
	// {"6k1/p1r5/1p1p4/3Pp1B1/1P2P3/P5P1/6K1/2R5 b - - 2 36", []string{}, 3, 10483},
	// {"r5r1/5k1p/1b3p1B/p1p1p2p/PpP1P3/1P1PR2P/6P1/R5K1 w - - 0 33", []string{}, 2, 641},
	// {"r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 0", []string{}, 3, 97862},
	{"r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 0", []string{"a2-a4"}, 2, 2149},
	// {"r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 0", []string{"a2-a4"}, 1, 44},
	// {"r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 0", []string{"a2-a4", "b4-b3"}, 1, 49},
	// {"r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 0", []string{"a2-a4", "g6-g5"}, 1, 48},
	// {"r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 0", []string{"a2-a4", "c7-c6"}, 1, 50},
	// {"r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 0", []string{"a2-a4", "d7-d6"}, 1, 48},
	// {"r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 0", []string{"a2-a4", "c7-c5"}, 1, 50},
	// {"r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 0", []string{"a2-a4", "h3-g2"}, 1, 47},
	// {"r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 0", []string{"a2-a4", "e6-d5"}, 1, 49},
	// {"r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 0", []string{"a2-a4", "b4-c3"}, 1, 49},
	// {"r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 0", []string{"a2-a4", "b4-a3"}, 1, 51},
	// {"r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 0", []string{"a2-a4", "b6-a4"}, 1, 50},
	// {"r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 0", []string{"a2-a4", "b6-c4"}, 1, 47},
	// {"r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 0", []string{"a2-a4", "b6-d5"}, 1, 49},
	// {"r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 0", []string{"a2-a4", "b6-c8"}, 1, 49},
	// {"r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 0", []string{"a2-a4", "f6-e4"}, 1, 52},
	// {"r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 0", []string{"a2-a4", "f6-g4"}, 1, 48},
	// {"r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 0", []string{"a2-a4", "f6-d5"}, 1, 50},
	// {"r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 0", []string{"a2-a4", "f6-h5"}, 1, 50},
	// {"r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 0", []string{"a2-a4", "f6-h7"}, 1, 50},
	// {"r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 0", []string{"a2-a4", "f6-g8"}, 1, 50},
	// {"r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 0", []string{"a2-a4", "a6-e2"}, 1, 42},
	// {"r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 0", []string{"a2-a4", "a6-d3"}, 1, 47},
	// {"r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 0", []string{"a2-a4", "a6-c4"}, 1, 47},
	// {"r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 0", []string{"a2-a4", "a6-b5"}, 1, 49},
	// {"r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 0", []string{"a2-a4", "a6-b7"}, 1, 49},
	// {"r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 0", []string{"a2-a4", "a6-c8"}, 1, 49},
	// {"r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 0", []string{"a2-a4", "g7-h6"}, 1, 49},
	// {"r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 0", []string{"a2-a4", "g7-f8"}, 1, 49},
	// {"r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 0", []string{"a2-a4", "a8-b8"}, 1, 49},
	// {"r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 0", []string{"a2-a4", "a8-c8"}, 1, 49},
	// {"r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 0", []string{"a2-a4", "a8-d8"}, 1, 49},
	// {"r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 0", []string{"a2-a4", "h8-h4"}, 1, 49},
	// {"r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 0", []string{"a2-a4", "h8-h5"}, 1, 49},
	// {"r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 0", []string{"a2-a4", "h8-h6"}, 1, 49},
	// {"r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 0", []string{"a2-a4", "h8-h7"}, 1, 49},
	// {"r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 0", []string{"a2-a4", "h8-f8"}, 1, 49},
	// {"r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 0", []string{"a2-a4", "h8-g8"}, 1, 49},
	// {"r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 0", []string{"a2-a4", "e7-c5"}, 1, 49},
	// {"r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 0", []string{"a2-a4", "e7-d6"}, 1, 48},
	// {"r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 0", []string{"a2-a4", "e7-d8"}, 1, 49},
	// {"r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 0", []string{"a2-a4", "e7-f8"}, 1, 49},
	// {"r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 0", []string{"a2-a4", "e8-d8"}, 1, 49},
	// {"r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 0", []string{"a2-a4", "e8-f8"}, 1, 49},
	// {"r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 0", []string{"a2-a4", "e8-g8"}, 1, 49},
	// {"r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 0", []string{"a2-a4", "e8-c8"}, 1, 49},
	// {"r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 0", []string{"e2-b5", "b4-b3"}, 1, 54},
	// {"r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 0", []string{"e2-b5", "g6-g5"}, 1, 52},
	// {"r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 0", []string{"e2-b5", "c7-c6"}, 1, 53},
	// {"r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 0", []string{"e2-b5", "e8-c8"}, 1, 53},
	// {"r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 0", []string{"e2-b5", "c7-c5"}, 1, 54},
	// {"r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 0", []string{"e2-b5", "h3-g2"}, 1, 51},
	// {"r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 0", []string{"e2-b5", "e6-d5"}, 1, 53},
	// {"r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 0", []string{"e2-b5", "b4-c3"}, 1, 53},
	// {"r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 0", []string{"e2-b5", "b6-a4"}, 1, 52},
	// {"r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 0", []string{"e2-b5", "b6-c4"}, 1, 50},
	// {"r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 0", []string{"e2-b5", "b6-d5"}, 1, 53},
	// {"r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 0", []string{"e2-b5", "b6-c8"}, 1, 53},
	// {"r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 0", []string{"e2-b5", "f6-e4"}, 1, 56},
	// {"r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 0", []string{"e2-b5", "f6-g4"}, 1, 52},
	// {"r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 0", []string{"e2-b5", "f6-d5"}, 1, 54},
	// {"r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 0", []string{"e2-b5", "f6-h5"}, 1, 54},
	// {"r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 0", []string{"e2-b5", "f6-h7"}, 1, 54},
	// {"r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 0", []string{"e2-b5", "f6-g8"}, 1, 54},
	// {"r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 0", []string{"e2-b5", "a6-b5"}, 1, 43},
	// {"r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 0", []string{"e2-b5", "a6-b7"}, 1, 53},
	// {"r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 0", []string{"e2-b5", "a6-c8"}, 1, 53},
	// {"r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 0", []string{"e2-b5", "g7-h6"}, 1, 53},
	// {"r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 0", []string{"e2-b5", "g7-f8"}, 1, 53},
	// {"r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 0", []string{"e2-b5", "a8-b8"}, 1, 53},
	// {"r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 0", []string{"e2-b5", "a8-c8"}, 1, 53},
	// {"r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 0", []string{"e2-b5", "a8-d8"}, 1, 53},
	// {"r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 0", []string{"e2-b5", "h8-h4"}, 1, 53},
	// {"r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 0", []string{"e2-b5", "h8-h5"}, 1, 53},
	// {"r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 0", []string{"e2-b5", "h8-h6"}, 1, 53},
	// {"r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 0", []string{"e2-b5", "h8-h7"}, 1, 53},
	// {"r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 0", []string{"e2-b5", "h8-f8"}, 1, 53},
	// {"r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 0", []string{"e2-b5", "h8-g8"}, 1, 53},
	// {"r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 0", []string{"e2-b5", "e7-c5"}, 1, 53},
	// {"r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 0", []string{"e2-b5", "e7-d6"}, 1, 52},
	// {"r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 0", []string{"e2-b5", "e7-d8"}, 1, 53},
	// {"r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 0", []string{"e2-b5", "e7-f8"}, 1, 53},
	// {"r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 0", []string{"e2-b5", "e8-d8"}, 1, 53},
	// {"r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 0", []string{"e2-b5", "e8-f8"}, 1, 53},
	// {"r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 0", []string{"e2-b5", "e8-g8"}, 1, 53},
}
