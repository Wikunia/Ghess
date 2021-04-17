package ghess

type engineMoves struct {
	fen        string
	engineName string
	possible   []string
}

var engineMovesTests = []engineMoves{
	{"8/8/4k3/8/4K3/p7/8/8 w - - 0 1", "random", []string{"e4d4", "e4d3", "e4e3", "e4f3", "e4f4"}},
	{"8/8/4k3/8/4K3/8/p7/8 w - - 0 1", "alphaBeta", []string{"e4d4", "e4d3", "e4e3", "e4f3", "e4f4"}},
	{"8/8/4k3/8/4K3/p7/8/8 w - - 0 1", "captureRandom", []string{"e4d4", "e4d3", "e4e3", "e4f3", "e4f4"}},
	{"8/8/4k3/8/4K3/p7/8/8 w - - 0 1", "checkCaptureRandom", []string{"e4d4", "e4d3", "e4e3", "e4f3", "e4f4"}},
	{"8/8/4k3/8/4K3/6r1/8/5N2 w - - 0 1", "random", []string{"e4d4", "e4f4", "f1d2", "f1e3", "f1g3", "f1h2"}},
	{"8/8/4k3/8/4K3/6r1/8/5N2 w - - 0 1", "captureRandom", []string{"f1g3"}},
	{"8/8/4k3/8/4K3/6r1/8/5N2 w - - 0 1", "checkCaptureRandom", []string{"f1g3"}},
	{"8/8/B3k3/8/4K3/6r1/8/5N2 w - - 0 1", "checkCaptureRandom", []string{"a6c8", "a6c4"}},
	{"8/8/B3k3/8/4K3/6r1/8/5N2 w - - 0 1", "checkCaptureRandom", []string{"a6c8", "a6c4"}},
	{"5k2/1P6/2P2K2/8/8/8/8/8 w - - 0 1", "random", []string{"b7b8q", "b7b8r", "b7b8b", "b7b8n", "c6c7", "f6e6", "f6e5", "f6f5", "f6g5", "f6g6"}},
	{"5k2/1P6/2P2K2/8/8/8/8/8 w - - 0 1", "captureRandom", []string{"b7b8q", "b7b8r", "b7b8b", "b7b8n", "c6c7", "f6e6", "f6e5", "f6ef5", "f6g5", "f6g6"}},
	{"5k2/1P6/2P2K2/8/8/8/8/8 w - - 0 1", "checkCaptureRandom", []string{"b7b8q", "b7b8r"}},
	{"5k2/1P6/2P2K2/8/8/8/8/8 w - - 0 1", "alphaBeta", []string{"b7b8q", "b7b8r"}},
}
