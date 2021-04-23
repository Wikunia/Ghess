package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
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
				fmt.Println("commands[9:]: ", commands[9:])
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
			fmt.Println("ERROR: ", err)
		}
		board.Move(&move)
		madeMoves = append(madeMoves, move)
	}
}

func handleGo(in string) {
	if currentlyPondering {
		stopPondering <- true
	}
	ready := <-isready
	fmt.Println("ready: ", ready)
	depth := 2
	lastMove := ""
	if len(madeMoves) > 0 {
		lastMove = ghess.GetAlgebraicFromMove(&madeMoves[len(madeMoves)-1])
	}
	completedOnce := false
	if len(madeMoves) > 0 && lastMove == ponderingMove {
		for i := 0; i < 30; i++ {
			if startPv[i].PieceId != 0 {
				depth = i + 2
			}
		}
		if depth < 2 {
			depth = 2
		}
		fmt.Println("Used pondering")
		completedOnce = true
	} else {
		startPv = [30]ghess.Move{}
	}

	fmt.Println("run with depth: ", depth)
	pv := board.AlphaBetaEngineMove(startPv, depth, 30, completedOnce, true, ghess.MAX_ENGINE_TIME)
	fmt.Printf("bestmove %s\n", ghess.GetAlgebraicFromMove(&pv[0]))
	ponderingMove = ghess.GetAlgebraicFromMove(&pv[1])

	copiedBoard := ghess.GetBoardFromFen(currentFEN)
	for _, move := range madeMoves {
		copiedBoard.Move(&move)
	}
	copiedBoard.Move(&pv[0])
	if pv[1].PieceId != 0 {
		copiedBoard.Move(&pv[1])
		ended, _, _ := copiedBoard.CheckGameEnded()
		if !ended {
			currentlyPondering = true
			stopPondering = make(chan bool)
			go copiedBoard.AlphaBetaEnginePonder(stopPondering, isready, currentBestPv)
		}
	}
}
