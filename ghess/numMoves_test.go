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
	// {[]string{}, 4, 197281},
	// {[]string{}, 5, 4865609},

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
	{"4k2r/5ppp/8/8/8/8/5PPP/4K2R b Kk - 0 1", []string{}, 1, 13},   // castle
	{"4k2r/5ppp/8/8/8/8/5PPP/4K2R w Kk - 0 1", []string{}, 2, 169},  // castle
	{"4k2r/5pp1/8/6Pp/8/8/6PP/4K2R w K h6 0 1", []string{}, 2, 169}, // castle + en passant
	{"4k2r/5pp1/8/6Pp/8/8/6PP/4K2R w K - 0 1", []string{}, 2, 156},  // castle + no en passant
	{"3R4/8/8/6K1/8/4k3/8/5Q2 b - - 0 1", []string{}, 1, 1},         // don't allow black to walk into check
	{"3r4/8/8/6k1/8/4K3/8/5q2 w - - 0 1", []string{}, 1, 1},         // don't allow white to walk into check
}
