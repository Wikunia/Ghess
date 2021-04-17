package ghess

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

type Eval struct {
	id    int
	move  Move
	score float64
}

func (board *Board) staticEvaluation() float64 {
	gameEnded, endType, _ := board.checkGameEnded()
	if gameEnded {
		if endType == "checkmate" {
			if board.isBlacksTurn {
				return 100000.0 - float64(board.nextMove)
			} else {
				return -(100000.0 - float64(board.nextMove))
			}
		} else if endType == "draw" {
			return 0.0
		}
	}
	// white pieces - black pieces
	return float64(board.countMaterialOfColor(false) - board.countMaterialOfColor(true))
}

func (board *Board) getPossibleMovesOrdered(usePv bool, pv [30]Move, currentDepth int) []Move {
	orderedMoves := board.getPossibleMoves()
	rand.Shuffle(len(orderedMoves), func(i, j int) { orderedMoves[i], orderedMoves[j] = orderedMoves[j], orderedMoves[i] })
	if pv[currentDepth].pieceId == 0 || !usePv {
		return orderedMoves
	}
	id := 0
	found := false
	for i, move := range orderedMoves {
		if move.isEqual(&pv[currentDepth]) {
			found = true
			id = i
			break
		}
	}
	if !found {
		fmt.Println("Could not find: ", pv[currentDepth])
	}
	orderedMoves[0], orderedMoves[id] = orderedMoves[id], orderedMoves[0]
	return orderedMoves
}

func (board *Board) getPossibleCaptures() []Move {
	moves := board.getPossibleMoves()
	capturedMoves := []Move{}
	for _, move := range moves {
		if move.captureId != 0 {
			capturedMoves = append(capturedMoves, move)
		}
	}
	return capturedMoves
}

func (board *Board) countNumberOfCaptureMoves() int {
	moves := board.getPossibleMoves()
	numCaptures := 0
	for _, move := range moves {
		if move.captureId != 0 {
			numCaptures++
		}
	}
	return numCaptures
}

func printPv(pv [30]Move) {
	str := ""
	for _, move := range pv {
		if move.pieceId == 0 {
			break
		}
		str += GetAlgebraicFromMove(&move) + " "
	}
	fmt.Println(str)
}

func (board *Board) AlphaBetaEngineMove() Move {
	currentDepth := 2
	startTime := time.Now()
	myColor := board.isBlacksTurn
	bestPv := [30]Move{}
	bestScore := 0.0
	maxTime := MAX_ENGINE_TIME * time.Millisecond
	completedOnce := false

	for time.Since(startTime) <= maxTime && currentDepth < 30 {
		completed, score, pv := board.alphaBetaPruning(completedOnce, 0, currentDepth, math.Inf(-1), math.Inf(1), !myColor, bestPv, true, startTime, maxTime)
		if completed {
			bestPv = pv
			bestScore = score
			currentDepth += 1
			completedOnce = true
		}
	}

	fmt.Printf("evaluated up to depth %d in %.02f sec.\n", currentDepth-1, time.Since(startTime).Seconds())
	printPv(bestPv)
	fmt.Println("score from whites perspective: ", bestScore)
	return bestPv[0]
}

func (board *Board) alphaBetaPruning(completedOnce bool, currentDepth, depth int, alpha, beta float64, maximizing bool, startPV [30]Move, usePv bool,
	startTime time.Time, maxTime time.Duration) (bool, float64, [30]Move) {

	gameEnded, _, _ := board.checkGameEnded()
	if depth == 0 || gameEnded {
		return true, board.staticEvaluation(), startPV
	}
	moves := board.getPossibleMovesOrdered(usePv, startPV, currentDepth)
	bestPv := startPV

	// maximizing player
	if maximizing {
		maxEval := math.Inf(-1)
		for i, move := range moves {
			if time.Since(startTime) >= maxTime && completedOnce {
				return false, 0.0, startPV
			}

			boardPrimitives := board.getBoardPrimitives()
			board.Move(&move)
			completed, eval, pv := board.alphaBetaPruning(completedOnce, currentDepth+1, depth-1, alpha, beta, !maximizing, startPV, usePv && i == 0, startTime, maxTime)
			board.reverseMove(&move, &boardPrimitives)
			if !completed {
				return false, 0.0, startPV
			}
			if eval > maxEval {
				maxEval = eval
				bestPv = pv
				bestPv[currentDepth] = move
			}
			alpha = math.Max(eval, alpha)
			if beta <= alpha {
				break
			}
		}
		return true, maxEval, bestPv
	} else { // minimizing player
		minEval := math.Inf(1)
		for i, move := range moves {
			if time.Since(startTime) >= maxTime && completedOnce {
				return false, 0.0, startPV
			}

			boardPrimitives := board.getBoardPrimitives()
			board.Move(&move)
			completed, eval, pv := board.alphaBetaPruning(completedOnce, currentDepth+1, depth-1, alpha, beta, !maximizing, startPV, usePv && i == 0, startTime, maxTime)
			board.reverseMove(&move, &boardPrimitives)
			if !completed {
				return false, 0.0, startPV
			}
			if eval < minEval {
				minEval = eval
				bestPv = pv
				bestPv[currentDepth] = move
			}
			beta = math.Min(eval, beta)
			if beta <= alpha {
				break
			}
		}
		return true, minEval, bestPv
	}
}
