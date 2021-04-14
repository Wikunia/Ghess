package ghess

func (piece *Piece) canMoveTo(pos int) bool {
	var posB uint64 = 1 << pos
	return piece.movementB&posB != 0
}
