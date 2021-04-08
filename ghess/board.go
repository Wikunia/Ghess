package ghess

import "fmt"

const NORTH = -8
const SOUTH = 8
const WEST = -1
const EAST = 1
const NORTH_EAST = -7
const NORTH_WEST = -9
const SOUTH_EAST = 9
const SOUTH_WEST = 7

const NORTH_ID = 0
const SOUTH_ID = 1
const WEST_ID = 2
const EAST_ID = 3

// getMovesTilEdge returns an 2d 64x8 array of number of squares from a starting square in the 8 directions
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

// hasBlackPieceOn returns whether there is a black piece on pos
func (board *Board) hasBlackPieceOn(pos int) bool {
	var posB uint64 = 1 << pos
	return posB&board.blackPiecePosB != 0
}

// hasWhitePieceOn returns whether there is a white piece on pos
func (board *Board) hasWhitePieceOn(pos int) bool {
	var posB uint64 = 1 << pos
	return posB&board.whitePiecePosB != 0
}

func (board *Board) sameColoredPieceOn(piece *Piece, pos int) bool {
	if piece.isBlack {
		return board.hasBlackPieceOn(pos)
	} else {
		return board.hasWhitePieceOn(pos)
	}
}

func (board *Board) oppositeColoredPieceOn(piece *Piece, pos int) bool {
	if piece.isBlack {
		return board.hasWhitePieceOn(pos)
	} else {
		return board.hasBlackPieceOn(pos)
	}
}

// combinePositionsOf combies the positions of all specified pieces and outputs the combined positions as uint64
func (board *Board) combinePositionsOf(pieceIds [16]int) uint64 {
	var posB uint64
	for _, pieceId := range pieceIds {
		posB |= board.pieces[pieceId].posB
	}
	return posB
}

// combineMovementsOf combies the movements of all specified pieces and outputs the combined movements as uint64
func (board *Board) combineMovementsOf(pieceIds [16]int) uint64 {
	var posB uint64
	for _, pieceId := range pieceIds {
		posB |= board.pieces[pieceId].movementB
	}
	return posB
}

// setMovement updates the movement of all pieces
func (board *Board) setMovement() {
	for pieceId := range board.pieces {
		board.pieces[pieceId].numMoves = 0
		// piece is not on board
		if board.pieces[pieceId].posB == 0 {
			continue
		}
		// reset current movement
		board.pieces[pieceId].movementB = 0
		switch board.pieces[pieceId].pieceType {
		case BISHOP, ROOK, QUEEN:
			board.setSlidingpieceMovement(&board.pieces[pieceId])
		case KNIGHT:
			board.setKnightMovement(&board.pieces[pieceId])
		case PAWN:
			board.setPawnMovement(&board.pieces[pieceId])
		case KING:
			board.setKingMovement(&board.pieces[pieceId])
		}
	}

	board.whitePieceMovB = board.combineMovementsOf(board.whiteIds)
	board.blackPieceMovB = board.combineMovementsOf(board.blackIds)
}

// setSlidingpieceMovement sets the possible movements for a queen, rook or bishop (does not check if it's a right piece)
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

	for dirId := startDir; dirId < endDir; dirId++ {
		dir := directions[dirId]
		for stepFactor := 1; stepFactor <= board.movesTilEdge[piece.pos][dirId]; stepFactor++ {
			step := stepFactor * dir
			pos := piece.pos + step
			if (piece.isBlack && board.hasBlackPieceOn(pos)) || (!piece.isBlack && board.hasWhitePieceOn(pos)) {
				break
			}
			board.setMovementIfNotCheck(piece, pos)
			// capture
			if (!piece.isBlack && board.hasBlackPieceOn(pos)) || (piece.isBlack && board.hasWhitePieceOn(pos)) {
				break
			}
		}
	}
}

