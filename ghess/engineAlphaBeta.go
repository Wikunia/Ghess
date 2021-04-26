package ghess

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
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

type OrderedMoves struct {
	move  Move
	score int
}

func (board *Board) getPieceVal(pieceId int) int {
	piece := board.pieces[pieceId]
	switch piece.pieceType {
	case 'p':
		return 10
	case 'n':
		return 30
	case 'b':
		return 35
	case 'r':
		return 50
	case 'q':
		return 100
	case 'k':
		return 100000
	}
	return 0
}

func (board *Board) getPossibleMovesOrdered(usePv bool, pv [30]Move, currentDepth int) []OrderedMoves {
	moves := board.getPossibleMoves()
	rand.Shuffle(len(moves), func(i, j int) { moves[i], moves[j] = moves[j], moves[i] })
	// order by capture value
	orderedMoves := make([]OrderedMoves, len(moves))
	for i, move := range moves {
		if move.captureId == 0 {
			orderedMoves[i] = OrderedMoves{move: move, score: 0}
		} else {
			score := board.getPieceVal(move.captureId) - board.getPieceVal(move.PieceId)
			orderedMoves[i] = OrderedMoves{move: move, score: score}
		}
	}
	sort.Slice(orderedMoves, func(i, j int) bool {
		return orderedMoves[i].score > orderedMoves[j].score
	})

	if pv[currentDepth].PieceId == 0 || !usePv {
		return orderedMoves
	}
	id := 0
	found := false
	for i, om := range orderedMoves {
		if om.move.isEqual(&pv[currentDepth]) {
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

func (board *Board) getPossibleCaptures() []OrderedMoves {
	moves := board.getPossibleMoves()
	numCaptureMoves := 0
	for _, move := range moves {
		if move.captureId != 0 {
			numCaptureMoves++
		}
	}

	// order by capture value
	orderedMoves := make([]OrderedMoves, numCaptureMoves)
	c := 0
	for _, move := range moves {
		if move.captureId != 0 {
			score := board.getPieceVal(move.captureId) - board.getPieceVal(move.PieceId)
			orderedMoves[c] = OrderedMoves{move: move, score: score}
			c++
		}
	}
	sort.Slice(orderedMoves, func(i, j int) bool {
		return orderedMoves[i].score > orderedMoves[j].score
	})

	return orderedMoves
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

type AlphaBetaOutput struct {
	Completed     bool
	Score         float64
	Pv            [30]Move
	NodesSearched int
	Depth         int
}

func (board *Board) AlphaBetaEngineMove(bestPv [30]Move, currentDepth int, maxDepth int, completedOnce bool, verbose bool, maxDuration int) AlphaBetaOutput {
	startTime := time.Now()
	myColor := board.IsBlacksTurn
	bestScore := 0.0
	maxTime := time.Duration(maxDuration) * time.Millisecond
	stopPondering := make(chan bool)
	if maxDepth > 30 {
		maxDepth = 30
	}
	numMoves := board.GetNumberOfMoves(1)
	if numMoves == 1 {
		fmt.Println("Only one move possible")
		moves := board.getPossibleMoves()
		bestPv[0] = moves[0]
		return AlphaBetaOutput{Score: math.NaN(), Pv: bestPv}
	}
	factor := time.Duration(1.0)
	lastRun := time.Duration(0.0)
	completeAb := AlphaBetaOutput{Score: math.NaN(), Pv: bestPv}

	for time.Since(startTime) <= maxTime && currentDepth <= maxDepth {
		startRun := time.Now()
		ab := board.alphaBetaPruning(stopPondering, completedOnce, 0, currentDepth, math.Inf(-1), math.Inf(1), !myColor, bestPv, true, startTime, maxTime, AlphaBetaOutput{})
		if ab.Completed {
			bestScore = ab.Score
			currentDepth += 1
			completedOnce = true
			completeAb.Completed = true
			completeAb.Score = bestScore
			completeAb.Depth = ab.Depth
			completeAb.Pv = ab.Pv
			completeAb.NodesSearched += ab.NodesSearched
		}
		if lastRun.Milliseconds() > 1 {
			factor = time.Since(startRun) / lastRun
		}
		lastRun = time.Since(startRun)
		if time.Since(startTime)+factor*lastRun >= maxTime {
			fmt.Println("factor break factor: ", factor)
			break
		}

		if board.IsBlacksTurn && ab.Score < -10000 {
			break
		} else if !board.IsBlacksTurn && ab.Score > 10000 {
			break
		}
	}

	if verbose {
		fmt.Printf("evaluated up to depth %d in %.02f sec.\n", currentDepth-1, time.Since(startTime).Seconds())
		printPv(bestPv)
		fmt.Println("score from whites perspective: ", bestScore)
	}

	return completeAb
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
			ab := board.alphaBetaPruning(stopPondering, completedOnce, 0, currentDepth, math.Inf(-1), math.Inf(1), !myColor, bestPv, true, startTime, maxTime, AlphaBetaOutput{})
			if ab.Completed {
				bestPv = ab.Pv
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

func (board *Board) quiesce(alpha, beta float64, maximizing bool) float64 {
	standPat := board.staticEvaluation()
	if maximizing {
		// it is whites turn now
		if standPat >= beta {
			return beta
		}
		if standPat > alpha {
			alpha = standPat
		}

		for _, om := range board.getPossibleCaptures() {
			boardPrimitives := board.getBoardPrimitives()
			board.Move(&om.move)
			score := board.quiesce(alpha, beta, !maximizing)
			board.reverseMove(&om.move, &boardPrimitives)

			if score >= beta {
				return beta
			}
			if score > alpha {
				alpha = score
			}
		}
		return alpha
	} else {
		// it is blacks turn now
		if standPat <= alpha {
			return alpha
		}
		if standPat < beta {
			beta = standPat
		}

		for _, om := range board.getPossibleCaptures() {
			boardPrimitives := board.getBoardPrimitives()
			board.Move(&om.move)
			score := board.quiesce(alpha, beta, !maximizing)
			board.reverseMove(&om.move, &boardPrimitives)

			if score <= alpha {
				return alpha
			}
			if score < beta {
				beta = score
			}
		}
		return beta
	}

}

func (board *Board) alphaBetaPruning(stopPondering chan bool, completedOnce bool, currentDepth, depth int, alpha, beta float64, maximizing bool, startPV [30]Move, usePv bool,
	startTime time.Time, maxTime time.Duration, output AlphaBetaOutput) AlphaBetaOutput {

	orderedMoves := board.getPossibleMovesOrdered(usePv, startPV, currentDepth)
	gameEnded, _, _ := board.CheckGameEnded()
	if gameEnded || depth == 0 {
		output.Completed = true
		if len(orderedMoves) > 0 {
			startPV[0] = orderedMoves[0].move
		}
	}
	if gameEnded {
		output.Score = board.staticEvaluation()
		return output
	}
	if depth == 0 {
		output.Score = board.quiesce(alpha, beta, maximizing)
		return output
	}
	bestPv := startPV
	notCompletedOutput := AlphaBetaOutput{Completed: false, Score: math.NaN(), Pv: startPV}
	completedOutput := AlphaBetaOutput{Completed: true}

	// maximizing player
	if maximizing {
		maxEval := math.Inf(-1)
		for i, om := range orderedMoves {
			select {
			case <-stopPondering:
				fmt.Println("stop pondering inside")
				return notCompletedOutput
			default:
				move := om.move
				if time.Since(startTime) >= maxTime && completedOnce {
					return notCompletedOutput
				}

				boardPrimitives := board.getBoardPrimitives()
				board.Move(&move)
				ab := board.alphaBetaPruning(stopPondering, completedOnce, currentDepth+1, depth-1, alpha, beta, !maximizing, startPV, usePv && i == 0, startTime, maxTime, output)
				board.reverseMove(&move, &boardPrimitives)
				if !ab.Completed {
					return notCompletedOutput
				}
				if ab.Score > maxEval {
					maxEval = ab.Score
					bestPv = ab.Pv
					bestPv[currentDepth] = move
				}
				alpha = math.Max(ab.Score, alpha)
				if beta <= alpha {
					completedOutput.Score = maxEval
					completedOutput.Pv = bestPv
					return completedOutput
				}
			}
		}
		completedOutput.Score = maxEval
		completedOutput.Pv = bestPv
		return completedOutput
	} else { // minimizing player
		minEval := math.Inf(1)
		for i, om := range orderedMoves {
			select {
			case <-stopPondering:
				fmt.Println("stop pondering inside")
				return notCompletedOutput
			default:
				move := om.move
				if time.Since(startTime) >= maxTime && completedOnce {
					return notCompletedOutput
				}

				boardPrimitives := board.getBoardPrimitives()
				board.Move(&move)
				ab := board.alphaBetaPruning(stopPondering, completedOnce, currentDepth+1, depth-1, alpha, beta, !maximizing, startPV, usePv && i == 0, startTime, maxTime, output)
				board.reverseMove(&move, &boardPrimitives)
				if !ab.Completed {
					return notCompletedOutput
				}
				if ab.Score < minEval {
					minEval = ab.Score
					bestPv = ab.Pv
					bestPv[currentDepth] = move
				}
				beta = math.Min(ab.Score, beta)
				if beta <= alpha {
					completedOutput.Score = minEval
					completedOutput.Pv = bestPv
					return completedOutput
				}
			}
		}
		completedOutput.Score = minEval
		completedOutput.Pv = bestPv
		return completedOutput
	}
}
