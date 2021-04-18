package ghess

import (
	"fmt"
	"os"
	"strconv"
	"time"
	"unicode"
)

func GetAlgebraicFromMove(m *Move) string {
	fromX, fromY := xy(m.from)
	toX, toY := xy(m.to)
	moveStr := string(rune('a'+fromX)) + strconv.Itoa(8-fromY)
	moveStr += string(rune('a'+toX)) + strconv.Itoa(8-toY)
	if m.promote != 0 {
		switch m.promote {
		case 1:
			moveStr += "q"
		case 2:
			moveStr += "r"
		case 3:
			moveStr += "b"
		case 4:
			moveStr += "n"
		}
	}
	return moveStr
}

func (board *Board) MoveLongAlgebraic(moveStr string) error {
	move, err := board.GetMoveFromLongAlgebraic(moveStr)
	if err != nil {
		return err
	}
	board.Move(&move)
	return nil
}

func (board *Board) GetMoveFromLongAlgebraic(moveStr string) (Move, error) {
	move := Move{}
	if len(moveStr) != 4 && len(moveStr) != 5 {
		return move, fmt.Errorf("currently only algebraic notation with 4 or 5 chars (with promotion) is supported")
	}
	fromX := int(moveStr[0] - 'a')
	fromY := 8 - int(moveStr[1]-'0')
	toX := int(moveStr[2] - 'a')
	toY := 8 - int(moveStr[3]-'0')
	pieceId := board.pos2PieceId[fromY*8+fromX]
	if pieceId == 0 {
		return move, fmt.Errorf("there is no piece at that position")
	}
	if board.pieces[pieceId].isBlack != board.isBlacksTurn {
		return move, fmt.Errorf("the piece has the wrong color")
	}
	promotionIdx := 0
	if len(moveStr) == 5 {
		switch moveStr[4] {
		case 'q':
			promotionIdx = 1
		case 'r':
			promotionIdx = 2
		case 'b':
			promotionIdx = 3
		case 'n':
			promotionIdx = 4
		default:
			return Move{}, fmt.Errorf("last char must be q,r,b,n for a 5 character string")
		}
	}
	move, _ = board.NewMove(pieceId, 0, toY*8+toX, promotionIdx)
	if board.isLegal(&move) {
		// capture will be filled automatically
		return move, nil
	}
	return move, fmt.Errorf("the move is not legal")
}

func getPromotionStr(promote int) string {
	switch promote {
	default:
		return ""
	case 1:
		return "=Q"
	case 2:
		return "=R"
	case 3:
		return "=B"
	case 4:
		return "=N"
	}
}

func (board *Board) getBasicStandardAlgebraic(m *Move) string {
	piece := board.pieces[m.pieceId]
	fromX, _ := xy(m.from)
	toX, toY := xy(m.to)
	endToStr := string(rune('a'+toX)) + strconv.Itoa(8-toY)
	// no capture:
	if m.captureId == 0 {
		if piece.pieceType == PAWN {
			return endToStr + getPromotionStr(m.promote)
		} else {
			return string(unicode.ToUpper(piece.pieceType)) + endToStr
		}
	} else { // capture
		if piece.pieceType == PAWN {
			return string(rune('a'+fromX)) + "x" + endToStr + getPromotionStr(m.promote)
		} else {
			return string(unicode.ToUpper(piece.pieceType)) + "x" + endToStr
		}
	}
}

func (board *Board) getDifferentFileStandardAlgebraic(m *Move) string {
	piece := board.pieces[m.pieceId]
	fromX, _ := xy(m.from)
	toX, toY := xy(m.to)
	endToStr := string(rune('a'+toX)) + strconv.Itoa(8-toY)
	if m.captureId == 0 {
		// can't happen for a pawn
		return string(unicode.ToUpper(piece.pieceType)) + string(rune('a'+fromX)) + endToStr
	} else {
		if piece.pieceType == PAWN {
			return string(rune('a'+fromX)) + "x" + endToStr + getPromotionStr(m.promote)
		} else {
			return string(unicode.ToUpper(piece.pieceType)) + string(rune('a'+fromX)) + "x" + endToStr
		}
	}
}

