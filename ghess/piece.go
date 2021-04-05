package ghess

func (board *Board) hasBlackCaptureVisionOn(pos Position) bool {
	for _, piece := range board.pieces {
		if !piece.onBoard || !piece.isBlack {
			continue
		}
		if piece.vision[pos.y][pos.x] {
			return piece.pieceType != PAWN || piece.position.x != pos.x
		}
	}
	return false
}

func (board *Board) hasWhiteCaptureVisionOn(pos Position) bool {
	for _, piece := range board.pieces {
		if !piece.onBoard || piece.isBlack {
			continue
		}
		if piece.vision[pos.y][pos.x] {
			return piece.pieceType != PAWN || piece.position.x != pos.x
		}
	}
	return false
}

func (board *Board) updateVision() {
	for pieceId := range board.pieces {
		piece := board.pieces[pieceId]
		if !piece.onBoard {
			continue
		}
		board.updatePieceVision(&board.pieces[pieceId])
	}
}

func (board *Board) updateVisionOnPosition(y, x int) {
	for pieceId, piece := range board.pieces {
		if !piece.onBoard {
			continue
		}
		if piece.vision[y][x] {
			board.updatePieceVision(&board.pieces[pieceId])
		}
	}
}

func (piece *Piece) setVision(y, x int) {
	if isInside(y, x) {
		piece.vision[y][x] = true
	}
}

// resetVision sets the vision to false
func resetVision(piece *Piece) {
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			piece.vision[i][j] = false
		}
	}
}

func (board *Board) resetUpdateMovement() {
	for pieceId := range board.pieces {
		if board.pieces[pieceId].pieceType != KING {
			board.pieces[pieceId].updateMovement = false
		}
	}
}

func (board *Board) setAllUpdateMovement() {
	for pieceId := range board.pieces {
		if board.pieces[pieceId].onBoard {
			board.pieces[pieceId].updateMovement = true
		}
	}
}

func (board *Board) updatePieceVision(piece *Piece) {
	piece.updateMovement = true
	resetVision(piece)
	piece.vision[piece.position.y][piece.position.x] = true
	switch piece.pieceType {
	case KING:
		board.updateKingVision(piece)
	case QUEEN:
		board.updateQueenVision(piece)
	case BISHOP:
		board.updateBishopVision(piece)
	case ROOK:
		board.updateRookVision(piece)
	case KNIGHT:
		board.updateKnightVision(piece)
	case PAWN:
		board.updatePawnVision(piece)
	}
}

func (board *Board) updateKingVision(piece *Piece) {
	cx := piece.position.x
	cy := piece.position.y
	for dy := -1; dy <= 1; dy++ {
		y := cy + dy
		for dx := -1; dx <= 1; dx++ {
			x := cx + dx
			piece.setVision(y, x)
		}
	}

	// castling
	// king side
	if (piece.isBlack && board.black_castle_king) || (!piece.isBlack && board.white_castle_king) {
		if board.isFree(cy, cx+1) {
			piece.setVision(cy, cx+2)
			if board.isFree(cy, cx+2) {
				piece.setVision(cy, cx+3)
			}
		}
	}
	// queen side
	if (piece.isBlack && board.black_castle_queen) || (!piece.isBlack && board.white_castle_queen) {
		if board.isFree(cy, cx-1) {
			piece.setVision(cy, cx-2)
			if board.isFree(cy, cx-2) {
				piece.setVision(cy, cx-3)
				if board.isFree(cy, cx-3) {
					piece.setVision(cy, cx-4)
				}
			}
		}
	}

}

func (board *Board) updateRookVision(piece *Piece) {
	cx := piece.position.x
	cy := piece.position.y
	for f := -1; f <= 1; f += 2 {
		for dx := 1; dx <= 7; dx++ {
			x := cx + f*dx
			piece.setVision(cy, x)
			if !board.isFree(cy, x) {
				break
			}
		}
	}
	// Y direction
	for f := -1; f <= 1; f += 2 {
		for dy := 1; dy <= 7; dy++ {
			y := cy + f*dy
			piece.setVision(y, cx)
			if !board.isFree(y, cx) {
				break
			}
		}
	}
}

func (board *Board) updateBishopVision(piece *Piece) {
	cx := piece.position.x
	cy := piece.position.y
	for fy := -1; fy <= 1; fy += 2 {
		for fx := -1; fx <= 1; fx += 2 {
			for d := 1; d <= 7; d++ {
				x := cx + fx*d
				y := cy + fy*d
				piece.setVision(y, x)
				if !board.isFree(y, x) {
					break
				}
			}
		}
	}
}

func (board *Board) updateQueenVision(piece *Piece) {
	board.updateRookVision(piece)
	board.updateBishopVision(piece)
}

