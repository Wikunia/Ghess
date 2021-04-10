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
