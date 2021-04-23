package ghess

import (
	"math/rand"
)

func (board *Board) captureEngineMove() Move {
	possibleMoves := []Move{}
	captureMoves := []Move{}
	n := 0
	nCaptures := 0
	var PieceIds [16]int
	if board.IsBlacksTurn {
		PieceIds = board.blackIds
	} else {
		PieceIds = board.whiteIds
	}
	for _, PieceId := range PieceIds {
		moves := board.pieces[PieceId].moves
		numMoves := board.pieces[PieceId].numMoves
		for mId := 0; mId < numMoves; mId++ {
			numTinyMoves := 1
			_, isPromotion := board.NewMove(PieceId, 0, moves[mId], 0)
			if isPromotion {
				numTinyMoves = 4
			}
			x := 0
			for i := 0; i < numTinyMoves; i++ {
				if isPromotion {
					x++
				}
				move, _ := board.NewMove(PieceId, 0, moves[mId], x)
				if move.captureId != 0 {
					captureMoves = append(captureMoves, move)
					nCaptures++
				}
				possibleMoves = append(possibleMoves, move)
				n++
			}
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
