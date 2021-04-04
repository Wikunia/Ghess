package ghess

import "fmt"

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

func (board *Board) getBoardPrimitives() BoardPrimitives {
	return BoardPrimitives{
		color:               board.color,
		white_castle_king:   board.white_castle_king,
		white_castle_queen:  board.white_castle_queen,
		black_castle_king:   board.black_castle_king,
		black_castle_queen:  board.black_castle_queen,
		en_passant_position: Position{x: board.en_passant_position.x, y: board.en_passant_position.y},
		halfMoves:           board.halfMoves,
		nextMove:            board.nextMove,
		whiteKingId:         board.whiteKingId,
		blackKingId:         board.blackKingId,
	}
}

func (board *Board) setBoardPrimitives(bp BoardPrimitives) {
	board.color = bp.color
	board.white_castle_king = bp.white_castle_king
	board.white_castle_queen = bp.white_castle_queen
	board.black_castle_king = bp.black_castle_king
	board.black_castle_queen = bp.black_castle_queen
	board.en_passant_position.x = bp.en_passant_position.x
	board.en_passant_position.y = bp.en_passant_position.y
	board.halfMoves = bp.halfMoves
	board.nextMove = bp.nextMove
	board.whiteKingId = bp.whiteKingId
	board.blackKingId = bp.blackKingId
}

func (board *Board) isEqual(otherBoard *Board) bool {
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			if board.position[i][j] != otherBoard.position[i][j] {
				fmt.Println("Not equal at x,y: ", i+1, j+1)
				fmt.Println("Now: ", board.position[i][j])
				fmt.Println("Other: ", otherBoard.position[i][j])
				return false
			}
		}
	}
	for key := range board.pieces {
		if board.pieces[key].id != otherBoard.pieces[key].id || board.pieces[key].onBoard != otherBoard.pieces[key].onBoard {
			fmt.Println("Not equal for piece: ", key)
			return false
		}
	}
	return true
}
