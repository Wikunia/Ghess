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
	gameEnded, endType, _ := board.CheckGameEnded()
	if gameEnded {
		if endType == "checkmate" {
			if board.IsBlacksTurn {
				return 100000.0 - float64(board.nextMove)
			} else {
				return -(100000.0 - float64(board.nextMove))
			}
		} else if endType == "draw" {
			return 0.0
		}
	}
	// white pieces - black pieces
	material := float64(board.countMaterialOfColor(false) - board.countMaterialOfColor(true))

	// piece activity
	activity := 0.0
	if board.IsBlacksTurn {
		activity += board.getWhiteMovementScore()
	} else {
		activity += board.getBlackMovementScore()
	}
	board.IsBlacksTurn = !board.IsBlacksTurn
	board.setMovement()
	if board.IsBlacksTurn {
		activity += board.getWhiteMovementScore()
	} else {
		activity += board.getBlackMovementScore()
	}
	board.IsBlacksTurn = !board.IsBlacksTurn
	board.setMovement()

	return material + activity
}

func (board *Board) getWhiteMovementScore() float64 {
	activity := 0.0
	for _, PieceId := range board.whiteIds {
		piece := board.pieces[PieceId]
		if piece.posB == 0 || (piece.pieceType == 'q' && board.nextMove < 10) {
			continue
		}
		pieceAct := 0.0
		for i := 0; i < 32; i++ {
			var bit uint64 = 1 << i
			if piece.movementB&bit != 0 && board.whitePiecePosB&bit == 0 {
				pieceAct += 1.0
			}
		}
		activity += pieceAct
	}

	return activity
}

func (board *Board) getBlackMovementScore() float64 {
	activity := 0.0
	for _, PieceId := range board.blackIds {
		piece := board.pieces[PieceId]
		if piece.posB == 0 || (piece.pieceType == 'q' && board.nextMove < 10) {
			continue
		}
		pieceAct := 0.0
		for i := 32; i < 64; i++ {
			var bit uint64 = 1 << i
			if piece.movementB&bit != 0 && board.blackPiecePosB&bit == 0 {
				pieceAct -= 1.0
			}
		}
		activity += pieceAct
	}
	return activity
}

func (board *Board) getPossibleMovesOrdered(usePv bool, pv [30]Move, currentDepth int) []Move {
	orderedMoves := board.getPossibleMoves()
	rand.Shuffle(len(orderedMoves), func(i, j int) { orderedMoves[i], orderedMoves[j] = orderedMoves[j], orderedMoves[i] })
	if pv[currentDepth].PieceId == 0 || !usePv {
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
		fmt.Println("board.isBlacksTurn: ", board.IsBlacksTurn)
		fmt.Println(orderedMoves)
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
		if move.PieceId == 0 {
			break
		}
		str += GetAlgebraicFromMove(&move) + " "
	}
	fmt.Println(str)
}

func (board *Board) AlphaBetaEngineMove(bestPv [30]Move, currentDepth int, maxDepth int, completedOnce bool, verbose bool, maxDuration int) [30]Move {
	startTime := time.Now()
	myColor := board.IsBlacksTurn
	bestScore := 0.0
	maxTime := time.Duration(maxDuration) * time.Millisecond
	stopPondering := make(chan bool)
	if maxDepth > 30 {
		maxDepth = 30
	}

	for time.Since(startTime) <= maxTime && currentDepth <= maxDepth {
		completed, score, pv := board.alphaBetaPruning(stopPondering, completedOnce, 0, currentDepth, math.Inf(-1), math.Inf(1), !myColor, bestPv, true, startTime, maxTime)
		if completed {
			bestPv = pv
			bestScore = score
			currentDepth += 1
			completedOnce = true
		}
	}

	if verbose {
		fmt.Printf("evaluated up to depth %d in %.02f sec.\n", currentDepth-1, time.Since(startTime).Seconds())
		printPv(bestPv)
		fmt.Println("score from whites perspective: ", bestScore)
	}

	return bestPv
}

func (board *Board) AlphaBetaEnginePonder(stopPondering chan bool, isready chan bool, currentBestPv chan [30]Move) {
	startTime := time.Now()
	myColor := board.IsBlacksTurn
	maxTime := 10000000 * time.Millisecond
	bestPv := [30]Move{}
	currentDepth := 2
	completedOnce := false

	for {
		select {
		case <-stopPondering:
			isready <- true
			return
		default:
			completed, _, pv := board.alphaBetaPruning(stopPondering, completedOnce, 0, currentDepth, math.Inf(-1), math.Inf(1), !myColor, bestPv, true, startTime, maxTime)
			if completed {
				bestPv = pv
				currentDepth += 1
				completedOnce = true
				currentBestPv <- bestPv
			} else {
				isready <- true
				return
			}
		}
	}
}

	sort.Slice(completelyEvaluated, func(i, j int) bool {
		return completelyEvaluated[i].score > completelyEvaluated[j].score
	})

func (board *Board) quiesce(maximizing bool) float64 {
	worstScore := board.staticEvaluation()
	// is there an immediate capture which will reduce my score?
	for _, move := range board.getPossibleCaptures() {
		boardPrimitives := board.getBoardPrimitives()
		board.Move(&move)
		scoreAfterwards := board.staticEvaluation()
		board.reverseMove(&move, &boardPrimitives)
		if !maximizing {
			if scoreAfterwards < worstScore {
				worstScore = scoreAfterwards
			}
		} else {
			if scoreAfterwards > worstScore {
				worstScore = scoreAfterwards
			}
		}
	}
	return worstScore
}

func (board *Board) alphaBetaPruning(stopPondering chan bool, completedOnce bool, currentDepth, depth int, alpha, beta float64, maximizing bool, startPV [30]Move, usePv bool,
	startTime time.Time, maxTime time.Duration) (bool, float64, [30]Move) {

	gameEnded, _, _ := board.CheckGameEnded()
	if depth == 0 || gameEnded {
		return true, board.quiesce(maximizing), startPV
	}
	moves := board.getPossibleMovesOrdered(usePv, startPV, currentDepth)
	bestPv := startPV

	// maximizing player
	if maximizing {
		maxEval := math.Inf(-1)
		for i, move := range moves {
			select {
			case <-stopPondering:
				fmt.Println("stop pondering inside")
				return false, 0.0, startPV
			default:
				if time.Since(startTime) >= maxTime && completedOnce {
					return false, 0.0, startPV
				}

				boardPrimitives := board.getBoardPrimitives()
				board.Move(&move)
				completed, eval, pv := board.alphaBetaPruning(stopPondering, completedOnce, currentDepth+1, depth-1, alpha, beta, !maximizing, startPV, usePv && i == 0, startTime, maxTime)
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
					return true, maxEval, bestPv
				}
			}
		}
		return true, maxEval, bestPv
	} else { // minimizing player
		minEval := math.Inf(1)
		for i, move := range moves {
			select {
			case <-stopPondering:
				fmt.Println("stop pondering inside")
				return false, 0.0, startPV
			default:
				if time.Since(startTime) >= maxTime && completedOnce {
					return false, 0.0, startPV
				}

				boardPrimitives := board.getBoardPrimitives()
				board.Move(&move)
				completed, eval, pv := board.alphaBetaPruning(stopPondering, completedOnce, currentDepth+1, depth-1, alpha, beta, !maximizing, startPV, usePv && i == 0, startTime, maxTime)
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
					return true, minEval, bestPv
				}
			}
		}
		return true, minEval, bestPv
	}
}
