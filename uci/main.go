package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/Wikunia/Ghess/ghess"
)

const ENGINE_NAME = "Ghess v0.1.0"
const AUTHOR_NAME = "Ole Kroeger"
const START_FEN = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"

var currentFEN = ""

var board = ghess.Board{}
var isready = make(chan bool)
var stopPondering = make(chan bool)
var currentBestPv = make(chan [30]ghess.Move)
var startPv = [30]ghess.Move{}
var currentlyPondering = false
var madeMoves = []ghess.Move{}
var ponderingMove = ""
var maxThinkingTime = 0

func main() {
	go func() {
		for {
			startPv = <-currentBestPv
		}
	}()

	input := make(chan string)
	go func(in chan string) {
		reader := bufio.NewReader(os.Stdin)
		for {
			s, err := reader.ReadString('\n')
			if err != nil {
				close(in)
				log.Println("Error in read string", err)
			}
			in <- s
		}
	}(input)
	for {
		select {
		case in := <-input:
			in = strings.TrimSpace(in)
			stillRunning := runCommand(in)
			if !stillRunning {
				return
			}
		}
	}
}

func runCommand(in string) bool {
	command := strings.Split(in, " ")
	switch command[0] {
	case "uci":
		printUCI()
	case "isready":
		fmt.Println("readyok")
		go func() {
			currentBestPv <- [30]ghess.Move{}
			isready <- true
		}()
	case "position":
		handlePosition(in)
	case "go":
		handleGo(in)
	case "stop":
		handleStop()
	case "ponderhit":
		handlePonderHit(in)
	case "quit":
		return false
	}
	return true
}

func printUCI() {
	fmt.Printf("id name %s\n", ENGINE_NAME)
	fmt.Printf("id author %s\n", AUTHOR_NAME)
	fmt.Println("uciok")
}

func handlePosition(in string) {
	commands := strings.Split(in, " ")
	switch commands[1] {
	case "startpos":
		currentFEN = START_FEN
		board = ghess.GetBoardFromFen(START_FEN)
		if len(commands) > 2 {
			if commands[2] == "moves" {
				makeMoves(commands[3:])
			}
		}
	case "fen":
		currentFEN = strings.Join(commands[2:8], " ")
		board = ghess.GetBoardFromFen(currentFEN)
		if len(commands) > 8 {
			if commands[8] == "moves" {
				makeMoves(commands[9:])
			}
		}
	default:
		fmt.Println("can't handle that command atm")
	}
}

func makeMoves(moves []string) {
	madeMoves = []ghess.Move{}
	for _, moveStr := range moves {
		move, err := board.GetMoveFromLongAlgebraic(moveStr)
		if err != nil {
			fmt.Println("ERROR for move ", moveStr, " ", err)
		}
		board.Move(&move)
		madeMoves = append(madeMoves, move)
	}
}

func handleStop() {
	if currentlyPondering {
		stopPondering <- true
		currentlyPondering = false
		// wait until pondering stopped
		ready := <-isready
		if ready {
			go func() {
				isready <- true
			}()
		}
		if startPv[0].PieceId != 0 {
			fmt.Printf("bestmove %s\n", ghess.GetAlgebraicFromMove(&startPv[0]))
		} else {
			// not sure whether there is a bestmove expected
			fmt.Println("bestmove h8h8")
		}
	} else {
		// not sure whether there is a bestmove expected
		fmt.Println("bestmove a8a8")
	}
}

func handleGo(in string) {
	parts := strings.Split(in, " ")
	if parts[1] == "ponder" {
		handleGoPonder(in)
		return
	}
	// default 4s*40
	wtime := 40 * 4000
	btime := 40 * 4000
	winc := 0
	binc := 0
	for i := 0; i < len(parts); i++ {
		if parts[i] == "wtime" {
			wtime, _ = strconv.Atoi(parts[i+1])
		}
		if parts[i] == "btime" {
			btime, _ = strconv.Atoi(parts[i+1])
		}
		if parts[i] == "winc" {
			winc, _ = strconv.Atoi(parts[i+1])
		}
		if parts[i] == "binc" {
			binc, _ = strconv.Atoi(parts[i+1])
		}
	}
	if board.IsBlacksTurn {
		maxThinkingTime = btime/40 + binc
	} else {
		maxThinkingTime = wtime/40 + winc
	}

	ready := <-isready
	fmt.Println("ready: ", ready)
	depth := 2
	startPv = [30]ghess.Move{}

	fmt.Println("run with depth: ", depth)
	ab := board.AlphaBetaEngineMove(startPv, depth, 30, false, true, maxThinkingTime)
	pv := ab.Pv
	if pv[1].PieceId != 0 {
		currentlyPondering = true
		ponderingMove = ghess.GetAlgebraicFromMove(&pv[1])
		fmt.Printf("bestmove %s ponder %s\n", ghess.GetAlgebraicFromMove(&pv[0]), ponderingMove)
	} else {
		fmt.Printf("bestmove %s\n", ghess.GetAlgebraicFromMove(&pv[0]))
		currentlyPondering = false
		go func() {
			isready <- true
		}()
	}
}

func handleGoPonder(in string) {
	ended, _, _ := board.CheckGameEnded()
	if !ended {
		currentlyPondering = true
		stopPondering = make(chan bool)
		go board.AlphaBetaEnginePonder(stopPondering, isready, currentBestPv)
	} else {
		currentlyPondering = false
		go func() {
			isready <- true
		}()
	}
}

func handlePonderHit(in string) {
	if currentlyPondering {
		stopPondering <- true
		currentlyPondering = false
	}
	ready := <-isready
	fmt.Println("ready: ", ready)
	currentFEN = START_FEN
	board = ghess.GetBoardFromFen(START_FEN)
	for _, move := range madeMoves {
		board.Move(&move)
	}

	depth := 2
	for i := 0; i < 30; i++ {
		if startPv[i].PieceId != 0 {
			depth = i + 2
		}
	}
	if depth < 2 {
		depth = 2
	}
	fmt.Println("Used pondering")
	ab := board.AlphaBetaEngineMove(startPv, depth, 30, true, true, maxThinkingTime)
	pv := ab.Pv
	if pv[1].PieceId != 0 {
		ponderingMove = ghess.GetAlgebraicFromMove(&pv[1])
		fmt.Printf("bestmove %s ponder %s\n", ghess.GetAlgebraicFromMove(&pv[0]), ponderingMove)
	} else {
		fmt.Printf("bestmove %s\n", ghess.GetAlgebraicFromMove(&pv[0]))
	}
}