// setKnightMovement sets the possible movements for a knight
func (board *Board) setKnightMovement(piece *Piece) {
	dirSouth := [8]int{2, 2, 1, 1, -1, -1, -2, -2}
	dirEast := [8]int{-1, 1, 2, -2, -2, 2, -1, 1}

	for dirId := 0; dirId < 8; dirId++ {
		dirS := dirSouth[dirId]
		dirE := dirEast[dirId]
		pos := piece.pos + dirS*SOUTH + dirE*EAST

		if dirS > 0 && dirE > 0 { // jump south east
			if board.movesTilEdge[piece.pos][SOUTH_ID] >= dirS && board.movesTilEdge[piece.pos][EAST_ID] >= dirE && !board.sameColoredPieceOn(piece, pos) {
				board.setMovementIfNotCheck(piece, pos)
			}
		} else if dirS > 0 && dirE < 0 { // jump south west
			if board.movesTilEdge[piece.pos][SOUTH_ID] >= dirS && board.movesTilEdge[piece.pos][WEST_ID] >= -dirE && !board.sameColoredPieceOn(piece, pos) {
				board.setMovementIfNotCheck(piece, pos)
			}
		} else if dirS < 0 && dirE > 0 { // jump north east
			if board.movesTilEdge[piece.pos][NORTH_ID] >= -dirS && board.movesTilEdge[piece.pos][EAST_ID] >= dirE && !board.sameColoredPieceOn(piece, pos) {
				board.setMovementIfNotCheck(piece, pos)
			}
		} else if dirS < 0 && dirE < 0 { // jump north west
			if board.movesTilEdge[piece.pos][NORTH_ID] >= -dirS && board.movesTilEdge[piece.pos][WEST_ID] >= -dirE && !board.sameColoredPieceOn(piece, pos) {
				board.setMovementIfNotCheck(piece, pos)
			}
		}
	}
}

// setPawnMovement sets the possible movements for a pawn
func (board *Board) setPawnMovement(piece *Piece) {
	forwardID := NORTH_ID
	forward := NORTH
	startRank := 6
	if piece.isBlack {
		forwardID = SOUTH_ID
		forward = SOUTH
		startRank = 1
	}
	_, rank := xy(piece.pos)

	// one move forward
	if board.movesTilEdge[piece.pos][forwardID] >= 1 && board.pos2PieceId[piece.pos+forward] == 0 {
		board.setMovementIfNotCheck(piece, (piece.pos + forward))
	}
	// two steps forward
	if rank == startRank && board.pos2PieceId[piece.pos+2*forward] == 0 && board.pos2PieceId[piece.pos+forward] == 0 {
		board.setMovementIfNotCheck(piece, (piece.pos + 2*forward))
	}

	// normal capture forward east
	if board.oppositeColoredPieceOn(piece, piece.pos+forward+EAST) && board.movesTilEdge[piece.pos][EAST_ID] >= 1 {
		board.setMovementIfNotCheck(piece, (piece.pos + forward + EAST))
	}
	// normal capture forward west
	if board.oppositeColoredPieceOn(piece, piece.pos+forward+WEST) && board.movesTilEdge[piece.pos][WEST_ID] >= 1 {
		board.setMovementIfNotCheck(piece, (piece.pos + forward + WEST))
	}
	// en passant capture
	if board.en_passant_pos == -1 {
		return
	}
	if piece.pos+forward+EAST == board.en_passant_pos {
		board.setMovementIfNotCheck(piece, (piece.pos + forward + EAST))
	} else if piece.pos+forward+WEST == board.en_passant_pos {
		board.setMovementIfNotCheck(piece, (piece.pos + forward + WEST))
	}
}

