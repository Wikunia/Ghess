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

var board = ghess.Board{}

func main() {
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
		board = ghess.GetBoardFromFen(START_FEN)
		if len(commands) > 2 {
			if commands[2] == "moves" {
				makeMoves(commands[3:])
			}
		}
	case "fen":
		board = ghess.GetBoardFromFen(strings.Join(commands[2:], " "))
	default:
		fmt.Println("can't handle that command atm")
	}
}

func makeMoves(moves []string) {
	for _, moveStr := range moves {
		board.MoveLongAlgebraic(moveStr)
	}
}

func handleGo(in string) {
	move := board.AlphaBetaEngineMove()
	fmt.Printf("bestmove %s\n", ghess.GetAlgebraicFromMove(&move))
}
