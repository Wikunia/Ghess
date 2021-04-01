package ghess

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func isKing(piece Piece) bool {
	return piece.c == 'K' || piece.c == 'k'
}

func isQueen(piece Piece) bool {
	return piece.c == 'Q' || piece.c == 'q'
}

func isRook(piece Piece) bool {
	return piece.c == 'R' || piece.c == 'r'
}

func isBishop(piece Piece) bool {
	return piece.c == 'B' || piece.c == 'b'
}

func isKnight(piece Piece) bool {
	return piece.c == 'N' || piece.c == 'n'
}

func isPawn(piece Piece) bool {
	return piece.c == 'P' || piece.c == 'p'
}
