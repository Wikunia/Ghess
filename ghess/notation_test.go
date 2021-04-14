package ghess

type standardAlgebraic struct {
	fen      string
	move     string
	expected string
}

var standardAlgebraicTests = []standardAlgebraic{
	{"5k2/8/3N1N2/8/4r3/2N3N1/4PPPP/4K2R w K - 0 1", "e1g1", "O-O"},
	{"5k2/8/3N1N2/8/4r3/6N1/3NPPPP/4K2R w K - 0 1", "d6e4", "Nd6xe4"},
	{"5k2/8/3N1N2/8/r7/6NP/3NPPP1/4K2R w K - 1 2", "d6e4", "Nd6e4"},
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
	{"2n5/1P4R1/8/1R1b1p2/2P1PK1k/8/1R4r1/6R1 w - - 0 1", "b2b3", "R2b3"},
	{"2n5/1P4R1/8/1R1b1p2/2P1PK1k/8/1R4r1/6R1 w - - 0 1", "b5b3", "R5b3"},
	{"2n5/1P4R1/8/1R1b1p2/2P1PK1k/8/1R4r1/6R1 w - - 0 1", "b2g2", "Rbxg2"},
	{"2n5/1P4R1/8/1R1b1p2/2P1PK1k/8/1R4r1/6R1 w - - 0 1", "g1g2", "R1xg2"},
	{"2n5/1P4R1/8/1R1b1p2/2P1PK1k/8/1R4r1/6R1 w - - 0 1", "g7g2", "R7xg2"},
	{"2n5/1P4R1/8/1R1P1p2/4PK1k/8/1R5r/6R1 w - - 1 2", "b2g2", "Rbg2"},
}