func (board *Board) updateKnightVision(piece *Piece) {
	cx := piece.position.x
	cy := piece.position.y
	for dx := -2; dx <= 2; dx++ {
		for dy := -2; dy <= 2; dy++ {
			if abs(dx)+abs(dy) != 3 {
				continue
			}
			x := cx + dx
			y := cy + dy
			piece.setVision(y, x)
		}
	}
}

func (board *Board) updatePawnVision(piece *Piece) {
	cx := piece.position.x
	cy := piece.position.y
	// two moves from start rank (keep in mind board starts with black row 0)
	if piece.isBlack && cy == 1 && board.isFree(cy+1, cx) {
		piece.setVision(cy+2, cx)
	} else if !piece.isBlack && cy == 6 && board.isFree(cy-1, cx) {
		piece.setVision(cy-2, cx)
	}
	if piece.isBlack {
		piece.setVision(cy+1, cx)
		piece.setVision(cy+1, cx-1)
		piece.setVision(cy+1, cx+1)
	} else {
		piece.setVision(cy-1, cx)
		piece.setVision(cy-1, cx-1)
		piece.setVision(cy-1, cx+1)
	}
}

func (board *Board) updateMovement() {
	for pieceId, piece := range board.pieces {
		if !piece.onBoard || !piece.updateMovement {
			continue
		}
		board.updatePieceMovement(&board.pieces[pieceId])
		// King updates possible movement every time to check for new checks
		if piece.pieceType != KING {
			board.pieces[pieceId].updateMovement = false
		}
	}
}

func (board *Board) setMovement(piece *Piece, y, x int) {
	if !isInside(y, x) {
		return
	}
	legal := true
	ox := piece.position.x
	oy := piece.position.y
	boardPrimitives := board.getBoardPrimitives()
	capturedId, rookMove := board.tempMove(piece, y, x)
	if board.isBlack && piece.isBlack {
		// check if white king can be captured
		blackKing := board.pieces[board.blackKingId]
		if board.hasWhiteCaptureVisionOn(blackKing.position) {
			legal = false
		}
		// castle
		if rookMove.pieceId != 0 {
			// castle can the rook after castling be captured?
			rook := board.pieces[rookMove.pieceId]
			if board.hasWhiteCaptureVisionOn(rook.position) {
				legal = false
			}
			// can the original king be captured? No castle when in check...
			if board.hasWhiteCaptureVisionOn(Position{y: oy, x: ox}) {
				legal = false
			}
		}
	} else if !board.isBlack && !piece.isBlack {
		// check if white king can be captured
		whiteKing := board.pieces[board.whiteKingId]
		if board.hasBlackCaptureVisionOn(whiteKing.position) {
			legal = false
		}
		// castle
		if rookMove.pieceId != 0 {
			// castle can the rook after castling be captured?
			rook := board.pieces[rookMove.pieceId]
			if board.hasBlackCaptureVisionOn(rook.position) {
				legal = false
			}
			// can the original king be captured? No castle when in check...
			if board.hasBlackCaptureVisionOn(Position{y: oy, x: ox}) {
				legal = false
			}
		}
	}
	board.reverseTempMove(piece, oy, ox, capturedId, rookMove.pieceId, &boardPrimitives)
	board.resetUpdateMovement()

	if legal {
		piece.movement[y][x] = true
		piece.moves[piece.numMoves] = Position{y: y, x: x}
		piece.numMoves += 1
	}
}

func (board *Board) move(piece *Piece, y, x int) (int, Move) {
	capturedId, rookMove := board.tempMove(piece, y, x)

	// remove castle privileges
	// if king or rook moved or rook captured
	board.updateCastlePrivileges()

	board.isBlack = !board.isBlack
	if board.lastMoveWasCheck {
		board.setAllUpdateMovement()
	}

	board.lastMoveWasCheck = false
	// update the vision/movement of new attacked pieces as well
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			if piece.vision[i][j] && board.position[i][j] != 0 {
				attackedPiece := board.pieces[board.position[i][j]]
				if attackedPiece.isBlack != piece.isBlack {
					board.updatePieceVision(&board.pieces[board.position[i][j]])
					// the enemy king is in check => update movement of all pieces on the board
					if board.pieces[board.position[i][j]].pieceType == KING {
						board.setAllUpdateMovement()
						board.lastMoveWasCheck = true
					}
				}
			}
		}
	}

	board.updateMovement()
	return capturedId, rookMove
}

