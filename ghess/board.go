package ghess

import (
	"math"
)

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

func (board *Board) oppositeHasVisionOn(piece *Piece, pos int) bool {
	if piece.isBlack {
		// check if white has vision on pos
		return board.whitePieceMovB&(1<<pos) != 0
	} else {
		// check if black has vision on pos
		return board.blackPieceMovB&(1<<pos) != 0
	}
}

func (board *Board) sameHasVisionOn(piece *Piece, pos int) bool {
	if piece.isBlack {
		// check if black has vision on pos
		return board.blackPieceMovB&(1<<pos) != 0
	} else {
		// check if white has vision on pos
		return board.whitePieceMovB&(1<<pos) != 0
	}
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
	// set movement for the color that just moved before setting movement for pieces that can move next
	// this helps with checking the movement of pieces due to check
	lastColor := !board.isBlacksTurn
	pieceIds := make([]int, 32)
	// last was white
	if !lastColor {
		copy(pieceIds[:16], board.whiteIds[:])
		copy(pieceIds[16:], board.blackIds[:])
	} else {
		copy(pieceIds[:16], board.blackIds[:])
		copy(pieceIds[16:], board.whiteIds[:])
	}

	// for the color that last moved we allow them to capture their own pieces
	// this helps as a defense strategy such that the king can't capture a piece in the next move if it's protected
	wasLastColor := false
	// reset several check values
	board.check = false
	board.doubleCheck = false
	board.blockCheckSquaresB = 0

	for pieceId := range board.pieces {
		board.pieces[pieceId].pinnedMoveB = math.MaxUint64
		board.pieces[pieceId].movementB = 0
	}
	board.whitePiecePosB = board.combinePositionsOf(board.whiteIds)
	board.blackPiecePosB = board.combinePositionsOf(board.blackIds)
	board.whitePieceMovB = board.combineMovementsOf(board.whiteIds)
	board.blackPieceMovB = board.combineMovementsOf(board.blackIds)
	// printBits(board.blackPiecePosB)

	for _, pieceId := range pieceIds {
		board.pieces[pieceId].numMoves = 0
		// reset current movement
		board.pieces[pieceId].movementB = 0
		// piece is not on board
		if board.pieces[pieceId].posB == 0 {
			continue
		}
		wasLastColor = board.pieces[pieceId].isBlack == lastColor
		switch board.pieces[pieceId].pieceType {
		case BISHOP, ROOK, QUEEN:
			board.setSlidingpieceMovement(&board.pieces[pieceId], wasLastColor)
		case KNIGHT:
			board.setKnightMovement(&board.pieces[pieceId], wasLastColor)
		case PAWN:
			board.setPawnMovement(&board.pieces[pieceId], wasLastColor)
		case KING:
			board.setKingMovement(&board.pieces[pieceId], wasLastColor)
		}
		if !lastColor {
			board.whitePieceMovB = board.combineMovementsOf(board.whiteIds)
		} else {
			board.blackPieceMovB = board.combineMovementsOf(board.blackIds)
		}
	}
	// update the other color
	if !lastColor {
		board.blackPieceMovB = board.combineMovementsOf(board.blackIds)
	} else {
		board.whitePieceMovB = board.combineMovementsOf(board.whiteIds)
	}
}

