package ghess

import (
	"unicode"
)

func (board *Board) GetNumberOfMoves(ply int) int {
	possibleMoves := board.getPossibleMoves()
	if ply == 1 {
		return len(possibleMoves)
	}
	if len(possibleMoves) == 0 {
		return 0
	}

	n := 0
	for _, m := range possibleMoves {
		fromY := board.pieces[m.PieceId].position.y
		fromX := board.pieces[m.PieceId].position.x
		boardPrimitives := board.getBoardPrimitives()
		capturedId, castledMove := board.move(&m)
		n += board.GetNumberOfMoves(ply - 1)
		board.reverseMove(&m, fromY, fromX, capturedId, &castledMove, boardPrimitives)
	}
	return n
}

func (board *Board) getPossibleMoves() []JSONMove {
	var moves []JSONMove
	for _, piece := range board.pieces {
		if !piece.onBoard {
			continue
		}
		pieceType := unicode.ToLower(piece.c)
		// only the current color can move
		if piece.color != board.color {
			continue
		}
		switch pieceType {
		case 'p':
			board.addPawnMoves(&piece, &moves)
		case 'k':
			board.addKingMoves(&piece, &moves)
		case 'q':
			board.addQueenMoves(&piece, &moves)
		case 'r':
			board.addRookMoves(&piece, &moves)
		case 'n':
			board.addKnightMoves(&piece, &moves)
		case 'b':
			board.addBishopMoves(&piece, &moves)
		}
	}
	// fmt.Println(moves)
	return moves
}

func (board *Board) addPawnMoves(piece *Piece, moves *[]JSONMove) {
	if !board.color {
		board.addWhitePawnMoves(piece, moves)
	} else {
		board.addBlackPawnMoves(piece, moves)
	}
}

func (board *Board) addWhitePawnMoves(piece *Piece, moves *[]JSONMove) {
	x := piece.position.x
	y := piece.position.y
	move := JSONMove{PieceId: piece.id, CaptureId: 0, ToX: x, ToY: y}
	move.ToY = y - 2
	if board.isLegal(&move) {
		*moves = append(*moves, move)
	}
	move.ToY = y - 1
	for dx := -1; dx <= 1; dx++ {
		move.ToX = x + dx
		move.CaptureId = 0
		if board.isLegal(&move) {
			*moves = append(*moves, move)
		}
	}
}

func (board *Board) addBlackPawnMoves(piece *Piece, moves *[]JSONMove) {
	x := piece.position.x
	y := piece.position.y
	move := JSONMove{PieceId: piece.id, CaptureId: 0, ToX: x, ToY: y}
	move.ToY = y + 2
	if board.isLegal(&move) {
		*moves = append(*moves, move)
	}
	move.ToY = y + 1
	for dx := -1; dx <= 1; dx++ {
		move.ToX = x + dx
		move.CaptureId = 0
		if board.isLegal(&move) {
			*moves = append(*moves, move)
		}
	}
}

func (board *Board) addKingMoves(piece *Piece, moves *[]JSONMove) {
	x := piece.position.x
	y := piece.position.y
	move := JSONMove{PieceId: piece.id, CaptureId: 0, ToX: x, ToY: y}
	for dx := -1; dx <= 1; dx++ {
		move.ToX = x + dx
		for dy := -1; dy <= 1; dy++ {
			move.ToY = y + dy
			move.CaptureId = 0
			if board.isLegal(&move) {
				*moves = append(*moves, move)
			}
		}
	}
	// check castling
	if (y == 1 && piece.color) || (y == 8 && !piece.color) {
		for dx := -2; dx <= 2; dx += 4 {
			move.ToY = y
			move.ToX = x + dx
			move.CaptureId = 0
			if board.isLegal(&move) {
				*moves = append(*moves, move)
			}
		}
	}
}

func (board *Board) addXMoves(piece *Piece, moves *[]JSONMove) {
	y := piece.position.y
	x := piece.position.x
	move := JSONMove{PieceId: piece.id, CaptureId: 0, ToX: x, ToY: y}

	for f := -1; f <= 1; f += 2 {
		for dx := 1; dx <= 7; dx++ {
			move.ToX = x + f*dx
			move.CaptureId = 0
			if board.isLegal(&move) {
				*moves = append(*moves, move)
			}
			if !board.isFree(move.ToY, move.ToX) {
				break
			}
		}
	}
}

func (board *Board) addYMoves(piece *Piece, moves *[]JSONMove) {
	y := piece.position.y
	x := piece.position.x
	move := JSONMove{PieceId: piece.id, CaptureId: 0, ToX: x, ToY: y}
	for f := -1; f <= 1; f += 2 {
		for dy := 1; dy <= 7; dy++ {
			move.ToY = y + f*dy
			move.CaptureId = 0
			if board.isLegal(&move) {
				*moves = append(*moves, move)
			}
			if !board.isFree(move.ToY, move.ToX) {
				break
			}
		}
	}
}

func (board *Board) addDiagMoves(piece *Piece, moves *[]JSONMove) {
	y := piece.position.y
	x := piece.position.x
	move := JSONMove{PieceId: piece.id, CaptureId: 0, ToX: x, ToY: y}
	for fx := -1; fx <= 1; fx += 2 {
		for fy := -1; fy <= 1; fy += 2 {
			for d := 1; d <= 7; d++ {
				move.ToY = y + fy*d
				move.ToX = x + fx*d
				move.CaptureId = 0
				if board.isLegal(&move) {
					*moves = append(*moves, move)
				}
				if !board.isFree(move.ToY, move.ToX) {
					break
				}
			}
		}
	}
}

func (board *Board) addQueenMoves(piece *Piece, moves *[]JSONMove) {
	board.addXMoves(piece, moves)
	board.addYMoves(piece, moves)
	board.addDiagMoves(piece, moves)
}

func (board *Board) addRookMoves(piece *Piece, moves *[]JSONMove) {
	board.addXMoves(piece, moves)
	board.addYMoves(piece, moves)
}

func (board *Board) addBishopMoves(piece *Piece, moves *[]JSONMove) {
	board.addDiagMoves(piece, moves)
}

func (board *Board) addKnightMoves(piece *Piece, moves *[]JSONMove) {
	y := piece.position.y
	x := piece.position.x
	move := JSONMove{PieceId: piece.id, CaptureId: 0, ToX: x, ToY: y}
	for dx := -2; dx <= 2; dx++ {
		for dy := -2; dy <= 2; dy++ {
			if abs(dx)+abs(dy) != 3 {
				continue
			}
			move.CaptureId = 0
			move.ToX = x + dx
			move.ToY = y + dy
			if board.isLegal(&move) {
				*moves = append(*moves, move)
			}
		}
	}
}
