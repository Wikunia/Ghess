package ghess

type numMoves struct {
	moves    []string
	ply      int
	expected int
}

var numMovesTests = []numMoves{
	{[]string{}, 1, 20},
	{[]string{}, 2, 400},
	{[]string{}, 4, 197281},
	{[]string{}, 5, 4865609},
	// check castling
	{[]string{"e2-e4", "e7-e5", "f1-c4", "f8-c5", "g1-f3", "g8-f6"}, 2, 1052},
}
