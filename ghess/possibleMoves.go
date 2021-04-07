package ghess

import (
	"fmt"
)

func (board *Board) getNumberOfMovesInternal(ply int, startPly int, print bool) int {
	possibleMoves := board.getPossibleMoves()
	if ply == 1 {
		return len(possibleMoves)
	}
	if len(possibleMoves) == 0 {
		return 0
	}

	n := 0
	for _, m := range possibleMoves {
		fromY := board.pieces[m.pieceId].position.y
		fromX := board.pieces[m.pieceId].position.x
		boardPrimitives := board.getBoardPrimitives()
		capturedId, castledMove := board.move(&board.pieces[m.pieceId], m.toY, m.toX)
		numMoves := board.getNumberOfMovesInternal(ply-1, startPly, print)
		if ply == startPly && print {
			fmt.Printf("Num moves starting with %s : %d\n", board.moveToToLongAlgebraic(fromY, fromX, m.toY, m.toX), numMoves)
		}
		n += numMoves
		board.reverseMove(&board.pieces[m.pieceId], fromY, fromX, capturedId, castledMove.pieceId, &boardPrimitives)
	}
	return n
}

func (board *Board) GetNumberOfMoves(ply int, print bool) int {
	return board.getNumberOfMovesInternal(ply, ply, print)
}

func (board *Board) getPossibleMoves() []Move {
	var moves []Move
	for pieceId := range board.pieces {
		if !board.pieces[pieceId].onBoard {
			continue
		}
		// only the current isBlack can move
		if board.pieces[pieceId].isBlack != board.isBlack {
			continue
		}
		positions := board.pieces[pieceId].moves[:board.pieces[pieceId].numMoves]
		for _, p := range positions {
			moves = append(moves, board.newMove(board.pieces[pieceId].id, 0, p.y, p.x, false))
		}
	}
	return moves
}
