package ghess

import "math/bits"

const NORTH = -8
const SOUTH = 8
const WEST = -1
const EAST = 1
const NORTH_EAST = -7
const NORTH_WEST = -9
const SOUTH_EAST = 9
const SOUTH_WEST = 7

func getMovesTilEdge() [64][8]int {
	var movesTilEdge [64][8]int
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			p := i*8 + j
			movesTilEdge[p] = [8]int{
				i,             // North
				7 - i,         // South
				j,             // West
				7 - j,         // East
				min(i, 7-j),   // NorthEast
				min(i, j),     // NorthWest
				min(7-i, 7-j), // SouthEast
				min(7-i, j),   // SouthWest
			}
		}
	}
	return movesTilEdge
}

func (board *Board) hasBlackPieceOn(pos int) bool {
	var posB uint64 = 1 << pos
	return posB&board.blackPiecePosB != 0
}

func (board *Board) hasWhitePieceOn(pos int) bool {
	var posB uint64 = 1 << pos
	return posB&board.whitePiecePosB != 0
}

// combinePositionsOf combies the positions of all specified pieces and outputs the combined position as uint64
func (board *Board) combinePositionsOf(pieceIds [16]int) uint64 {
	var posB uint64
	for _, pieceId := range pieceIds {
		posB |= board.pieces[pieceId].posB
	}
	return posB
}

func (board *Board) setMovement() {
	for pieceId := range board.pieces {
		// piece is not on board
		if board.pieces[pieceId].posB == 0 {
			continue
		}
		switch board.pieces[pieceId].pieceType {
		case BISHOP, ROOK, QUEEN:
			board.setSlidingpieceMovement(&board.pieces[pieceId])
		}
	}
}

func (board *Board) setSlidingpieceMovement(piece *Piece) {
	directions := [8]int{NORTH, SOUTH, WEST, EAST, NORTH_EAST, NORTH_WEST, SOUTH_EAST, SOUTH_WEST}
	startDir := 0
	endDir := 8
	switch piece.pieceType {
	case ROOK:
		endDir = 4
	case BISHOP:
		startDir = 4
	}
	// reset current movement
	piece.movementB = 0
	for dirId := startDir; dirId < endDir; dirId++ {
		dir := directions[dirId]
		for stepFactor := 1; stepFactor <= board.movesTilEdge[piece.pos][dirId]; stepFactor++ {
			step := stepFactor * dir
			pos := piece.pos + step
			if (piece.isBlack && board.hasBlackPieceOn(pos)) || (!piece.isBlack && board.hasWhitePieceOn(pos)) {
				break
			}
			piece.movementB |= 1 << pos
			// capture
			if (!piece.isBlack && board.hasBlackPieceOn(pos)) || (piece.isBlack && board.hasWhitePieceOn(pos)) {
				break
			}
		}
	}
}

func (board *Board) NewMove(pieceId int, captureId int, to int, isCapture bool) Move {
	from := board.pieces[pieceId].pos
	if isCapture {
		to = board.pieces[captureId].pos
	} else if captureId == 0 {
		if board.pos2PieceId[to] != 0 { // fill capture if there is a piece on that position
			captureId = board.pos2PieceId[to]
		} else if to == bits.TrailingZeros64(board.en_passant_posB) && board.pieces[pieceId].pieceType == PAWN {
			if board.isBlacksTurn {
				captureId = board.pos2PieceId[to-8]
			} else {
				captureId = board.pos2PieceId[to+8]
			}
		}
	}
	return Move{pieceId: pieceId, captureId: captureId, from: from, to: to}
}

func (board *Board) Move(m *Move) {
	if m.captureId != 0 {
		board.pieces[m.captureId].pos = -1
		board.pieces[m.captureId].posB = 0
	}
	board.pieces[m.pieceId].pos = m.to
	var posB uint64 = 1 << m.to
	board.pieces[m.pieceId].posB = posB
	board.pos2PieceId[m.from] = 0
	board.pos2PieceId[m.to] = m.pieceId

	board.whitePiecePosB = board.combinePositionsOf(board.whiteIds)
	board.blackPiecePosB = board.combinePositionsOf(board.blackIds)

	board.setMovement()
}