func (board *Board) getDifferentRankStandardAlgebraic(m *Move) string {
	piece := board.pieces[m.pieceId]
	_, fromY := xy(m.from)
	toX, toY := xy(m.to)
	endToStr := string(rune('a'+toX)) + strconv.Itoa(8-toY)
	// can't be a pawn
	if m.captureId == 0 {
		return string(unicode.ToUpper(piece.pieceType)) + strconv.Itoa(8-fromY) + endToStr
	} else {
		return string(unicode.ToUpper(piece.pieceType)) + strconv.Itoa(8-fromY) + "x" + endToStr
	}
}

func (board *Board) getStandardAlgebraicFromMove(m *Move) string {
	ambiguityMoves := []Move{}
	piece := board.pieces[m.pieceId]
	for _, move := range board.getPossibleMoves() {
		if move.to == m.to && piece.pieceType == board.pieces[move.pieceId].pieceType && m.pieceId != move.pieceId {
			ambiguityMoves = append(ambiguityMoves, move)
		}
	}
	fromX, fromY := xy(m.from)
	toX, toY := xy(m.to)
	endToStr := string(rune('a'+toX)) + strconv.Itoa(8-toY)

	// castle
	if piece.pieceType == KING && abs(fromX-toX) == 2 {
		if toX == 6 {
			return "O-O"
		} else {
			return "O-O-O"
		}
	}

	// no ambiguity
	if len(ambiguityMoves) == 0 {
		return board.getBasicStandardAlgebraic(m)
	}
	// single ambiguity
	if len(ambiguityMoves) == 1 {
		ambiguityMove := ambiguityMoves[0]
		ambFromX, ambFromY := xy(ambiguityMove.from)
		// different file
		if ambFromX != fromX {
			return board.getDifferentFileStandardAlgebraic(m)
		} else if ambFromY != fromY {
			return board.getDifferentRankStandardAlgebraic(m)
		}
	}
	// several ambiguities
	// check if all ambiguity moves have the same file
	haveSameFile := 0
	haveSameRank := 0
	for _, move := range ambiguityMoves {
		ambFromX, ambFromY := xy(move.from)
		if ambFromX == fromX {
			haveSameFile++
		}
		if ambFromY == fromY {
			haveSameRank++
		}
	}
	if haveSameFile == 0 {
		return board.getDifferentFileStandardAlgebraic(m)
	}
	if haveSameRank == 0 {
		return board.getDifferentRankStandardAlgebraic(m)
	}
	// rare case where we use the full length form
	// can't be a pawn
	if m.captureId == 0 {
		return string(unicode.ToUpper(piece.pieceType)) + string(rune('a'+fromX)) + strconv.Itoa(8-fromY) + endToStr
	} else {
		return string(unicode.ToUpper(piece.pieceType)) + string(rune('a'+fromX)) + strconv.Itoa(8-fromY) + "x" + endToStr
	}
}

// writePGNFile writes a file called YYYY-MM-DD-HH:MM.pgn from by converting playedMoves into the standard algebraic notation
func writePGNFile(playedMoves []Move) {
	dt := time.Now()
	dateTime := dt.Format("2006-01-30_15:04")
	f, _ := os.Create(dateTime + ".pgn")

	defer f.Close()

	board := GetBoardFromFen(START_FEN)
	moveNo := 0
	line := ""
	for nextMoveId := 0; nextMoveId < len(playedMoves); nextMoveId++ {
		if nextMoveId%2 == 0 {
			if moveNo != 0 {
				f.WriteString(line + "\n")
			}
			moveNo++
			line = strconv.Itoa(moveNo) + "."
		}
		line += " " + board.getStandardAlgebraicFromMove(&playedMoves[nextMoveId])
		board.Move(&playedMoves[nextMoveId])
	}
	if line[len(line)-1] != '.' {
		fmt.Println("line: ", line)
		f.WriteString(line + "\n")
	}
}