// setSlidingpieceMovement sets the possible movements for a queen, rook or bishop (does not check if it's a right piece)
func (board *Board) setSlidingpieceMovement(piece *Piece, wasLastColor bool) {
	directions := [8]int{NORTH, SOUTH, WEST, EAST, NORTH_EAST, NORTH_WEST, SOUTH_EAST, SOUTH_WEST}
	startDir := 0
	endDir := 8
	switch piece.pieceType {
	case ROOK:
		endDir = 4
	case BISHOP:
		startDir = 4
	}
	oppositeKingPos := 0

	// board.isBlacksTurn is already for next move at this stage
	if board.isBlacksTurn {
		oppositeKingPos = board.pieces[board.blackKingId].pos
	} else {
		oppositeKingPos = board.pieces[board.whiteKingId].pos
	}

	for dirId := startDir; dirId < endDir; dirId++ {
		dir := directions[dirId]
		for stepFactor := 1; stepFactor <= board.movesTilEdge[piece.pos][dirId]; stepFactor++ {
			step := stepFactor * dir
			pos := piece.pos + step
			// can't move through our own pieces
			if (piece.isBlack && board.hasBlackPieceOn(pos)) || (!piece.isBlack && board.hasWhitePieceOn(pos)) {
				// we can defend our own pieces though if we just moved -> this makes it impossible that the king can capture a defended piece
				if wasLastColor {
					board.setPieceCanMoveTo(piece, pos)
					// en passant case where the pawn of the color that just moved is in front of the other pawn that can capture it
					// we need to disallow capturing if it would be check otherwise
					if (board.isBlacksTurn && board.en_passant_pos == pos+8) || (!board.isBlacksTurn && board.en_passant_pos == pos-8) {
						// the en passant can be possible if it's either east or west direction
						if dirId == WEST_ID || dirId == EAST_ID {
							// is there a pawn of opposite color at pos+dir ?
							if stepFactor+1 < board.movesTilEdge[piece.pos][dirId] && board.pos2PieceId[pos+dir] != 0 {
								possiblePawnPiece := board.pieces[board.pos2PieceId[pos+dir]]
								if possiblePawnPiece.pieceType == PAWN && piece.isBlack != possiblePawnPiece.isBlack {
									// en passant case fully possible, but is there even a king in sight? :D
									for stepFactorKingHunt := stepFactor + 2; stepFactorKingHunt <= board.movesTilEdge[piece.pos][dirId]; stepFactorKingHunt++ {
										stepKingHunt := stepFactorKingHunt * dir
										posKingHunt := piece.pos + stepKingHunt
										// okay there is another piece in between -> all is safe
										if board.pos2PieceId[posKingHunt] != 0 && posKingHunt != oppositeKingPos {
											break
										}
										// disallow en passant capture
										if posKingHunt == oppositeKingPos {
											board.pieces[board.pos2PieceId[pos+dir]].pinnedMoveB &= ^(1 << board.en_passant_pos)
										}
									}
								}
							}
						}
					}
				}
				break
			}

			board.setPieceCanMoveTo(piece, pos)
			// capture
			if (!piece.isBlack && board.hasBlackPieceOn(pos)) || (piece.isBlack && board.hasWhitePieceOn(pos)) {
				// only do some more checks if we just played the last move
				if wasLastColor {
					// check if we check the king
					if pos == oppositeKingPos {
						if !board.check {
							board.check = true
						} else {
							board.doubleCheck = true
						}
						// add last moves and piece itself to blockCheckSquares
						for stepFactorCheck := 0; stepFactorCheck < stepFactor; stepFactorCheck++ {
							stepCheck := stepFactorCheck * dir
							posCheck := piece.pos + stepCheck
							board.blockCheckSquaresB |= 1 << posCheck
						}
						// add the one square after the king to possible moves of the piece (to avoid letting the king run away in the same direction backwards)
						if stepFactor < board.movesTilEdge[piece.pos][dirId] {
							board.setPieceCanMoveTo(piece, pos+dir)
						}
					}
					enPassantCase := false
					// set pinned pieces if we can catch the king afterwards
					for stepFactorKingHunt := stepFactor + 1; stepFactorKingHunt <= board.movesTilEdge[piece.pos][dirId]; stepFactorKingHunt++ {
						stepKingHunt := stepFactorKingHunt * dir
						posKingHunt := piece.pos + stepKingHunt

						// if there is another piece in between => no pin (besides special en passant case)
						if board.pos2PieceId[posKingHunt] != 0 && posKingHunt != oppositeKingPos {
							// if there is a pawn of the previous color in between which can be captured en passant.
							// Like i.e 8/2p5/3p4/KP5r/1R3pPk/8/4P3/8 b - g3 0 1
							// we need to disallow en passant
							if board.en_passant_pos != -1 {
								if (board.isBlacksTurn && board.en_passant_pos == posKingHunt+8) || (!board.isBlacksTurn && board.en_passant_pos == posKingHunt-8) {
									enPassantCase = true
								}
							}
							if !enPassantCase {
								break
							}
						}

						if posKingHunt == oppositeKingPos {
							// piece at pos is pinned
							pinnedPieceId := board.pos2PieceId[pos]
							if !enPassantCase {
								// reset pinnedMove to be able to set 1s to the positions it can actually move to
								board.pieces[pinnedPieceId].pinnedMoveB = 0
								// set the pinnedMove bitset starting from the piece through the pinned piece up to the king (not including but doesn't matter)
								for stepFactorPin := 0; stepFactorPin < stepFactorKingHunt; stepFactorPin++ {
									stepPin := stepFactorPin * dir
									posPin := piece.pos + stepPin
									board.pieces[pinnedPieceId].pinnedMoveB |= 1 << posPin
								}
							} else {
								// we can move forward or capture but not capture en passant
								board.pieces[pinnedPieceId].pinnedMoveB &= ^(1 << board.en_passant_pos)
							}
							break
						}
					}

				}

				break
			}
		}
	}
}

