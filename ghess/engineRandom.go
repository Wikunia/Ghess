package ghess

import (
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
			numTinyMoves := 1
			_, isPromotion := board.NewMove(pieceId, 0, moves[mId], 0)
			if isPromotion {
				numTinyMoves = 4
			}
			x := 0
			for i := 0; i < numTinyMoves; i++ {
				if isPromotion {
					x++
				}
				move, _ := board.NewMove(pieceId, 0, moves[mId], x)
				possibleMoves = append(possibleMoves, move)
				n++
			}
		}
	}
	moveId := rand.Intn(n)
	return possibleMoves[moveId]
}
