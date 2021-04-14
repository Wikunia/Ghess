package ghess

type engineMoves struct {
	fen        string
	engineName string
	possible   []string
}

var engineMovesTests = []engineMoves{
	{"8/8/4k3/8/4K3/8/8/8 w - - 0 1", "random", []string{"e4d4", "e4d3", "e4e3", "e4f3", "e4f4"}},
	{"8/8/4k3/8/4K3/8/8/8 w - - 0 1", "alphaBeta", []string{"e4d4", "e4d3", "e4e3", "e4f3", "e4f4"}},
	{"8/8/4k3/8/4K3/8/8/8 w - - 0 1", "captureRandom", []string{"e4d4", "e4d3", "e4e3", "e4f3", "e4f4"}},
	{"8/8/4k3/8/4K3/8/8/8 w - - 0 1", "checkCaptureRandom", []string{"e4d4", "e4d3", "e4e3", "e4f3", "e4f4"}},
	{"8/8/4k3/8/4K3/6r1/8/5N2 w - - 0 1", "random", []string{"e4d4", "e4f4", "f1d2", "f1e3", "f1g3", "f1h2"}},
	{"8/8/4k3/8/4K3/6r1/8/5N2 w - - 0 1", "captureRandom", []string{"f1g3"}},
	{"8/8/4k3/8/4K3/6r1/8/5N2 w - - 0 1", "checkCaptureRandom", []string{"f1g3"}},
	{"8/8/B3k3/8/4K3/6r1/8/5N2 w - - 0 1", "checkCaptureRandom", []string{"a6c8", "a6c4"}},
}
