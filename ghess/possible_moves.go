package ghess

import "fmt"

func (board *Board) GetNumberOfMoves(ply int, isBlacksTurn bool) int {
	n := 0
	var pieceIds [16]int
	if isBlacksTurn {
		pieceIds = board.blackIds
	} else {
		pieceIds = board.whiteIds
	}
	for _, pieceId := range pieceIds {
		moves := board.pieces[pieceId].moves
		numMoves := board.pieces[pieceId].numMoves
		for mId := 0; mId < numMoves; mId++ {
			if ply == 1 {
				n += 1
				continue
			}
			boardPrimitives := board.getBoardPrimitives()
			move := board.NewMove(pieceId, 0, moves[mId])
			fmt.Println("Move: ", move)
			board.Move(&move)
			n += board.GetNumberOfMoves(ply-1, !isBlacksTurn)
			board.reverseMove(&move, &boardPrimitives)
		}
	}
	return n
}
