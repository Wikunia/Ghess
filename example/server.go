package main

import (
	"fmt" // "github.com/Wikunia/Ghess/ghess"

	"github.com/pkg/profile"
)

func main() {

	// ghess.Run()
	p := profile.Start(profile.CPUProfile)
	defer p.Stop()

	var x uint64
	x |= 1 << 2
	x |= 1 << 3
	fmt.Printf("x: %064b\n", x)

	/*
		startFEN := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
		board := ghess.GetBoardFromFen(startFEN)
		// board.MoveLongAlgebraic("e2-e4")
		board.GetNumberOfMoves(4, false)
	*/
}
