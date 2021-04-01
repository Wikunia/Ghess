package ghess

import (
	"fmt"
	"unicode"
)

func isFree(y, x int) bool {
	return currentBoard.position[y-1][x-1] == 0
}

func isLegalPawn(m *JSONMove) bool {
	color := currentBoard.color
	if !color {
		return isLegalPawnWhite(m)
	}
	return isLegalPawnBlack(m)

}

func isLegalPawnWhite(m *JSONMove) bool {
	fromX := currentBoard.pieces[m.PieceId].position.x
	fromY := currentBoard.pieces[m.PieceId].position.y

	// normal move
	if m.CaptureId == 0 {
		if fromX != m.ToX {
			return false
		}
		if fromY-m.ToY == 2 {
			if fromY != 7 {
				return false
			}
			return isFree(fromY-1, fromX)
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

func isLegalPawnBlack(m *JSONMove) bool {
	fromX := currentBoard.pieces[m.PieceId].position.x
	fromY := currentBoard.pieces[m.PieceId].position.y

	// normal move
	if m.CaptureId == 0 {
		if fromX != m.ToX {
			return false
		}
		if m.ToY-fromY == 2 {
			if fromY != 2 {
				return false
			}
			return isFree(fromY+1, fromX)
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

func isLegalKing(m *JSONMove) bool {
	fromX := currentBoard.pieces[m.PieceId].position.x
	fromY := currentBoard.pieces[m.PieceId].position.y
	toX := m.ToX
	toY := m.ToY
	diffx := abs(toX - fromX)
	diffy := abs(toY - fromY)
	if diffx <= 1 && diffy <= 1 {
		return true
	}
	diffx = toX - fromX
	color := getPieceColor(m.PieceId)
	if !color {
		if fromY != 8 {
			return false
		}
		// check king side castle
		if diffx == 2 {
			if !currentBoard.white_castle_king {
				return false
			}
			return isFree(fromY, fromX+1) && isFree(fromY, fromX+2)
		} else if diffx == -2 { // check queen side castle
			if !currentBoard.white_castle_queen {
				return false
			}
			return isFree(fromY, fromX-1) && isFree(fromY, fromX-2) && isFree(fromY, fromX-3)
		}
		return false
	} else {
		if fromY != 1 {
			return false
		}
		// check king side castle
		if diffx == 2 {
			if !currentBoard.black_castle_king {
				return false
			}
			return isFree(fromY, fromX+1) && isFree(fromY, fromX+2)
		} else if diffx == -2 { // check queen side castle
			if !currentBoard.black_castle_queen {
				return false
			}
			return isFree(fromY, fromX-1) && isFree(fromY, fromX-2) && isFree(fromY, fromX-3)
		}
		return false
	}

	return false
}

func isLegalInY(m *JSONMove) bool {
	fromX := currentBoard.pieces[m.PieceId].position.x
	fromY := currentBoard.pieces[m.PieceId].position.y
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
		if !isFree(y, fromX) {
			return false
		}
	}
	return true
}

func isLegalInX(m *JSONMove) bool {
	fromX := currentBoard.pieces[m.PieceId].position.x
	fromY := currentBoard.pieces[m.PieceId].position.y
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
		if !isFree(fromY, x) {
			return false
		}
	}
	return true
}

func isLegalInDiag(m *JSONMove) bool {
	fromX := currentBoard.pieces[m.PieceId].position.x
	fromY := currentBoard.pieces[m.PieceId].position.y
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
		if !isFree(y, x) {
			return false
		}
	}
	return true
}

func isLegalKnight(m *JSONMove) bool {
	fromX := currentBoard.pieces[m.PieceId].position.x
	fromY := currentBoard.pieces[m.PieceId].position.y
	toX := m.ToX
	toY := m.ToY
	diffx := abs(toX - fromX)
	diffy := abs(toY - fromY)
	return (diffx == 1 && diffy == 2) || (diffx == 2 && diffy == 1)
}

func isLegal(m *JSONMove) bool {
	pieceType := unicode.ToLower(rune(currentBoard.pieces[m.PieceId].c))
	// check that the other piece is of the other color
	if m.CaptureId != 0 {
		if currentBoard.pieces[m.CaptureId].color == currentBoard.pieces[m.PieceId].color {
			return false
		}
	}

	fmt.Println(pieceType == 'p')
	switch pieceType {
	case 'p':
		return isLegalPawn(m)
	case 'k':
		return isLegalKing(m)
	case 'q':
		return isLegalInX(m) || isLegalInY(m) || isLegalInDiag(m)
	case 'r':
		return isLegalInX(m) || isLegalInY(m)
	case 'b':
		return isLegalInDiag(m)
	case 'n':
		return isLegalKnight(m)
	}
	return false
}
