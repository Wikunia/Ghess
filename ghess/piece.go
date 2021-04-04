package ghess

func (board *Board) updateVision() {
	for pieceId, piece := range board.pieces {
		board.updatePieceVision(&piece)
		board.pieces[pieceId] = piece
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

func (board *Board) updatePieceVision(piece *Piece) {
	resetVision(piece)
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
	if (piece.isBlack && board.black_castle_king) || (!piece.isBlack && board.white_castle_king) {
		piece.setVision(cy, cx+1)
		piece.setVision(cy, cx+2)
	}
	if (piece.isBlack && board.black_castle_queen) || (!piece.isBlack && board.white_castle_queen) {
		piece.setVision(cy, cx-1)
		piece.setVision(cy, cx-2)
		piece.setVision(cy, cx-3)
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
