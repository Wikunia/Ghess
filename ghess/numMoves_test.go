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
	// {[]string{}, 5, 4865609},
	// {[]string{"f2-f3"}, 4, 178889},
	{[]string{"f2-f3", "a7-a6"}, 3, 7697},
	{[]string{"f2-f3", "b7-b6"}, 3, 8505},
	{[]string{"f2-f3", "c7-c6"}, 3, 8407},
	{[]string{"f2-f3", "d7-d6"}, 3, 10902},
	{[]string{"f2-f3", "e7-e6"}, 3, 11632},
	{[]string{"f2-f3", "f7-f6"}, 3, 7697},
	{[]string{"f2-f3", "g7-g6"}, 3, 8504},
	{[]string{"f2-f3", "h7-h6"}, 3, 7697},
	{[]string{"f2-f3", "a7-a5"}, 3, 8490},
	{[]string{"f2-f3", "b7-b5"}, 3, 8492},
	{[]string{"f2-f3", "c7-c5"}, 3, 8871},
	{[]string{"f2-f3", "d7-d5"}, 3, 11334},
	{[]string{"f2-f3", "e7-e5"}, 3, 11679},
	{[]string{"f2-f3", "f7-f5"}, 3, 8124},
	{[]string{"f2-f3", "g7-g5"}, 3, 8507},
	{[]string{"f2-f3", "h7-h5"}, 3, 8490},
	{[]string{"f2-f3", "b8-a6"}, 3, 8086},
	{[]string{"f2-f3", "b8-c6"}, 3, 8877},
	{[]string{"f2-f3", "g8-f6"}, 3, 8834},
	{[]string{"f2-f3", "g8-h6"}, 3, 8064},
	{[]string{"f2-f3", "f7-f5", "e1-f2"}, 2, 437},
	{[]string{"f2-f3", "f7-f5", "e1-f2", "f5-f4"}, 1, 19},
	{[]string{"f2-f3", "f7-f5", "e1-f2", "a7-a6"}, 1, 22},
	{[]string{"f2-f3", "f7-f5", "e1-f2", "b7-b6"}, 1, 22},
	{[]string{"f2-f3", "f7-f5", "e1-f2", "c7-c6"}, 1, 22},
	{[]string{"f2-f3", "f7-f5", "e1-f2", "d7-d6"}, 1, 22},
	{[]string{"f2-f3", "f7-f5", "e1-f2", "e7-e6"}, 1, 22},
	{[]string{"f2-f3", "f7-f5", "e1-f2", "g7-g6"}, 1, 22},
	{[]string{"f2-f3", "f7-f5", "e1-f2", "h7-h6"}, 1, 22},
	{[]string{"f2-f3", "f7-f5", "e1-f2", "a7-a5"}, 1, 22},
	{[]string{"f2-f3", "f7-f5", "e1-f2", "b7-b5"}, 1, 22},
	{[]string{"f2-f3", "f7-f5", "e1-f2", "c7-c5"}, 1, 22},
	{[]string{"f2-f3", "f7-f5", "e1-f2", "d7-d5"}, 1, 22},
	{[]string{"f2-f3", "f7-f5", "e1-f2", "e7-e5"}, 1, 22},
	{[]string{"f2-f3", "f7-f5", "e1-f2", "g7-g5"}, 1, 22},
	{[]string{"f2-f3", "f7-f5", "e1-f2", "h7-h5"}, 1, 22},
	{[]string{"f2-f3", "f7-f5", "e1-f2", "b8-a6"}, 1, 22},
	{[]string{"f2-f3", "f7-f5", "e1-f2", "b8-c6"}, 1, 22},
	{[]string{"f2-f3", "f7-f5", "e1-f2", "g8-f6"}, 1, 22},
	{[]string{"f2-f3", "f7-f5", "e1-f2", "g8-h6"}, 1, 22},
	{[]string{"f2-f3", "f7-f5", "e1-f2", "e8-f7"}, 1, 22},
	{[]string{"f2-f3", "g7-g5", "a2-a3"}, 2, 378},
	{[]string{"f2-f3", "g7-g5", "b2-b3"}, 2, 420},
	{[]string{"f2-f3", "g7-g5", "c2-c3"}, 2, 420},
	{[]string{"f2-f3", "g7-g5", "d2-d3"}, 2, 526},
	{[]string{"f2-f3", "g7-g5", "e2-e3"}, 2, 545},
	{[]string{"f2-f3", "g7-g5", "g2-g3"}, 2, 420},
	{[]string{"f2-f3", "g7-g5", "h2-h3"}, 2, 379},
	{[]string{"f2-f3", "g7-g5", "f3-f4"}, 2, 458},
	{[]string{"f2-f3", "g7-g5", "a2-a4"}, 2, 420},
	{[]string{"f2-f3", "g7-g5", "b2-b4"}, 2, 421},
	{[]string{"f2-f3", "g7-g5", "c2-c4"}, 2, 442},
	{[]string{"f2-f3", "g7-g5", "d2-d4"}, 2, 548},
	{[]string{"f2-f3", "g7-g5", "e2-e4"}, 2, 546},
	{[]string{"f2-f3", "g7-g5", "g2-g4"}, 2, 382},
	{[]string{"f2-f3", "g7-g5", "h2-h4"}, 2, 459},
	{[]string{"f2-f3", "g7-g5", "b1-a3"}, 2, 399},
	{[]string{"f2-f3", "g7-g5", "b1-c3"}, 2, 441},
	{[]string{"f2-f3", "g7-g5", "g1-h3"}, 2, 441},
	{[]string{"f2-f3", "g7-g5", "e1-f2"}, 2, 462},
	{[]string{"f2-f3", "e7-e5", "a2-a3"}, 2, 522},
	{[]string{"f2-f3", "e7-e5", "b2-b3"}, 2, 575},
	{[]string{"f2-f3", "e7-e5", "c2-c3"}, 2, 579},
	{[]string{"f2-f3", "e7-e5", "d2-d3"}, 2, 732},
	{[]string{"f2-f3", "e7-e5", "e2-e3"}, 2, 751},
	{[]string{"f2-f3", "e7-e5", "g2-g3"}, 2, 594},
	{[]string{"f2-f3", "e7-e5", "h2-h3"}, 2, 518},
	{[]string{"f2-f3", "e7-e5", "f3-f4"}, 2, 623},
	{[]string{"f2-f3", "e7-e5", "a2-a4"}, 2, 578},
	{[]string{"f2-f3", "e7-e5", "b2-b4"}, 2, 559},
	{[]string{"f2-f3", "e7-e5", "c2-c4"}, 2, 605},
	{[]string{"f2-f3", "e7-e5", "d2-d4"}, 2, 815},
	{[]string{"f2-f3", "e7-e5", "e2-e4"}, 2, 698},
	{[]string{"f2-f3", "e7-e5", "g2-g4"}, 2, 575},
	{[]string{"f2-f3", "e7-e5", "h2-h4"}, 2, 578},
	{[]string{"f2-f3", "e7-e5", "b1-a3"}, 2, 546},
	{[]string{"f2-f3", "e7-e5", "b1-c3"}, 2, 607},
	{[]string{"f2-f3", "e7-e5", "g1-h3"}, 2, 606},
	{[]string{"f2-f3", "e7-e5", "e1-f2"}, 2, 618},
	{[]string{"f2-f3", "e7-e5", "e1-f2", "e5-e4"}, 1, 22},
	{[]string{"f2-f3", "e7-e5", "e1-f2", "a7-a6"}, 1, 22},
	{[]string{"f2-f3", "e7-e5", "e1-f2", "b7-b6"}, 1, 22},
	{[]string{"f2-f3", "e7-e5", "e1-f2", "c7-c6"}, 1, 22},
	{[]string{"f2-f3", "e7-e5", "e1-f2", "d7-d6"}, 1, 22},
	{[]string{"f2-f3", "e7-e5", "e1-f2", "f7-f6"}, 1, 22},
	{[]string{"f2-f3", "e7-e5", "e1-f2", "g7-g6"}, 1, 22},
	{[]string{"f2-f3", "e7-e5", "e1-f2", "h7-h6"}, 1, 22},
	{[]string{"f2-f3", "e7-e5", "e1-f2", "a7-a5"}, 1, 22},
	{[]string{"f2-f3", "e7-e5", "e1-f2", "b7-b5"}, 1, 22},
	{[]string{"f2-f3", "e7-e5", "e1-f2", "c7-c5"}, 1, 22},
	{[]string{"f2-f3", "e7-e5", "e1-f2", "d7-d5"}, 1, 22},
	{[]string{"f2-f3", "e7-e5", "e1-f2", "f7-f5"}, 1, 22},
	{[]string{"f2-f3", "e7-e5", "e1-f2", "g7-g5"}, 1, 22},
	{[]string{"f2-f3", "e7-e5", "e1-f2", "h7-h5"}, 1, 22},
	{[]string{"f2-f3", "e7-e5", "e1-f2", "b8-a6"}, 1, 22},
	{[]string{"f2-f3", "e7-e5", "e1-f2", "b8-c6"}, 1, 22},
	{[]string{"f2-f3", "e7-e5", "e1-f2", "g8-f6"}, 1, 22},
	{[]string{"f2-f3", "e7-e5", "e1-f2", "g8-h6"}, 1, 22},
	{[]string{"f2-f3", "e7-e5", "e1-f2", "g8-e7"}, 1, 22},
	{[]string{"f2-f3", "e7-e5", "e1-f2", "f8-a3"}, 1, 21},
	{[]string{"f2-f3", "e7-e5", "e1-f2", "f8-b4"}, 1, 21},
	{[]string{"f2-f3", "e7-e5", "e1-f2", "f8-c5"}, 1, 4},
	{[]string{"f2-f3", "e7-e5", "e1-f2", "f8-d6"}, 1, 22},
	{[]string{"f2-f3", "e7-e5", "e1-f2", "f8-e7"}, 1, 22},
	{[]string{"f2-f3", "e7-e5", "e1-f2", "d8-h4"}, 1, 2},
	{[]string{"f2-f3", "e7-e5", "e1-f2", "d8-g5"}, 1, 20},
	{[]string{"f2-f3", "e7-e5", "e1-f2", "d8-f6"}, 1, 22},
	{[]string{"f2-f3", "e7-e5", "e1-f2", "d8-e7"}, 1, 22},
	{[]string{"f2-f3", "e7-e5", "e1-f2", "e8-e7"}, 1, 22},

	// no castle through check
	// {[]string{"e2-e4", "b7-b6", "g1-f3", "c8-a6", "g2-g3", "d7-d5", "f1-g2", "d5-e4"}, 1, 23},
	// no castle if in check
	// {[]string{"e2-e4", "b7-b6", "g1-f3", "c8-a6", "g2-g3", "d7-d5", "f1-g2", "d5-e4", "f3-d4", "e7-e5", "d2-d3", "f8-b4"}, 1, 7},
}

