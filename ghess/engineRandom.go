package ghess

import (
	"fmt"
	"math/rand"
)

func (board *Board) randomEngineMove() Move {
	possibleMoves := []Move{}
	n := 0
	var pieceIds [16]int
	if board.isBlacksTurn {
		pieceIds = board.blackIds
	} else {
		pieceIds = board.whiteIds
	}
	for _, pieceId := range pieceIds {
		moves := board.pieces[pieceId].moves
		numMoves := board.pieces[pieceId].numMoves
		for mId := 0; mId < numMoves; mId++ {
			move, _ := board.NewMove(pieceId, 0, moves[mId], 0)
			possibleMoves = append(possibleMoves, move)
			n++
		}
	}
	moveId := rand.Intn(n)
	fmt.Println("moveId: ", moveId)
	fmt.Println("possibleMoves: ", n)
	return possibleMoves[moveId]
}
