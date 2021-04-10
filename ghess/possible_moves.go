package ghess

func (board *Board) getNumberOfMoves(startPly, ply int, isBlacksTurn bool) int {
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
			move, _ := board.NewMove(pieceId, 0, moves[mId], 0)
			/*
				padding := (startPly - ply) * 2
				for i := 0; i < padding; i++ {
					fmt.Print(" ")
				}
				fmt.Println(getAlgebraicFromMove(&move))
			*/
			if ply == 1 {
				n += 1
				continue
			}
			boardPrimitives := board.getBoardPrimitives()
			board.Move(&move)
			n += board.getNumberOfMoves(startPly, ply-1, !isBlacksTurn)
			board.reverseMove(&move, &boardPrimitives)
		}
	}
	return n
}

func (board *Board) GetNumberOfMoves(ply int) int {
	return board.getNumberOfMoves(ply, ply, board.isBlacksTurn)
}