type numMovesFromFEN struct {
	fen      string
	moves    []string
	ply      int
	expected int
}

var numMovesFromFENTests = []numMovesFromFEN{
	// {"4k2r/5ppp/8/8/8/8/5PPP/4K2R b Kk - 0 1", []string{}, 1, 13},   // castle
	// {"4k2r/5ppp/8/8/8/8/5PPP/4K2R w Kk - 0 1", []string{}, 2, 169},  // castle
	// {"4k2r/5pp1/8/6Pp/8/8/6PP/4K2R w K h6 0 1", []string{}, 2, 169}, // castle + en passant
	// {"4k2r/5pp1/8/6Pp/8/8/6PP/4K2R w K - 0 1", []string{}, 2, 156},  // castle + no en passant
	// {"3R4/8/8/6K1/8/4k3/8/5Q2 b - - 0 1", []string{}, 1, 1},         // don't allow black to walk into check
	// {"3r4/8/8/6k1/8/4K3/8/5q2 w - - 0 1", []string{}, 1, 1},         // don't allow white to walk into check
	// {"4k2r/8/8/8/5R2/3K4/8/5Q2 b k - 0 1", []string{}, 1, 12},       // don't allow to move through castle check king side
	// {"r3k2r/8/8/8/2K2R2/8/8/3Q4 b kq - 0 1", []string{}, 1, 20},     // don't allow to move through castle check both sides
	// {"4kb1r/3ppppp/7n/8/8/3P4/4PPPP/2BQKBNR w - - 0 1", []string{}, 2, 294},
	// {"4k2r/3ppppp/7n/8/8/3P4/4PPPP/2BQK3 w - - 0 1", []string{}, 2, 325},
	// {"4k2r/3pp2p/7n/8/8/3P4/4P3/2BQK3 w - - 0 1", []string{}, 2, 213},
	// {"4k2r/3pp2p/7n/8/8/3P4/4P3/2B1K3 w - - 0 1", []string{}, 2, 177},
	// {"4k3/3pp2p/7n/8/8/8/4P3/2B1K3 w - - 0 1", []string{}, 2, 138},
	// {"4k3/4p2p/7n/8/8/8/8/2B1K3 w - - 0 1", []string{}, 2, 115},
	{"7k/7p/7n/8/8/8/8/2B1K3 w - - 0 1", []string{}, 2, 62},
}
