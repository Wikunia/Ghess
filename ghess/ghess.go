package ghess

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"unicode"

	fiber "github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/mustache"
	websocket "github.com/gofiber/websocket/v2"
)

type Position struct {
	x int
	y int
}

type Piece struct {
	position Position
	c        rune
	color    bool
}

type Board struct {
	position [8][8]int
	pieces   map[int]Piece
	color    bool
}

type JSONMove struct {
	PieceId   int `json:"pieceId"`
	CaptureId int `json:"captureId"`
	ToY       int `json:"toY"`
	ToX       int `json:"toX"`
}

var currentBoard Board
var websocketConns map[int]*websocket.Conn
var nextConnectionId = 1

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
	setFen(fen)
	return displayBoard(currentBoard)
}

func displayBoard(board Board) string {
	result := displayGround()
	short2full := short2full_name()
	for pieceId, piece := range currentBoard.pieces {
		pieceName := short2full[piece.c]
		position := piece.position
		left := strconv.Itoa((position.x - 1) * 10)
		top := strconv.Itoa((position.y - 1) * 10)
		result += `<div class="piece" draggable="true" ondragstart="onDragStart(event);" ondrop="onDrop(event);" ondragover="onDragOver(event);"> 
				<img id="piece_` + strconv.Itoa(pieceId) + `" src="images/` + pieceName + `.png" style="left: ` + left + `vmin; top: ` + top + `vmin;"/>
			</div>`
	}
	return result
}

func displayGround() string {
	result := ""
	colors := []string{"white", "black"}
	for i := 1; i < 9; i++ {
		result += `<div class="board_row">`
		for j := 1; j < 9; j++ {
			color := colors[(i+j)%2]
			result += `<div id="square_` + strconv.Itoa(i) + `_` + strconv.Itoa(j) + `" class="square square_` + color + `" ondrop="onDrop(event);" ondragover="onDragOver(event);"> </div>`
		}
		result += `</div>`
	}
	return result
}

func setFen(fen string) {
	parts := strings.Split(fen, " ")
	fen_pieces := parts[0]
	rows := strings.Split(fen_pieces, "/")
	var position [8][8]int
	pieces := make(map[int]Piece)
	pieceId := 1
	for r, row := range rows {
		cpos := 0
		for _, p := range row {
			if (p > 'a' && p < 'z') || (p > 'A' && p < 'Z') {
				position[r][cpos] = pieceId
				pieces[pieceId] = Piece{position: Position{x: cpos + 1, y: r + 1}, c: p, color: p == unicode.ToLower(p)}
				pieceId += 1
				cpos += 1
			} else {
				// convert rune to integer
				n, _ := strconv.Atoi(string(p))
				for i := 0; i < n; i++ {
					position[r][cpos+i] = 0
				}
				cpos += n
			}
		}
	}
	color := false
	if parts[1][0] == 'b' {
		color = true
	}

	currentBoard = Board{position: position, pieces: pieces, color: color}
}

func fillMove(m *JSONMove) error {
	if m.ToX == 0 {
		captureId := m.CaptureId
		m.ToX = currentBoard.pieces[captureId].position.x
		m.ToY = currentBoard.pieces[captureId].position.y
	}

	return nil
}

func move(m *JSONMove) error {
	fromX := currentBoard.pieces[m.PieceId].position.x
	fromY := currentBoard.pieces[m.PieceId].position.y
	currentBoard.position[m.ToY-1][m.ToX-1] = currentBoard.position[fromY-1][fromX-1]
	currentBoard.position[fromY-1][fromX-1] = 0
	if thisPiece, ok := currentBoard.pieces[m.PieceId]; ok {
		thisPiece.position.x = m.ToX
		thisPiece.position.y = m.ToY
		currentBoard.pieces[m.PieceId] = thisPiece
	} else {
		return fmt.Errorf("Should exist")
	}
	if m.CaptureId != 0 {
		delete(currentBoard.pieces, m.CaptureId)
	}
	return nil
}

func engineMove() JSONMove {
	return JSONMove{PieceId: 13, CaptureId: 0, ToY: 4, ToX: 5}
}

