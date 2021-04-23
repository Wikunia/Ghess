package ghess

import "fmt"

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func min(x, y int) int {
	if x <= y {
		return x
	}
	return y
}

func xy(n int) (x, y int) {
	y = n / 8
	x = n % 8
	return
}

func bits2array(bits uint64) [8][8]bool {
	var res [8][8]bool
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			pos := i*8 + j
			if bits&(1<<pos) != 0 {
				res[i][j] = true
			}
		}
	}
	return res
}

func (board *Board) getBoardPrimitives() BoardPrimitives {
	return BoardPrimitives{
		IsBlacksTurn:       board.IsBlacksTurn,
		white_castle_king:  board.white_castle_king,
		white_castle_queen: board.white_castle_queen,
		black_castle_king:  board.black_castle_king,
		black_castle_queen: board.black_castle_queen,
		en_passant_pos:     board.en_passant_pos,
		halfMoves:          board.halfMoves,
		nextMove:           board.nextMove,
		whiteKingId:        board.whiteKingId,
		blackKingId:        board.blackKingId,
		fens:               board.fens[:],
	}
}

func (board *Board) setBoardPrimitives(bp *BoardPrimitives) {
	board.IsBlacksTurn = bp.IsBlacksTurn
	board.white_castle_king = bp.white_castle_king
	board.white_castle_queen = bp.white_castle_queen
	board.black_castle_king = bp.black_castle_king
	board.black_castle_queen = bp.black_castle_queen
	board.en_passant_pos = bp.en_passant_pos
	board.halfMoves = bp.halfMoves
	board.nextMove = bp.nextMove
	board.whiteKingId = bp.whiteKingId
	board.blackKingId = bp.blackKingId
	board.fens = bp.fens[:]
}

func printBits(bits uint64) {
	fmt.Println("==================")
	arr := bits2array(bits)
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			if arr[i][j] {
				fmt.Print("X")
			} else {
				fmt.Print("-")
			}
		}
		fmt.Print("\n")
	}
	fmt.Println("==================")
}

func makeRange(min, max int) []int {
	a := make([]int, max-min+1)
	for i := range a {
		a[i] = min + i
	}
	return a
}

func (move *Move) isEqual(m *Move) bool {
	return move.PieceId == m.PieceId && move.from == m.from && move.captureId == m.captureId && move.to == m.to && move.promote == m.promote
}
