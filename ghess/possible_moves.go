package ghess

func (board *Board) getPossibleMoves() []Move {
	var possibleMoves []Move

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
				possibleMoves = append(possibleMoves, move)
			}
		}
	}
	return possibleMoves
}

func (board *Board) getNumberOfMoves(startPly, ply int, IsBlacksTurn bool) int {
	n := 0
	var PieceIds [16]int
	if IsBlacksTurn {
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
			/*
				padding := (startPly - ply) * 2
				for i := 0; i < padding; i++ {
					fmt.Print(" ")
				}
				fmt.Println(GetAlgebraicFromMove(&move))
			*/
			x := 0
			for i := 0; i < numTinyMoves; i++ {
				if ply == 1 {
					n += 1
					continue
				}
				if isPromotion {
					x++
				}
				move, _ := board.NewMove(PieceId, 0, moves[mId], x)
				boardPrimitives := board.getBoardPrimitives()
				board.Move(&move)
				n += board.getNumberOfMoves(startPly, ply-1, !IsBlacksTurn)
				board.reverseMove(&move, &boardPrimitives)
			}
		}
	}
	return n
}

func (board *Board) GetNumberOfMoves(ply int) int {
	return board.getNumberOfMoves(ply, ply, board.IsBlacksTurn)
}