func getPieceColor(piece int) bool {
	return currentBoard.pieces[piece].color
}

func isFree(y, x int) bool {
	return currentBoard.position[y-1][x-1] == 0
}

func isLegalPawn(m *JSONMove) bool {
	color := currentBoard.color
	if !color {
		return isLegalPawnWhite(m)
	}
	return isLegalPawnBlack(m)

}

func isLegalPawnWhite(m *JSONMove) bool {
	fromX := currentBoard.pieces[m.PieceId].position.x
	fromY := currentBoard.pieces[m.PieceId].position.y

	// normal move
	if m.CaptureId == 0 {
		if fromX != m.ToX {
			return false
		}
		if fromY-m.ToY == 2 {
			if fromY != 7 {
				return false
			}
			return isFree(fromY-1, fromX)
		}
		return fromY-m.ToY == 1
	} else { // capture
		diff := fromX - m.ToX
		if (diff == 1 || diff == -1) && fromY-m.ToY == 1 {
			return true
		}
	}
	return false
}

func isLegalPawnBlack(m *JSONMove) bool {
	fromX := currentBoard.pieces[m.PieceId].position.x
	fromY := currentBoard.pieces[m.PieceId].position.y

	// normal move
	if m.CaptureId == 0 {
		if fromX != m.ToX {
			return false
		}
		if m.ToY-fromY == 2 {
			if fromY != 2 {
				return false
			}
			return isFree(fromY+1, fromX)
		}
		return m.ToY-fromY == 1
	} else { // capture
		diff := fromX - m.ToX
		if (diff == 1 || diff == -1) && m.ToY-fromY == 1 {
			return true
		}
	}
	return false
}

func isLegal(m *JSONMove) bool {
	pieceType := unicode.ToLower(rune(currentBoard.pieces[m.PieceId].c))
	fmt.Println(pieceType == 'p')
	switch pieceType {
	case 'p':
		return isLegalPawn(m)
	}
	return false
}

func Run() {

	websocketConns = make(map[int]*websocket.Conn)
	// Create a new engine
	engine := mustache.NewFileSystem(http.Dir("./../ghess/public/templates"), ".mustache")

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Static("/", "./../ghess/public")

	setFen("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")

	app.Get("/", func(c *fiber.Ctx) error {
		// Render index
		return c.Render("index", fiber.Map{
			"board": displayBoard(currentBoard),
		})
	})

	// websocket
	app.Use("/ws", func(c *fiber.Ctx) error {
		// IsWebSocketUpgrade returns true if the client
		// requested upgrade to the WebSocket protocol.
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	// Upgraded websocket request
	app.Get("/ws", websocket.New(func(c *websocket.Conn) {
		type JSONWelcome struct {
			Id int `json:"connectionId"`
		}
		websocketConns[nextConnectionId] = c
		c.WriteJSON(JSONWelcome{Id: nextConnectionId})
		nextConnectionId += 1
		for {

			var moveObj JSONMove
			err := c.ReadJSON(&moveObj)
			if err != nil {
				log.Println("read:", err)
				break
			}
			log.Printf("recv: %v\n", moveObj)
			if moveObj.PieceId != 0 && getPieceColor(moveObj.PieceId) == currentBoard.color {
				err = fillMove(&moveObj)
				if err != nil {
					log.Println("read:", err)
					break
				}
				fmt.Printf("move: %v\n", moveObj)
				legal := isLegal(&moveObj)

				if legal {
					err = move(&moveObj)
					if err != nil {
						log.Println("read:", err)
						break
					}
					err = c.WriteJSON(moveObj)
					if err != nil {
						log.Println("read:", err)
						break
					}
					currentBoard.color = !currentBoard.color
				}

				/*
					// calculate response move
					engineMoveObj := engineMove()
					move(&engineMoveObj)
					c.WriteJSON(engineMoveObj)

					currentBoard.color = !currentBoard.color
				*/
			}
		}
	}))

	log.Fatal(app.Listen(":3000"))
}
