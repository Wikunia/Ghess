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
	{"r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1", []string{"g1h1"}, 3, 81638},
	{"r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1", []string{"c4c5"}, 3, 60769},
	{"r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1", []string{"d2d4"}, 3, 72051},
	{"r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1", []string{"f3d4"}, 3, 75736},
	{"r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1", []string{"b4c5"}, 3, 58167},
	{"r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1", []string{"f1f2"}, 3, 73972},
	{"r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1", []string{"b4c5", "c7c6"}, 2, 1527},
	{"r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1", []string{"b4c5", "d7d6"}, 2, 1567},
	{"r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1", []string{"b4c5", "d7d5"}, 2, 1603},
	{"r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1", []string{"b4c5", "b2a1q"}, 2, 1478},
	{"r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1", []string{"b4c5", "b2a1r"}, 2, 1353},
	{"r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1", []string{"b4c5", "b2a1b"}, 2, 1346},
	{"r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1", []string{"b4c5", "b2a1n"}, 2, 1285},
	{"r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1", []string{"b4c5", "b2b1q"}, 2, 1559},
	{"r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1", []string{"b4c5", "b2b1r"}, 2, 1439},
	{"r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1", []string{"b4c5", "b2b1b"}, 2, 1340},
	{"r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1", []string{"b4c5", "b2b1n"}, 2, 1282},
	{"r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1", []string{"b4c5", "g7h6"}, 2, 1288},
	{"r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1", []string{"b4c5", "a5b3"}, 2, 1408},
	{"r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1", []string{"b4c5", "a5c4"}, 2, 1561},
	{"r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1", []string{"b4c5", "a5c6"}, 2, 1665},
	{"r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1", []string{"b4c5", "f6e4"}, 2, 1468},
	{"r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1", []string{"b4c5", "f6g4"}, 2, 1456},
	{"r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1", []string{"b4c5", "f6d5"}, 2, 1598},
	{"r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1", []string{"b4c5", "f6h5"}, 2, 1412},
	{"r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1", []string{"b4c5", "f6g8"}, 2, 1381},
	{"r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1", []string{"b4c5", "b6c5"}, 2, 203},
	{"r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1", []string{"b4c5", "b6a7"}, 2, 1574},
	{"r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1", []string{"b4c5", "g6e4"}, 2, 1581},
	{"r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1", []string{"b4c5", "g6f5"}, 2, 1622},
	{"r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1", []string{"b4c5", "g6h5"}, 2, 1488},
	{"r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1", []string{"b4c5", "a8a7"}, 2, 1319},
	{"r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1", []string{"b4c5", "a8b8"}, 2, 1620},
	{"r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1", []string{"b4c5", "a8c8"}, 2, 1546},
	{"r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1", []string{"b4c5", "a8d8"}, 2, 1509},
	{"r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1", []string{"b4c5", "h8f8"}, 2, 1446},
	{"r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1", []string{"b4c5", "h8g8"}, 2, 1419},
	{"r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1", []string{"b4c5", "a3a2"}, 2, 1367},
	{"r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1", []string{"b4c5", "a3b3"}, 2, 1527},
	{"r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1", []string{"b4c5", "a3c3"}, 2, 1624},
	{"r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1", []string{"b4c5", "a3d3"}, 2, 1596},
	{"r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1", []string{"b4c5", "a3e3"}, 2, 168},
	{"r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1", []string{"b4c5", "a3f3"}, 2, 1552},
	{"r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1", []string{"b4c5", "a3a4"}, 2, 1396},
	{"r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1", []string{"b4c5", "a3b4"}, 2, 1463},
	{"r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1", []string{"b4c5", "a3c5"}, 2, 197},
	{"r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1", []string{"b4c5", "e8c8"}, 2, 1495},
	{"r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1", []string{"b4c5", "e8d8"}, 2, 1439},
	{"r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1", []string{"b4c5", "b2a1r", "d2d3"}, 1, 38},
	{"r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1", []string{"b4c5", "b2a1r", "g2g3"}, 1, 40},
	{"r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1", []string{"b4c5", "b2a1r", "h2h3"}, 1, 40},
	{"r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1", []string{"b4c5", "b2a1r", "e4e5"}, 1, 43},
	{"r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1", []string{"b4c5", "b2a1r", "d2d4"}, 1, 40},
	{"r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1", []string{"b4c5", "b2a1r", "g2g4"}, 1, 40},
	{"r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1", []string{"b4c5", "b2a1r", "h2h4"}, 1, 40},
	{"r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1", []string{"b4c5", "b2a1r", "f3e1"}, 1, 42},
	{"r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1", []string{"b4c5", "b2a1r", "f3d4"}, 1, 42},
	{"r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1", []string{"b4c5", "b2a1r", "f3h4"}, 1, 42},
	{"r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1", []string{"b4c5", "b2a1r", "f3e5"}, 1, 42},
	{"r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1", []string{"b4c5", "b2a1r", "f3g5"}, 1, 42},
	{"r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1", []string{"b4c5", "b2a1r", "h6g4"}, 1, 41},
	{"r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1", []string{"b4c5", "b2a1r", "h6f5"}, 1, 40},
	{"r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1", []string{"b4c5", "b2a1r", "h6f7"}, 1, 41},
	{"r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1", []string{"b4c5", "b2a1r", "h6g8"}, 1, 40},
	{"r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1", []string{"b4c5", "b2a1r", "a4c2"}, 1, 40},
	{"r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1", []string{"b4c5", "b2a1r", "a4b3"}, 1, 36},
	{"r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1", []string{"b4c5", "b2a1r", "c5f2"}, 1, 49},
	{"r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1", []string{"b4c5", "b2a1r", "g1f2"}, 1, 40},
	{"r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1", []string{"b4c5", "b2a1r", "c5e3"}, 1, 47},
	{"r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1", []string{"b4c5", "b2a1r", "g1h1"}, 1, 40},
	{"r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1", []string{"b4c5", "b2a1r", "c5d4"}, 1, 47},
	{"r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1", []string{"b4c5", "b2a1r", "c5b6"}, 1, 45},
	{"r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1", []string{"b4c5", "b2a1r", "d1b3"}, 1, 38},
	{"r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1", []string{"b4c5", "b2a1r", "d1e2"}, 1, 42},
	{"r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1", []string{"b4c5", "b2a1r", "d1c2"}, 1, 42},
	{"r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1", []string{"b4c5", "b2a1r", "f1e1"}, 1, 40},
	{"r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1", []string{"b4c5", "b2a1r", "f1f2"}, 1, 40},
	{"r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1", []string{"b4c5", "b2a1r", "d1a1"}, 1, 36},
	{"r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1", []string{"b4c5", "b2a1r", "d1b1"}, 1, 38},
	{"r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1", []string{"b4c5", "b2a1r", "d1c1"}, 1, 39},
	{"r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1", []string{"b4c5", "b2a1r", "d1e1"}, 1, 41},
}
