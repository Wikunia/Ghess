package main

import (
	"github.com/Wikunia/Ghess/ghess"
	"github.com/pkg/profile"
)

func main() {

	// ghess.Run()
	p := profile.Start(profile.CPUProfile)
	defer p.Stop()

	startFEN := "8/5r2/8/8/2B5/8/4Q3/8 w - - 0 1"
	board := ghess.GetBoardFromFen(startFEN)
	move := board.NewMove(3, 0, 55, false)
	board.Move(&move)
}
