package ghess

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
	"time"
)

type Eval struct {
	move  Move
	score float64
}

func (board *Board) staticEvaluation() float64 {
	gameEnded, endType, _ := board.checkGameEnded()
	if gameEnded {
		if endType == "checkmate" {
			if board.isBlacksTurn {
				return 100000.0
			} else {
				return -100000.0
			}
		} else if endType == "draw" {
			return 0.0
		}
	}
	return float64(board.countMaterialOfColor(false) - board.countMaterialOfColor(true))
}

func (board *Board) getPossibleMovesOrdered(isMaximizing bool) []Move {
	evals := []Eval{}
	orderedMoves := board.getPossibleMoves()
	for _, move := range orderedMoves {
		score := board.staticEvaluation()
		if !isMaximizing {
			score = -score
		}
		// we look for highest score so if not maximizing
		evals = append(evals, Eval{move: move, score: score})
	}
	sort.Slice(evals, func(i, j int) bool {
		return evals[i].score > evals[j].score
	})
	for i := 0; i < len(orderedMoves); i++ {
		orderedMoves[i] = evals[i].move
	}
	return orderedMoves
}

func (board *Board) alphaBetaEngineMove() Move {
	currentDepth := 2
	orderedMoves := board.getPossibleMoves()
	evals := []Eval{}
	startTime := time.Now()
	completelyEvaluated := []Eval{}
	calcuatedAtLeastOnce := false
	myColor := board.isBlacksTurn

	for time.Since(startTime) <= MAX_ENGINE_TIME*time.Second {
		evals = []Eval{}
		inTime := true
		for _, move := range orderedMoves {
			if time.Since(startTime) >= MAX_ENGINE_TIME*time.Second && calcuatedAtLeastOnce {
				inTime = false
				break
			}
			boardPrimitives := board.getBoardPrimitives()
			board.Move(&move)
			eval := board.alphaBetaPruning(currentDepth-1, math.Inf(-1), math.Inf(1), !myColor)
			evals = append(evals, Eval{move: move, score: -eval})
			board.reverseMove(&move, &boardPrimitives)
		}
		if inTime {
			sort.Slice(evals, func(i, j int) bool {
				return evals[i].score > evals[j].score
			})
			for i := 0; i < len(orderedMoves); i++ {
				orderedMoves[i] = evals[i].move
			}
			currentDepth += 1
			completelyEvaluated = evals
			calcuatedAtLeastOnce = true
		}

	}

	fmt.Printf("evalutad up to depth %d in %.02f\n", currentDepth-1, time.Since(startTime).Seconds())

	sort.Slice(completelyEvaluated, func(i, j int) bool {
		return completelyEvaluated[i].score > completelyEvaluated[j].score
	})

	bestScore := completelyEvaluated[0].score
	fmt.Println("bestScore: ", bestScore)
	whenWorse := len(completelyEvaluated)
	for i := range completelyEvaluated {
		if completelyEvaluated[i].score < bestScore {
			whenWorse = i
			break
		}
	}
	moveId := rand.Intn(whenWorse)
	fmt.Println("score: ", completelyEvaluated[moveId])
	return completelyEvaluated[moveId].move
}

func (board *Board) alphaBetaPruning(depth int, alpha, beta float64, maximizing bool) float64 {
	gameEnded, _, _ := board.checkGameEnded()
	if depth == 0 || gameEnded {
		return board.staticEvaluation()
	}

	// maximizing player
	if maximizing {
		maxEval := math.Inf(-1)
		for _, move := range board.getPossibleMovesOrdered(true) {
			boardPrimitives := board.getBoardPrimitives()
			board.Move(&move)
			eval := board.alphaBetaPruning(depth-1, alpha, beta, !maximizing)
			board.reverseMove(&move, &boardPrimitives)
			maxEval = math.Max(eval, maxEval)
			alpha = math.Max(alpha, eval)
			if beta <= alpha {
				break
			}
		}
		return maxEval
	} else { // minimizing player
		minEval := math.Inf(1)
		for _, move := range board.getPossibleMovesOrdered(false) {
			boardPrimitives := board.getBoardPrimitives()
			board.Move(&move)
			eval := board.alphaBetaPruning(depth-1, alpha, beta, !maximizing)
			board.reverseMove(&move, &boardPrimitives)
			minEval = math.Min(eval, minEval)
			beta = math.Min(beta, eval)
			if beta <= alpha {
				break
			}
		}
		return minEval
	}
}
