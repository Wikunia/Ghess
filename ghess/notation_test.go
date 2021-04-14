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
	{"8/1P6/8/8/5K1k/8/8/8 w - - 0 1", "b7b8q", "b8=Q"},
	{"8/1P6/8/8/5K1k/8/8/8 w - - 0 1", "b7b8n", "b8=N"},
	{"8/1P6/8/8/5K1k/8/8/8 w - - 0 1", "b7b8b", "b8=B"},
	{"8/1P6/8/8/5K1k/8/8/8 w - - 0 1", "b7b8r", "b8=R"},
	{"8/1P6/8/8/5K1k/8/8/8 w - - 0 1", "f4f3", "Kf3"},
	{"8/1P6/8/5p2/5K1k/8/8/8 w - - 0 1", "f4f5", "Kxf5"},
	{"2n5/1P6/8/5p2/5K1k/8/8/8 w - - 0 1", "b7c8q", "bxc8=Q"},
	{"2n5/1P6/8/5p2/5K1k/8/8/8 w - - 0 1", "b7c8b", "bxc8=B"},
	{"2n5/1P6/8/5p2/5K1k/8/8/8 w - - 0 1", "b7c8n", "bxc8=N"},
	{"2n5/1P6/8/5p2/5K1k/8/8/8 w - - 0 1", "b7c8r", "bxc8=R"},
	{"2n5/1P6/8/3b1p2/2P1PK1k/8/8/8 w - - 0 1", "c4d5", "cxd5"},
}
