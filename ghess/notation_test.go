package ghess

type standardAlgebraic struct {
	fen      string
	move     string
	expected string
}

var standardAlgebraicTests = []standardAlgebraic{
	{"5k2/8/3N1N2/8/4r3/2N3N1/4PPPP/4K2R w K - 0 1", "e1g1", "O-O"},
	{"5k2/8/3N1N2/8/4r3/6N1/3NPPPP/4K2R w K - 0 1", "d6e4", "Nd6xe4"},
	{"5k2/8/3N1N2/8/4r3/6N1/3NPPPP/4K2R w K - 0 1", "g3e4", "Ngxe4"},
	{"5k2/8/3N1N2/8/4r3/6N1/3NPPPP/R3K2R w KQ - 0 1", "e1c1", "O-O-O"},
}