func (board *Board) setBlockingSquaresIfKingAt(startPos, pos int) {
	// board.isBlacksTurn is already for next move at this stage
	oppositeKingPos := 0
	if board.isBlacksTurn {
		oppositeKingPos = board.pieces[board.blackKingId].pos
	} else {
		oppositeKingPos = board.pieces[board.whiteKingId].pos
	}
	if pos == oppositeKingPos {
		if !board.check {
			board.check = true
		} else {
			board.doubleCheck = true
		}
		// add position of piece to block squares
		board.blockCheckSquaresB |= 1 << startPos
	}
}

// setKnightMovement sets the possible movements for a knight
func (board *Board) setKnightMovement(piece *Piece, wasLastColor bool) {
	dirSouth := [8]int{2, 2, 1, 1, -1, -1, -2, -2}
	dirEast := [8]int{-1, 1, 2, -2, -2, 2, -1, 1}

	for dirId := 0; dirId < 8; dirId++ {
		dirS := dirSouth[dirId]
		dirE := dirEast[dirId]
		pos := piece.pos + dirS*SOUTH + dirE*EAST

		if dirS > 0 && dirE > 0 { // jump south east
			if board.movesTilEdge[piece.pos][SOUTH_ID] >= dirS && board.movesTilEdge[piece.pos][EAST_ID] >= dirE {
				if !board.sameColoredPieceOn(piece, pos) || wasLastColor {
					board.setPieceCanMoveTo(piece, pos)
					board.setBlockingSquaresIfKingAt(piece.pos, pos)
				}
			}
		} else if dirS > 0 && dirE < 0 { // jump south west
			if board.movesTilEdge[piece.pos][SOUTH_ID] >= dirS && board.movesTilEdge[piece.pos][WEST_ID] >= -dirE {
				if !board.sameColoredPieceOn(piece, pos) || wasLastColor {
					board.setPieceCanMoveTo(piece, pos)
					board.setBlockingSquaresIfKingAt(piece.pos, pos)
				}
			}
		} else if dirS < 0 && dirE > 0 { // jump north east
			if board.movesTilEdge[piece.pos][NORTH_ID] >= -dirS && board.movesTilEdge[piece.pos][EAST_ID] >= dirE {
				if !board.sameColoredPieceOn(piece, pos) || wasLastColor {
					board.setPieceCanMoveTo(piece, pos)
					board.setBlockingSquaresIfKingAt(piece.pos, pos)
				}
			}
		} else if dirS < 0 && dirE < 0 { // jump north west
			if board.movesTilEdge[piece.pos][NORTH_ID] >= -dirS && board.movesTilEdge[piece.pos][WEST_ID] >= -dirE {
				if !board.sameColoredPieceOn(piece, pos) || wasLastColor {
					board.setPieceCanMoveTo(piece, pos)
					board.setBlockingSquaresIfKingAt(piece.pos, pos)
				}
			}
		}
	}
}

// setPawnMovement sets the possible movements for a pawn
func (board *Board) setPawnMovement(piece *Piece, wasLastColor bool) {
	forwardID := NORTH_ID
	forward := NORTH
	startRank := 6
	if piece.isBlack {
		forwardID = SOUTH_ID
		forward = SOUTH
		startRank = 1
	}
	_, rank := xy(piece.pos)

	// for the previous color we are only interested in capture vision
	if !wasLastColor {
		// one move forward
		if board.movesTilEdge[piece.pos][forwardID] >= 1 && board.pos2PieceId[piece.pos+forward] == 0 {
			board.setPieceCanMoveTo(piece, (piece.pos + forward))
		}
		// two steps forward
		if rank == startRank && board.pos2PieceId[piece.pos+2*forward] == 0 && board.pos2PieceId[piece.pos+forward] == 0 {
			board.setPieceCanMoveTo(piece, (piece.pos + 2*forward))
		}
	}

	// normal capture forward east
	if board.movesTilEdge[piece.pos][EAST_ID] >= 1 {
		if board.oppositeColoredPieceOn(piece, piece.pos+forward+EAST) || wasLastColor {
			board.setPieceCanMoveTo(piece, piece.pos+forward+EAST)
			board.setBlockingSquaresIfKingAt(piece.pos, piece.pos+forward+EAST)
		}
	}
	// normal capture forward west
	if board.movesTilEdge[piece.pos][WEST_ID] >= 1 {
		if board.oppositeColoredPieceOn(piece, piece.pos+forward+WEST) || wasLastColor {
			board.setPieceCanMoveTo(piece, piece.pos+forward+WEST)
			board.setBlockingSquaresIfKingAt(piece.pos, piece.pos+forward+WEST)
		}
	}
	// en passant capture
	if board.en_passant_pos == -1 {
		return
	}
	if piece.pos+forward+EAST == board.en_passant_pos && board.movesTilEdge[piece.pos][EAST_ID] >= 1 {
		board.setPieceCanMoveTo(piece, (piece.pos + forward + EAST))
	} else if piece.pos+forward+WEST == board.en_passant_pos && board.movesTilEdge[piece.pos][WEST_ID] >= 1 {
		board.setPieceCanMoveTo(piece, (piece.pos + forward + WEST))
	}
}

