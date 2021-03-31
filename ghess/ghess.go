package ghess

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	fiber "github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/mustache"
)

type Position struct {
	x int
	y int
}

type Board struct {
	fen            string
	position       [8][8]rune
	piece2Position map[int]Position
	piece2Rune     map[int]rune
}

var currentBoard Board

func short2full_name() map[rune]string {
	return map[rune]string{
		'K': "white_king",
		'Q': "white_queen",
		'B': "white_bishop",
		'N': "white_knight",
		'R': "white_rook",
		'P': "white_pawn",
		'k': "black_king",
		'q': "black_queen",
		'b': "black_bishop",
		'n': "black_knight",
		'r': "black_rook",
		'p': "black_pawn",
	}
}

func displayFen(fen string) string {
	currentBoard = parseFen(fen)
	return displayBoard(currentBoard)
}

func displayBoard(board Board) string {
	result := displayGround()
	short2full := short2full_name()
	for pieceId, position := range currentBoard.piece2Position {
		pieceName := short2full[currentBoard.piece2Rune[pieceId]]
		left := strconv.Itoa(position.x * 10)
		top := strconv.Itoa(position.y * 10)
		result += `<div class="piece" draggable="true" ondragstart="onDragStart(event);" ondrop="onDrop(event);" ondragover="onDragOver(event);"> 
				<img id="piece_` + strconv.Itoa(pieceId) + `" src="images/` + pieceName + `.png" style="left: ` + left + `vmin; top: ` + top + `vmin;"/>
			</div>`
	}
	return result
}

func displayGround() string {
	result := ""
	colors := []string{"white", "black"}
	for i := 0; i < 8; i++ {
		result += `<div class="board_row">`
		for j := 0; j < 8; j++ {
			color := colors[(i+j)%2]
			result += `<div id="square_` + strconv.Itoa(i) + `_` + strconv.Itoa(j) + `" class="square square_` + color + `" ondrop="onDrop(event);" ondragover="onDragOver(event);"> </div>`
		}
		result += `</div>`
	}
	return result
}

func parseFen(fen string) Board {
	parts := strings.Split(fen, " ")
	pieces := parts[0]
	rows := strings.Split(pieces, "/")
	var position [8][8]rune
	piece2Position := make(map[int]Position)
	piece2Rune := make(map[int]rune)
	pieceId := 0
	for r, row := range rows {
		cpos := 0
		for _, p := range row {
			if (p > 'a' && p < 'z') || (p > 'A' && p < 'Z') {
				position[r][cpos] = p
				piece2Position[pieceId] = Position{x: cpos, y: r}
				piece2Rune[pieceId] = p
				pieceId += 1
				cpos += 1
			} else {
				// convert rune to integer
				n, _ := strconv.Atoi(string(p))
				for i := 0; i < n; i++ {
					position[r][cpos+i] = '0'
				}
				cpos += n
			}
		}
	}
	return Board{fen: fen, position: position, piece2Position: piece2Position, piece2Rune: piece2Rune}
}

func apiMakeMove(c *fiber.Ctx, capture bool) error {
	type JSONMove struct {
		PieceId   int `json:"pieceId"`
		captureId int `json:"captureId"`
		ToY       int `json:"to_y"`
		ToX       int `json:"to_x"`
	}

	var move JSONMove
	fmt.Println(string(c.Body()))

	err := c.BodyParser(&move)

	// if error
	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Cannot parse JSON",
		})
	}

	if !capture {
		makeMove(move.PieceId, move.ToY, move.ToX)
	} else {
		makeCapture(move.PieceId, move.captureId)
	}

	return c.SendString("success")
}

func makeMove(pieceId, toY, toX int) string {
	fromX := currentBoard.piece2Position[pieceId].x
	fromY := currentBoard.piece2Position[pieceId].y
	currentBoard.position[toY][toX] = currentBoard.position[fromY][fromX]
	currentBoard.position[fromY][fromX] = '0'
	currentBoard.piece2Position[pieceId] = Position{x: toX, y: toY}
	return displayBoard(currentBoard)
}

func makeCapture(pieceId, captureId int) string {
	fromX := currentBoard.piece2Position[pieceId].x
	fromY := currentBoard.piece2Position[pieceId].y
	toX := currentBoard.piece2Position[captureId].x
	toY := currentBoard.piece2Position[captureId].y
	currentBoard.position[toY][toX] = currentBoard.position[fromY][fromX]
	currentBoard.position[fromY][fromX] = '0'
	currentBoard.piece2Position[pieceId] = Position{x: toX, y: toY}
	delete(currentBoard.piece2Position, pieceId)
	delete(currentBoard.piece2Rune, pieceId)
	return displayBoard(currentBoard)
}

func Run() {
	// Create a new engine
	engine := mustache.NewFileSystem(http.Dir("./../ghess/public/templates"), ".mustache")

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Static("/", "./../ghess/public")
	currentBoard = parseFen("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")

	app.Get("/", func(c *fiber.Ctx) error {
		// Render index
		return c.Render("index", fiber.Map{
			"board": displayBoard(currentBoard),
		})
	})

	app.Post("/api/move", func(c *fiber.Ctx) error { return apiMakeMove(c, false) })
	app.Post("/api/capture", func(c *fiber.Ctx) error { return apiMakeMove(c, true) })

	log.Fatal(app.Listen(":3000"))
}
