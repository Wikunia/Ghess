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
	position           [8][8]int
	pieces             map[int]Piece
	color              bool
	white_castle_king  bool
	white_castle_queen bool
	black_castle_king  bool
	black_castle_queen bool
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

	currentBoard = Board{
		position:           position,
		pieces:             pieces,
		color:              color,
		white_castle_king:  strings.ContainsRune(parts[2], 'K'),
		white_castle_queen: strings.ContainsRune(parts[2], 'Q'),
		black_castle_king:  strings.ContainsRune(parts[2], 'k'),
		black_castle_queen: strings.ContainsRune(parts[2], 'q'),
	}
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
	piece := currentBoard.pieces[m.PieceId]
	fromX := piece.position.x
	fromY := piece.position.y
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

	// disallow castling
	if isKing(piece) {
		if !piece.color {
			currentBoard.white_castle_king = false
			currentBoard.white_castle_queen = false
		} else {
			currentBoard.black_castle_king = false
			currentBoard.black_castle_queen = false
		}
	}
	return nil
}

func engineMove() JSONMove {
	return JSONMove{PieceId: 13, CaptureId: 0, ToY: 4, ToX: 5}
}

func getPieceColor(piece int) bool {
	return currentBoard.pieces[piece].color
}

func getRookMoveIfCastle(m *JSONMove) (JSONMove, bool) {
	piece := currentBoard.pieces[m.PieceId]
	rm := JSONMove{PieceId: 0, CaptureId: 0, ToX: 0, ToY: 0}
	if !isKing(piece) {
		return rm, false
	}
	color := piece.color
	fromX := currentBoard.pieces[m.PieceId].position.x
	toX := m.ToX
	diffx := toX - fromX
	fmt.Println("diffx: ", diffx)

	if !color {
		if diffx == 2 {
			rook_id := currentBoard.position[7][7]
			rm = JSONMove{PieceId: rook_id, CaptureId: 0, ToX: 6, ToY: 8}
		} else {
			rook_id := currentBoard.position[7][0]
			rm = JSONMove{PieceId: rook_id, CaptureId: 0, ToX: 4, ToY: 8}
		}
	} else {
		if diffx == 2 {
			rook_id := currentBoard.position[0][7]
			rm = JSONMove{PieceId: rook_id, CaptureId: 0, ToX: 6, ToY: 1}
		} else {
			rook_id := currentBoard.position[0][0]
			rm = JSONMove{PieceId: rook_id, CaptureId: 0, ToX: 4, ToY: 1}
		}
	}
	return rm, true
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
					rm, isCastle := getRookMoveIfCastle(&moveObj)
					if isCastle {
						err = move(&rm)
						if err != nil {
							log.Println("read:", err)
							break
						}
						err = c.WriteJSON(rm)
						if err != nil {
							log.Println("read:", err)
							break
						}
					}
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