// setKingMovement sets the possible movements for a king
func (board *Board) setKingMovement(piece *Piece, wasLastColor bool) {
	directions := [8]int{NORTH, SOUTH, WEST, EAST, NORTH_EAST, NORTH_WEST, SOUTH_EAST, SOUTH_WEST}

	// normal movement
	for dirId := 0; dirId < 8; dirId++ {
		if board.movesTilEdge[piece.pos][dirId] >= 1 {
			pos := piece.pos + directions[dirId]
			if (wasLastColor || !board.sameColoredPieceOn(piece, pos)) && !board.oppositeHasVisionOn(piece, pos) {
				board.setPieceCanMoveTo(piece, pos)
			}
		}
	}

	// don't allow castle out of check
	if board.check {
		return
	}

	// castle
	if piece.isBlack {
		if board.black_castle_king {
			// check if positions are free
			if board.pos2PieceId[piece.pos+EAST] == 0 && board.pos2PieceId[piece.pos+2*EAST] == 0 {
				// check that we don't castle through check
				if !board.oppositeHasVisionOn(piece, piece.pos+EAST) && !board.oppositeHasVisionOn(piece, piece.pos+2*EAST) {
					board.setPieceCanMoveTo(piece, (piece.pos + 2*EAST))
				}
			}
		}
		if board.black_castle_queen {
			// check if positions are free
			if board.pos2PieceId[piece.pos+WEST] == 0 && board.pos2PieceId[piece.pos+2*WEST] == 0 && board.pos2PieceId[piece.pos+3*WEST] == 0 {
				if !board.oppositeHasVisionOn(piece, piece.pos+WEST) && !board.oppositeHasVisionOn(piece, piece.pos+2*WEST) {
					board.setPieceCanMoveTo(piece, (piece.pos + 2*WEST))
				}
			}
		}
	} else {
		// and for white
		// todo: refactor
		if board.white_castle_king {
			// check if positions are free
			if board.pos2PieceId[piece.pos+EAST] == 0 && board.pos2PieceId[piece.pos+2*EAST] == 0 {
				if !board.oppositeHasVisionOn(piece, piece.pos+EAST) && !board.oppositeHasVisionOn(piece, piece.pos+2*EAST) {
					board.setPieceCanMoveTo(piece, (piece.pos + 2*EAST))
				}
			}
		}
		if board.white_castle_queen {
			// check if positions are free
			if board.pos2PieceId[piece.pos+WEST] == 0 && board.pos2PieceId[piece.pos+2*WEST] == 0 && board.pos2PieceId[piece.pos+3*WEST] == 0 {
				if !board.oppositeHasVisionOn(piece, piece.pos+WEST) && !board.oppositeHasVisionOn(piece, piece.pos+2*WEST) {
					board.setPieceCanMoveTo(piece, (piece.pos + 2*WEST))
				}
			}
		}
	}
}

// NewMove creates a move object given a pieceId, to and checks whether the move is a capture. If isCapture is set to true
func (board *Board) NewMove(pieceId int, captureId int, to int, promote int) (Move, bool) {
	needsPromotionType := false
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
	// check for promotion
	if board.pieces[pieceId].pieceType == PAWN && promote == 0 {
		_, y := xy(to)
		if (y == 0 && !board.pieces[pieceId].isBlack) || (y == 7 && board.pieces[pieceId].isBlack) {
			needsPromotionType = true
		}
	}
	return Move{pieceId: pieceId, captureId: captureId, from: from, to: to, promote: promote}, needsPromotionType
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
	// promotion
	if m.promote != 0 {
		switch m.promote {
		case -1:
			board.pieces[m.pieceId].pieceType = 'p' // reverse promotion
		case 1:
			board.pieces[m.pieceId].pieceType = 'q'
		case 2:
			board.pieces[m.pieceId].pieceType = 'r'
		case 3:
			board.pieces[m.pieceId].pieceType = 'b'
		case 4:
			board.pieces[m.pieceId].pieceType = 'n'
		}
	}

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
			rookMove, _ = board.NewMove(board.pos2PieceId[m.from+3*EAST], 0, m.from+EAST, 0)
		} else if m.from+2*WEST == m.to {
			isCastle = true
			rookMove, _ = board.NewMove(board.pos2PieceId[m.from+4*WEST], 0, m.from+WEST, 0)
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
	if board.isBlacksTurn {
		board.nextMove++
	}

	if m.captureId != 0 || board.pieces[m.pieceId].pieceType != PAWN {
		board.halfMoves++
	} else {
		board.halfMoves = 0
	}

	board.isBlacksTurn = !board.isBlacksTurn
	board.setMovement()
	return rookMove
}