func (board *Board) tempMove(piece *Piece, y, x int) (int, Move) {
	ox := piece.position.x
	oy := piece.position.y

	// check if capture
	isCapture := false
	capturedId := 0
	captureY := -1
	captureX := -1
	// normal capture
	if board.position[y][x] != 0 {
		isCapture = true
		captureY = y
		captureX = x
		// en passant capture
	} else if piece.pieceType == PAWN && board.en_passant_position.x == x && board.en_passant_position.y == y {
		isCapture = true
		captureX = x
		if piece.isBlack {
			captureY = y + 1
		} else {
			captureY = y - 1
		}
	}
	if isCapture {
		capturedId = board.position[y][x]
		board.pieces[capturedId].onBoard = false
	}

	// castle
	rookMove := Move{}
	if piece.pieceType == KING && abs(x-ox) == 2 && ox == 4 {
		if x == 6 {
			rookId := board.position[y][7]
			rook := board.pieces[rookId]
			rookMove = board.newMove(rookId, 0, y, 5, false)
			board.tempMove(&rook, y, 5)
			board.pieces[rookId] = rook
		} else {
			rookId := board.position[y][0]
			rook := board.pieces[rookId]
			rookMove = board.newMove(rookId, 0, y, 3, false)
			board.tempMove(&rook, y, 3)
			board.pieces[rookId] = rook
		}
	}

	// remove the figure from its current position
	board.position[piece.position.y][piece.position.x] = 0
	// add the figure to its new position
	board.position[y][x] = piece.id
	piece.position.x = x
	piece.position.y = y
	// fmt.Println("captureY, captureX: ", captureY, captureX)

	// update vision
	board.updateVisionOnPosition(oy, ox)
	board.updateVisionOnPosition(y, x)
	if isCapture && (captureX != x || captureY != y) {
		board.updateVisionOnPosition(captureY, captureX)
	}

	// en passant possible
	if piece.pieceType == PAWN {
		if (piece.isBlack && y-oy == 2) || (!piece.isBlack && y-oy == -2) {
			board.en_passant_position.y = (y + oy) / 2
			board.en_passant_position.x = x
			board.updateVisionOnPosition((y+oy)/2, x)
		} else {
			board.en_passant_position.y = 0
			board.en_passant_position.x = 0
		}
	} else {
		board.en_passant_position.y = 0
		board.en_passant_position.x = 0
	}

	return capturedId, rookMove
}

func (board *Board) reverseMove(piece *Piece, y, x int, capturedId int, rookId int, boardPrimitives *BoardPrimitives) {
	board.isBlack = !board.isBlack
	board.reverseTempMove(piece, y, x, capturedId, rookId, boardPrimitives)

	board.updateVision()

	board.updateCastlePrivileges()
	board.setAllUpdateMovement()

	board.updateMovement()
}

func (board *Board) reverseTempMove(piece *Piece, y, x int, capturedId int, rookId int, boardPrimitives *BoardPrimitives) {
	board.currentlyReverseMode = true
	// reverse castle rook move
	if rookId != 0 {
		rook := board.pieces[rookId]
		if rook.position.x == 3 {
			board.tempMove(&rook, y, 0)
		} else {
			board.tempMove(&rook, y, 7)
		}
		board.pieces[rookId] = rook
	}
	// move the piece itself back
	board.tempMove(piece, y, x)

	if capturedId != 0 {
		capturedPiece := board.pieces[capturedId]
		board.position[capturedPiece.position.y][capturedPiece.position.x] = capturedId
		board.pieces[capturedId].onBoard = true
		board.updateVisionOnPosition(capturedPiece.position.y, capturedPiece.position.x)
	}
	board.currentlyReverseMode = false
	board.setBoardPrimitives(boardPrimitives)
}

// resetMovement sets the movement to false
func resetMovement(piece *Piece) {
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			piece.movement[i][j] = false
		}
	}
	piece.numMoves = 0
}

func (board *Board) updatePieceMovement(piece *Piece) {
	resetMovement(piece)
	// handle pawn extra as it can only move diagonal if capture incl en passant
	// can only move forward if no capture
	if piece.pieceType == PAWN {
		cx := piece.position.x
		for i := 0; i < 8; i++ {
			for j := 0; j < 8; j++ {
				if piece.vision[i][j] {
					diffx := abs(cx - j)
					// normal move
					if diffx == 0 {
						if board.isFree(i, j) {
							board.setMovement(piece, i, j)
						}
					} else { // capture
						if !board.isFree(i, j) {
							pieceAtPosition := board.pieces[board.position[i][j]]
							if pieceAtPosition.isBlack != piece.isBlack {
								board.setMovement(piece, i, j)
							}
						} else if board.en_passant_position.x == j && board.en_passant_position.y == i && piece.isBlack == board.isBlack {
							board.setMovement(piece, i, j)
						}
					}
				}
			}
		}
		return
	}

	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			if piece.vision[i][j] {
				if board.isFree(i, j) {
					if piece.pieceType == KING {
						// even for queen side castle we can only jump 2 fields the other is just for vision
						if abs(piece.position.x-j) <= 2 {
							board.setMovement(piece, i, j)
						}
					} else {
						board.setMovement(piece, i, j)
					}
				} else {
					pieceAtPosition := board.pieces[board.position[i][j]]
					if pieceAtPosition.isBlack != piece.isBlack {
						board.setMovement(piece, i, j)
					}
				}
			}
		}
	}
}
