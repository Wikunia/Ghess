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
	{[]string{}, 6, 119060324},

	// no castle through check
	{[]string{"e2e4", "b7b6", "g1f3", "c8a6", "g2g3", "d7d5", "f1g2", "d5e4"}, 1, 23},
	// no castle if in check
	{[]string{"e2e4", "b7b6", "g1f3", "c8a6", "g2g3", "d7d5", "f1g2", "d5e4", "f3d4", "e7e5", "d2d3", "f8b4"}, 1, 7},
}

type numMovesFromFEN struct {
	fen      string
	moves    []string
	ply      int
	expected int
}

var numMovesFromFENTests = []numMovesFromFEN{
	{"4k2r/5ppp/8/8/8/8/5PPP/4K2R b Kk - 0 1", []string{}, 1, 13},   // castle
	{"4k2r/5ppp/8/8/8/8/5PPP/4K2R w Kk - 0 1", []string{}, 2, 169},  // castle
	{"4k2r/5pp1/8/6Pp/8/8/6PP/4K2R w K h6 0 1", []string{}, 2, 169}, // castle + en passant
	{"4k2r/5pp1/8/6Pp/8/8/6PP/4K2R w K - 0 1", []string{}, 2, 156},  // castle + no en passant
	{"3R4/8/8/6K1/8/4k3/8/5Q2 b - - 0 1", []string{}, 1, 1},         // don't allow black to walk into check
	{"3r4/8/8/6k1/8/4K3/8/5q2 w - - 0 1", []string{}, 1, 1},         // don't allow white to walk into check
	{"4k2r/8/8/8/5R2/3K4/8/5Q2 b k - 0 1", []string{}, 1, 12},       // don't allow to move through castle check king side
	{"r3k2r/8/8/8/2K2R2/8/8/3Q4 b kq - 0 1", []string{}, 1, 20},     // don't allow to move through castle check both sides
	{"4kb1r/3ppppp/7n/8/8/3P4/4PPPP/2BQKBNR w - - 0 1", []string{}, 2, 294},
	{"4k2r/3ppppp/7n/8/8/3P4/4PPPP/2BQK3 w - - 0 1", []string{}, 2, 325},
	{"4k2r/3pp2p/7n/8/8/3P4/4P3/2BQK3 w - - 0 1", []string{}, 2, 213},
	{"4k2r/3pp2p/7n/8/8/3P4/4P3/2B1K3 w - - 0 1", []string{}, 2, 177},
	{"4k3/3pp2p/7n/8/8/8/4P3/2B1K3 w - - 0 1", []string{}, 2, 138},
	{"4k3/4p2p/7n/8/8/8/8/2B1K3 w - - 0 1", []string{}, 2, 115},
	{"7k/7p/7n/8/8/8/8/2B1K3 w - - 0 1", []string{}, 2, 62},
	{"8/2p5/3p4/KP5r/1R2Pp1k/8/6P1/8 b - e3 0 1", []string{}, 1, 16},
	{"8/2p5/3p4/KP5r/1R3pPk/8/4P3/8 b - g3 0 1", []string{}, 1, 17},
	// https://www.chessprogramming.org/Perft_Results
	// Kiwipete
	{"r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1", []string{}, 3, 97862},   // last without promotion
	{"r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1", []string{}, 4, 4085603}, // first with promotion and queen side castle where rook king side was captured :D
	// position 3
	{"8/2p5/3p4/KP5r/1R3p1k/8/4P1P1/8 w - - 0 1", []string{}, 4, 43238}, // en passant madness
	// position 4
	{"r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1", []string{}, 2, 264},
	{"r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1", []string{}, 3, 9467},
	{"r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1", []string{}, 4, 422333},
	// position 5
	{"rnbq1k1r/pp1Pbppp/2p5/8/2B5/8/PPP1NnPP/RNBQK2R w KQ - 1 8", []string{}, 4, 2103487},
	// position 6
	{"r4rk1/1pp1qppp/p1np1n2/2b1p1B1/2B1P1b1/P1NP1N2/1PP1QPPP/R4RK1 w - - 0 10", []string{}, 4, 3894594},
}