// setKingMovement sets the possible movements for a king
func (board *Board) setKingMovement(piece *Piece) {
	directions := [8]int{NORTH, SOUTH, WEST, EAST, NORTH_EAST, NORTH_WEST, SOUTH_EAST, SOUTH_WEST}

	// normal movement
	for dirId := 0; dirId < 8; dirId++ {
		if board.movesTilEdge[piece.pos][dirId] >= 1 {
			pos := piece.pos + directions[dirId]
			if !board.sameColoredPieceOn(piece, pos) {
				board.setMovementIfNotCheck(piece, (pos))
			}
		}
	}

	// castle
	if piece.isBlack {
		if board.black_castle_king {
			// check if positions are free
			if board.pos2PieceId[piece.pos+EAST] == 0 && board.pos2PieceId[piece.pos+2*EAST] == 0 {
				board.setMovementIfNotCheck(piece, (piece.pos + 2*EAST))
			}
		}
		if board.black_castle_queen {
			// check if positions are free
			if board.pos2PieceId[piece.pos+WEST] == 0 && board.pos2PieceId[piece.pos+2*WEST] == 0 && board.pos2PieceId[piece.pos+3*WEST] == 0 {
				board.setMovementIfNotCheck(piece, (piece.pos + 3*WEST))
			}
		}
	} else {
		// and for white
		// todo: refactor
		if board.white_castle_king {
			// check if positions are free
			if board.pos2PieceId[piece.pos+EAST] == 0 && board.pos2PieceId[piece.pos+2*EAST] == 0 {
				board.setMovementIfNotCheck(piece, (piece.pos + 2*EAST))
			}
		}
		if board.white_castle_queen {
			// check if positions are free
			if board.pos2PieceId[piece.pos+WEST] == 0 && board.pos2PieceId[piece.pos+2*WEST] == 0 && board.pos2PieceId[piece.pos+3*WEST] == 0 {
				board.setMovementIfNotCheck(piece, (piece.pos + 3*WEST))
			}
		}
	}
}

// NewMove creates a move object given a pieceId, to and checks whether the move is a capture. If isCapture is set to true
func (board *Board) NewMove(pieceId int, captureId int, to int) Move {
	from := board.pieces[pieceId].pos
	if captureId != 0 {
		to = board.pieces[captureId].pos
	} else {
		if board.pos2PieceId[to] != 0 { // fill capture if there is a piece on that position
			captureId = board.pos2PieceId[to]
		} else if to == board.en_passant_pos && board.pieces[pieceId].pieceType == PAWN {
			if board.isBlacksTurn {
				captureId = board.pos2PieceId[to-8]
			} else {
				captureId = board.pos2PieceId[to+8]
			}
		}
	}
	return Move{pieceId: pieceId, captureId: captureId, from: from, to: to}
}

func (board *Board) TempMove(m *Move) Move {
	forward := NORTH
	if board.pieces[m.pieceId].isBlack {
		forward = SOUTH
	}

	if m.captureId != 0 {
		// important for en passant
		board.pos2PieceId[board.pieces[m.captureId].pos] = 0
		board.pieces[m.captureId].pos = -1
		board.pieces[m.captureId].posB = 0
	}
	board.pieces[m.pieceId].pos = m.to
	board.pieces[m.pieceId].posB = 1 << m.to
	board.pos2PieceId[m.from] = 0
	board.pos2PieceId[m.to] = m.pieceId

	board.whitePiecePosB = board.combinePositionsOf(board.whiteIds)
	board.blackPiecePosB = board.combinePositionsOf(board.blackIds)

	if board.pieces[m.pieceId].pieceType == PAWN && (m.to-m.from) == 2*forward {
		board.en_passant_pos = m.from + forward
	} else {
		board.en_passant_pos = -1
	}

	// check if castled
	piece := board.pieces[m.pieceId]
	isCastle := false
	rookMove := Move{}
	// should not trigger for rverse castle
	x, _ := xy(m.from)
	if piece.pieceType == KING && x == 4 {
		if m.from+2*EAST == m.to {
			isCastle = true
			rookMove = board.NewMove(board.pos2PieceId[m.from+3*EAST], 0, m.from+EAST)
		} else if m.from+2*WEST == m.to {
			isCastle = true
			rookMove = board.NewMove(board.pos2PieceId[m.from+4*WEST], 0, m.from+WEST)
		}
	}
	if isCastle {
		board.TempMove(&rookMove)
	}
	return rookMove
}