func (board *Board) reverseMove(m *Move, boardPrimitives *BoardPrimitives) {
	revPromoteInto := 0
	if m.promote != 0 {
		revPromoteInto = -1 // reverse promote back into a pawn
	}
	revMmove, _ := board.NewMove(m.pieceId, 0, m.from, revPromoteInto)
	if board.pieces[m.pieceId].pieceType == KING {
		_, y := xy(m.from)
		// king side
		if m.to-m.from == 2 {
			rookId := board.pos2PieceId[y*8+5]
			reverseCastleMove, _ := board.NewMove(rookId, 0, y*8+7, 0)
			board.TempMove(&reverseCastleMove)
		} else if m.to-m.from == -2 { // queen side
			rookId := board.pos2PieceId[y*8+3]
			reverseCastleMove, _ := board.NewMove(rookId, 0, y*8, 0)
			board.TempMove(&reverseCastleMove)
		}
	}

	board.TempMove(&revMmove)
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
			board.pieces[m.captureId].pos = posOfCapturedPawn
			board.pieces[m.captureId].posB = 1 << posOfCapturedPawn
			board.pos2PieceId[posOfCapturedPawn] = m.captureId
		} else {
			board.pieces[m.captureId].pos = m.to
			board.pieces[m.captureId].posB = 1 << m.to
			board.pos2PieceId[m.to] = m.captureId
		}
		// need to update the combined positions
		board.whitePiecePosB = board.combinePositionsOf(board.whiteIds)
		board.blackPiecePosB = board.combinePositionsOf(board.blackIds)
	}

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

func (board *Board) setPieceCanMoveTo(piece *Piece, pos int) {
	// check whether we are in check
	// if that is the case we can only block, but if we are a king we can move out of check
	if piece.isBlack == board.isBlacksTurn && piece.pieceType != KING {
		if board.check && !board.doubleCheck {
			if board.blockCheckSquaresB&(1<<pos) == 0 {
				return
			}
		} else if board.doubleCheck {
			return
		}
	}
	// check whether the piece can move there or whether it's pinned to some line
	if piece.pinnedMoveB&(1<<pos) != 0 {
		piece.movementB |= 1 << pos
		piece.moves[piece.numMoves] = pos
		piece.numMoves++
	}
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
		file, rank := xy(m.from)
		// ask for rank to avoid weird bug where black promotes to a rook
		if file == 7 {
			if piece.isBlack && rank == 0 {
				board.black_castle_king = false
			} else if !piece.isBlack && rank == 7 {
				board.white_castle_king = false
			}
		} else if file == 0 {
			if piece.isBlack && rank == 0 {
				board.black_castle_queen = false
			} else if !piece.isBlack && rank == 7 {
				board.white_castle_queen = false
			}
		}
	}

	// if rook gets captured
	if m.captureId > 0 {
		capturedPiece := board.pieces[m.captureId]
		file, rank := xy(m.to)
		// ask for rank to avoid weird bug where black promotes to a rook
		if capturedPiece.pieceType == ROOK {
			if capturedPiece.isBlack && rank == 0 {
				if file == 7 {
					board.black_castle_king = false
				} else if file == 0 {
					board.black_castle_queen = false
				}
			} else if !capturedPiece.isBlack && rank == 7 {
				if file == 7 {
					board.white_castle_king = false
				} else if file == 0 {
					board.white_castle_queen = false
				}
			}
		}
	}
}

func (board *Board) countMaterialOfColor(isBlack bool) int {
	pieceIds := board.whiteIds
	if isBlack {
		pieceIds = board.blackIds
	}
	materialCount := 0
	for _, pieceId := range pieceIds {
		if board.pieces[pieceId].posB != 0 {
			materialCount += materialCountMap[board.pieces[pieceId].pieceType]
		}
	}
	return materialCount
}
