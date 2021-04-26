package ghess

import (
	"math/rand"
)

// checkCaptureEngineMove tries to check (if the piece checking can't be captured, or double check) then capture and then random
func (board *Board) checkCaptureEngineMove() Move {
	possibleMoves := []Move{}
	captureMoves := []Move{}
	checkMoves := []Move{}
	highestMaterialGain := -100000
	n := 0
	nCaptures := 0
	nCheck := 0
	var PieceIds [16]int
	if board.IsBlacksTurn {
		PieceIds = board.blackIds
	} else {
		PieceIds = board.whiteIds
	}
	myColor := board.IsBlacksTurn
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
				boardPrimitives := board.getBoardPrimitives()
				board.Move(&move)
				if board.check {
					if !board.oppositeHasVisionOn(&board.pieces[PieceId], move.to) || board.doubleCheck {
						checkMoves = append(checkMoves, move)
						nCheck++
					} else if move.captureId == 0 {
						// don't choose this move
						board.reverseMove(&move, &boardPrimitives)
						continue
					}
				}

				if move.captureId != 0 {
					cMaterialGain := board.countMaterialOfColor(myColor) - board.countMaterialOfColor(!myColor)
					if cMaterialGain > highestMaterialGain {
						highestMaterialGain = cMaterialGain
						captureMoves = []Move{move}
						nCaptures = 1
					} else if cMaterialGain == highestMaterialGain {
						captureMoves = append(captureMoves, move)
						nCaptures++
					}
				}
				possibleMoves = append(possibleMoves, move)
				n++

				board.reverseMove(&move, &boardPrimitives)
			}
		}
	}
	// first check
	if nCheck != 0 {
		moveId := rand.Intn(nCheck)
		return checkMoves[moveId]
	}

	// second capture something
	if nCaptures != 0 {
		moveId := rand.Intn(nCaptures)
		return captureMoves[moveId]
	}
	// else just a move...
	moveId := rand.Intn(n)
	return possibleMoves[moveId]
}
