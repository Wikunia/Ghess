package ghess

import (
	"fmt"
	"unicode"
)

func (board *Board) isFree(y, x int) bool {
	return board.position[y-1][x-1] == 0
}

func (board *Board) isLegalPawn(m *JSONMove) bool {
	color := board.color
	if !color {
		return board.isLegalPawnWhite(m)
	}
	return board.isLegalPawnBlack(m)
}

func (board *Board) isLegalPawnWhite(m *JSONMove) bool {
	fromX := board.pieces[m.PieceId].position.x
	fromY := board.pieces[m.PieceId].position.y

	// normal move
	if m.CaptureId == 0 {
		if fromX != m.ToX {
			// check en passant
			diff := abs(fromX - m.ToX)
			return diff == 1 && fromY-m.ToY == 1 && board.en_passant_position.x == m.ToX && board.en_passant_position.y == m.ToY
		}
		if fromY-m.ToY == 2 {
			if fromY != 7 {
				return false
			}
			return board.isFree(fromY-1, fromX)
		}
		return fromY-m.ToY == 1
	} else { // capture
		diff := abs(fromX - m.ToX)
		if diff == 1 && fromY-m.ToY == 1 {
			return true
		}
	}
	return false
}

func (board *Board) isLegalPawnBlack(m *JSONMove) bool {
	fromX := board.pieces[m.PieceId].position.x
	fromY := board.pieces[m.PieceId].position.y

	// normal move
	if m.CaptureId == 0 {
		if fromX != m.ToX {
			// check en passant
			diff := abs(fromX - m.ToX)
			return diff == 1 && m.ToY-fromY == 1 && board.en_passant_position.x == m.ToX && board.en_passant_position.y == m.ToY
		}
		if m.ToY-fromY == 2 {
			if fromY != 2 {
				return false
			}
			return board.isFree(fromY+1, fromX)
		}
		return m.ToY-fromY == 1
	} else { // capture
		diff := abs(fromX - m.ToX)
		if diff == 1 && m.ToY-fromY == 1 {
			return true
		}
	}
	return false
}

func (board *Board) isLegalKing(m *JSONMove) bool {
	fromX := board.pieces[m.PieceId].position.x
	fromY := board.pieces[m.PieceId].position.y
	toX := m.ToX
	toY := m.ToY
	diffx := abs(toX - fromX)
	diffy := abs(toY - fromY)
	if diffx <= 1 && diffy <= 1 {
		return true
	}
	diffx = toX - fromX
	color := board.getPieceColor(m.PieceId)
	if !color {
		if fromY != 8 {
			return false
		}
		// check king side castle
		if diffx == 2 {
			if !board.white_castle_king {
				return false
			}
			return board.isFree(fromY, fromX+1) && board.isFree(fromY, fromX+2)
		} else if diffx == -2 { // check queen side castle
			if !board.white_castle_queen {
				return false
			}
			return board.isFree(fromY, fromX-1) && board.isFree(fromY, fromX-2) && board.isFree(fromY, fromX-3)
		}
		return false
	} else {
		if fromY != 1 {
			return false
		}
		// check king side castle
		if diffx == 2 {
			if !board.black_castle_king {
				return false
			}
			return board.isFree(fromY, fromX+1) && board.isFree(fromY, fromX+2)
		} else if diffx == -2 { // check queen side castle
			if !board.black_castle_queen {
				return false
			}
			return board.isFree(fromY, fromX-1) && board.isFree(fromY, fromX-2) && board.isFree(fromY, fromX-3)
		}
		return false
	}

	return false
}

func (board *Board) isLegalInY(m *JSONMove) bool {
	fromX := board.pieces[m.PieceId].position.x
	fromY := board.pieces[m.PieceId].position.y
	toX := m.ToX
	toY := m.ToY
	diffx := abs(toX - fromX)
	diffy := abs(toY - fromY)
	if diffx != 0 {
		return false
	}
	var f int
	if f = 1; fromY > toY {
		f = -1
	}
	for dy := 1; dy < diffy; dy++ {
		y := fromY + f*dy
		if !board.isFree(y, fromX) {
			return false
		}
	}
	return true
}

func (board *Board) isLegalInX(m *JSONMove) bool {
	fromX := board.pieces[m.PieceId].position.x
	fromY := board.pieces[m.PieceId].position.y
	toX := m.ToX
	toY := m.ToY
	diffx := abs(toX - fromX)
	diffy := abs(toY - fromY)
	if diffy != 0 {
		return false
	}
	var f int
	if f = 1; fromX > toX {
		f = -1
	}
	for dx := 1; dx < diffx; dx++ {
		x := fromX + f*dx
		if !board.isFree(fromY, x) {
			return false
		}
	}
	return true
}

func (board *Board) isLegalInDiag(m *JSONMove) bool {
	fromX := board.pieces[m.PieceId].position.x
	fromY := board.pieces[m.PieceId].position.y
	toX := m.ToX
	toY := m.ToY
	diffx := abs(toX - fromX)
	diffy := abs(toY - fromY)
	if diffx != diffy {
		return false
	}
	var fx int
	if fx = 1; fromX > toX {
		fx = -1
	}
	var fy int
	if fy = 1; fromY > toY {
		fy = -1
	}
	for d := 1; d < diffx; d++ {
		x := fromX + fx*d
		y := fromY + fy*d
		if !board.isFree(y, x) {
			return false
		}
	}
	return true
}

func (board *Board) isLegalKnight(m *JSONMove) bool {
	fromX := board.pieces[m.PieceId].position.x
	fromY := board.pieces[m.PieceId].position.y
	toX := m.ToX
	toY := m.ToY
	diffx := abs(toX - fromX)
	diffy := abs(toY - fromY)
	return (diffx == 1 && diffy == 2) || (diffx == 2 && diffy == 1)
}

func (board *Board) isLegal(m *JSONMove) bool {
	pieceType := unicode.ToLower(rune(board.pieces[m.PieceId].c))
	// check that the other piece is of the other color
	if m.CaptureId != 0 {
		if board.pieces[m.CaptureId].color == board.pieces[m.PieceId].color {
			return false
		}
	}

	fmt.Println(pieceType == 'p')
	switch pieceType {
	case 'p':
		return board.isLegalPawn(m)
	case 'k':
		return board.isLegalKing(m)
	case 'q':
		return board.isLegalInX(m) || board.isLegalInY(m) || board.isLegalInDiag(m)
	case 'r':
		return board.isLegalInX(m) || board.isLegalInY(m)
	case 'b':
		return board.isLegalInDiag(m)
	case 'n':
		return board.isLegalKnight(m)
	}
	return false
}
