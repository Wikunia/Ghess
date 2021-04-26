package main

import (
	"github.com/Wikunia/Ghess/ghess"
	// "github.com/pkg/profile"
)

func main() {

	ghess.Run()
	/*
		p := profile.Start(profile.CPUProfile)
		defer p.Stop()

		board := ghess.GetBoardFromFen(ghess.START_FEN)
		board.GetNumberOfMoves(5)
	*/
}
