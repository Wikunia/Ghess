package ghess

import (
	"math/rand"
)

func (board *Board) captureEngineMove() Move {
	possibleMoves := []Move{}
	captureMoves := []Move{}
	n := 0
	nCaptures := 0
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
			if move.captureId != 0 {
				captureMoves = append(captureMoves, move)
				nCaptures++
			}
			possibleMoves = append(possibleMoves, move)
			n++
		}
	}
	// first capture something
	if nCaptures != 0 {
		moveId := rand.Intn(nCaptures)
		return captureMoves[moveId]
	}
	// else just a move...

	moveId := rand.Intn(n)
	return possibleMoves[moveId]
}