func (board *Board) Move(m *Move) Move {
	rookMove := board.TempMove(m)

	board.updateCastleRights(m)

	board.setMovement()

	board.isBlacksTurn = !board.isBlacksTurn
	return rookMove
}

func (board *Board) reverseMove(m *Move, boardPrimitives *BoardPrimitives) {
	move := board.NewMove(m.pieceId, 0, m.from)
	fmt.Println("m.PieceId: ", m.pieceId)
	fmt.Println("m.captureId: ", m.captureId)
	fmt.Println("boardPrimitives.en_passant_pos: ", boardPrimitives.en_passant_pos)
	fmt.Println("m.to ", m.to)
	fmt.Println("====================================")
	if m.captureId != 0 {
		// en passant capture
		if boardPrimitives.en_passant_pos == m.to && board.pieces[m.pieceId].pieceType == PAWN {
			posOfCapturedPawn := 0
			if board.pieces[m.pieceId].isBlack {
				// white pawn was captured
				posOfCapturedPawn = boardPrimitives.en_passant_pos - 8
			} else {
				// black pawn was captured
				posOfCapturedPawn = boardPrimitives.en_passant_pos + 8
			}
			fmt.Println("pos of captured pawn: ", posOfCapturedPawn)
			board.pieces[m.captureId].pos = posOfCapturedPawn
			board.pieces[m.captureId].posB = 1 << posOfCapturedPawn
			board.pos2PieceId[posOfCapturedPawn] = m.captureId
		} else {
			board.pieces[m.captureId].pos = m.to
			board.pieces[m.captureId].posB = 1 << m.to
			board.pos2PieceId[m.to] = m.captureId
		}
	}
	if board.pieces[m.pieceId].pieceType == KING {
		_, y := xy(m.from)
		// king side
		if m.to-m.from == 2 {
			rookId := board.pos2PieceId[y*8+5]
			reverseCastleMove := board.NewMove(rookId, 0, y*8+7)
			board.TempMove(&reverseCastleMove)
		} else if m.to-m.from == -2 { // queen side
			rookId := board.pos2PieceId[y*8+3]
			reverseCastleMove := board.NewMove(rookId, 0, y*8)
			board.TempMove(&reverseCastleMove)
		}
	}

	board.TempMove(&move)

	// resets castle and en passant rights which is important for setMovement
	board.setBoardPrimitives(boardPrimitives)
	board.setMovement()
}

func (board *Board) isLegal(m *Move) bool {
	piece := board.pieces[m.pieceId]
	if piece.isBlack == board.isBlacksTurn {
		return piece.canMoveTo(m.to)
	}
	return false
}

func (board *Board) setMovementIfNotCheck(piece *Piece, pos int) {
	piece.movementB |= 1 << pos
	piece.moves[piece.numMoves] = pos
	piece.numMoves++
}

func (board *Board) updateCastleRights(m *Move) {
	// if king moved remove the right for both sides
	piece := board.pieces[m.pieceId]
	if piece.pieceType == KING {
		if piece.isBlack {
			board.black_castle_king = false
			board.black_castle_queen = false
		} else {
			board.white_castle_king = false
			board.white_castle_queen = false
		}
	}
	// if rook moves remove the castle right for that side
	if piece.pieceType == ROOK {
		file, _ := xy(m.from)
		if file == 7 {
			if piece.isBlack {
				board.black_castle_king = false
			} else {
				board.white_castle_king = false
			}
		} else if file == 0 {
			if piece.isBlack {
				board.black_castle_queen = false
			} else {
				board.white_castle_queen = false
			}
		}
	}

	// if rook gets captured
	if m.captureId > 0 {
		capturedPiece := board.pieces[m.captureId]
		file, _ := xy(m.to)
		if capturedPiece.pieceType == ROOK {
			if capturedPiece.isBlack {
				if file == 7 {
					board.black_castle_king = false
				} else {
					board.black_castle_queen = false
				}
			} else {
				if file == 7 {
					board.white_castle_king = false
				} else {
					board.white_castle_queen = false
				}
			}
		}
	}
}