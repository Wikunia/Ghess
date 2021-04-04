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
	id       int
	position Position
	c        rune
	color    bool
	onBoard  bool
}

type Board struct {
	position            [8][8]int
	pieces              map[int]Piece
	color               bool
	white_castle_king   bool
	white_castle_queen  bool
	black_castle_king   bool
	black_castle_queen  bool
	en_passant_position Position // will be 0,0 if not possible
	halfMoves           int
	nextMove            int
	whiteKingId         int
	blackKingId         int
}

type BoardPrimitives struct {
	color               bool
	white_castle_king   bool
	white_castle_queen  bool
	black_castle_king   bool
	black_castle_queen  bool
	en_passant_position Position // will be 0,0 if not possible
	halfMoves           int
	nextMove            int
	whiteKingId         int
	blackKingId         int
}

type JSONMove struct {
	PieceId   int `json:"pieceId"`
	CaptureId int `json:"captureId"`
	ToY       int `json:"toY"`
	ToX       int `json:"toX"`
}

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

func (board *Board) getFen() string {
	fen := ""
	n := 0
	pieceId := 0
	for y := 0; y < 8; y++ {
		n = 0
		for x := 0; x < 8; x++ {
			pieceId = board.position[y][x]
			if pieceId == 0 {
				n += 1
			} else {
				if n != 0 {
					fen += strconv.Itoa(n)
					n = 0
				}
				fen += string(board.pieces[pieceId].c)
			}
		}
		if n != 0 {
			fen += strconv.Itoa(n)
		}
		if y != 7 {
			fen += "/"
		}
	}

	fen += " "
	colorInitial := "w"
	if board.color {
		colorInitial = "b"
	}
	fen += colorInitial

	fen += " "
	if board.white_castle_king {
		fen += "K"
	}
	if board.white_castle_queen {
		fen += "Q"
	}
	if board.black_castle_king {
		fen += "k"
	}
	if board.black_castle_queen {
		fen += "q"
	}
	fen += " "
	if board.en_passant_position.x != 0 {
		fen += string(rune('a' - 1 + board.en_passant_position.x))
		fen += strconv.Itoa(9 - board.en_passant_position.y)
	} else {
		fen += "-"
	}

	fen += " "
	fen += strconv.Itoa(board.halfMoves)
	fen += " "
	fen += strconv.Itoa(board.nextMove)

	return fen
}

func displayFen(fen string) string {
	board := GetBoardFromFen(fen)
	return board.display()
}

func (board *Board) display() string {
	result := displayGround()
	short2full := short2full_name()
	for pieceId, piece := range board.pieces {
		if !piece.onBoard {
			continue
		}
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

func GetBoardFromFen(fen string) Board {
	parts := strings.Split(fen, " ")
	fen_pieces := parts[0]
	rows := strings.Split(fen_pieces, "/")
	var position [8][8]int
	pieces := make(map[int]Piece)
	pieceId := 1
	whiteKingId := 0
	blackKingId := 0
	for r, row := range rows {
		cpos := 0
		for _, p := range row {
			if (p > 'a' && p < 'z') || (p > 'A' && p < 'Z') {
				position[r][cpos] = pieceId
				pieces[pieceId] = Piece{id: pieceId, position: Position{x: cpos + 1, y: r + 1}, c: p, color: p == unicode.ToLower(p), onBoard: true}
				if p == 'K' {
					whiteKingId = pieceId
				}
				if p == 'k' {
					blackKingId = pieceId
				}
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
	en_passant_position := Position{x: 0, y: 0}
	if parts[3][0] != '-' {
		en_passant_position.y = 9 - int(parts[3][1]-'0')
		en_passant_position.x = int(parts[3][0]-'a') + 1
	}
	// fmt.Println("en_passant_position: ", en_passant_position)

	halfMoves, err := strconv.Atoi(parts[4])
	if err != nil {
		fmt.Println("could not convert number of half moves to integer")
	}
	nextMove, err := strconv.Atoi(parts[5])
	if err != nil {
		fmt.Println("could not convert next move number to integer")
	}
	return Board{
		position:            position,
		pieces:              pieces,
		color:               color,
		white_castle_king:   strings.ContainsRune(parts[2], 'K'),
		white_castle_queen:  strings.ContainsRune(parts[2], 'Q'),
		black_castle_king:   strings.ContainsRune(parts[2], 'k'),
		black_castle_queen:  strings.ContainsRune(parts[2], 'q'),
		en_passant_position: en_passant_position,
		halfMoves:           halfMoves,
		nextMove:            nextMove,
		whiteKingId:         whiteKingId,
		blackKingId:         blackKingId,
	}
}

func (board *Board) fillMove(m *JSONMove) {
	if m.ToX == 0 {
		captureId := m.CaptureId
		m.ToX = board.pieces[captureId].position.x
		m.ToY = board.pieces[captureId].position.y
	} else if m.CaptureId == 0 { // fill capture if there is a piece on that position
		m.CaptureId = board.position[m.ToY-1][m.ToX-1]
	}
}

func (board *Board) move(m *JSONMove) (int, JSONMove) {
	capturedId, castledMove := board.moveTemp(m)
	// only count once for castling
	if isPawn(board.pieces[m.PieceId]) || m.CaptureId != 0 {
		board.halfMoves = 0
	} else {
		board.halfMoves += 1
	}
	if !board.color {
		board.nextMove += 1
	}
	board.color = !board.color
	return capturedId, castledMove
}

// Move the piece and return the id of a captured piece and return castled rook move 0, empty otherwise
func (board *Board) moveTemp(m *JSONMove) (int, JSONMove) {
	// before moving anything check if this is a castling move
	castledMove := JSONMove{}
	rm, isCastle := board.getRookMoveIfCastle(m)
	if isCastle {
		board.moveTemp(&rm)
		castledMove = rm
	}
	piece := board.pieces[m.PieceId]
	fromX := piece.position.x
	fromY := piece.position.y
	board.position[m.ToY-1][m.ToX-1] = board.position[fromY-1][fromX-1]
	board.position[fromY-1][fromX-1] = 0
	if thisPiece, ok := board.pieces[m.PieceId]; ok {
		thisPiece.position.x = m.ToX
		thisPiece.position.y = m.ToY
		board.pieces[m.PieceId] = thisPiece
	}
	if m.CaptureId != 0 {
		if capturedPiece, ok := board.pieces[m.CaptureId]; ok {
			capturedPiece.onBoard = false
			board.pieces[m.CaptureId] = capturedPiece
		}
		board.en_passant_position.x = 0
		board.en_passant_position.y = 0
		return m.CaptureId, castledMove
	} else {
		if isPawn(piece) {
			diffx := abs(fromX - m.ToX)
			if diffx == 1 {
				// en passant
				m.CaptureId = board.position[fromY-1][m.ToX-1]
				board.en_passant_position.x = 0
				board.en_passant_position.y = 0
				return m.CaptureId, castledMove
			}
		}
	}

	// disallow castling
	if isKing(piece) {
		if !piece.color {
			board.white_castle_king = false
			board.white_castle_queen = false
		} else {
			board.black_castle_king = false
			board.black_castle_queen = false
		}
	}

	if isRook(piece) {
		if !piece.color {
			if fromY == 8 {
				if fromX == 8 {
					board.white_castle_king = false
				} else if fromX == 1 {
					board.white_castle_queen = false
				}
			}
		} else { //black
			if fromY == 1 {
				if fromX == 8 {
					board.black_castle_king = false
				} else if fromX == 1 {
					board.black_castle_queen = false
				}
			}
		}
	}

	// check en passant
	if isPawn(piece) {
		diffy := abs(fromY - m.ToY)
		if diffy == 2 {
			board.en_passant_position.x = fromX
			board.en_passant_position.y = (fromY + m.ToY) / 2
		} else {
			board.en_passant_position.x = 0
			board.en_passant_position.y = 0
		}
	} else {
		board.en_passant_position.x = 0
		board.en_passant_position.y = 0
	}

	return 0, castledMove
}

func (board *Board) reverseMove(m *JSONMove, fromY int, fromX int, capturedId int, castledMove *JSONMove, boardPrimitives BoardPrimitives) {
	move := JSONMove{PieceId: board.position[m.ToY-1][m.ToX-1], CaptureId: capturedId, ToY: fromY, ToX: fromX}
	board.fillMove(&move)
	board.move(&move)
	if capturedId != 0 {
		if capturedPiece, ok := board.pieces[capturedId]; ok {
			capturedPiece.onBoard = true
			board.pieces[capturedId] = capturedPiece
			board.position[capturedPiece.position.y-1][capturedPiece.position.x-1] = capturedId
		}
	}
	if castledMove.PieceId != 0 {
		move.PieceId = castledMove.PieceId
		move.CaptureId = 0
		if castledMove.ToX == 6 {
			move.ToX = 8
		} else {
			move.ToX = 1
		}
		move.ToY = castledMove.ToY
		board.fillMove(&move)
		board.move(&move)
	}
	board.setBoardPrimitives(boardPrimitives)
}

func engineMove() JSONMove {
	return JSONMove{PieceId: 13, CaptureId: 0, ToY: 4, ToX: 5}
}

func (board *Board) getPieceColor(piece int) bool {
	return board.pieces[piece].color
}

func (board *Board) getRookMoveIfCastle(m *JSONMove) (JSONMove, bool) {
	piece := board.pieces[m.PieceId]
	rm := JSONMove{PieceId: 0, CaptureId: 0, ToX: 0, ToY: 0}
	if !isKing(piece) {
		return rm, false
	}
	color := piece.color
	fromX := board.pieces[m.PieceId].position.x
	toX := m.ToX
	diffx := toX - fromX
	if abs(diffx) != 2 {
		return rm, false
	}
	// this can happen in temporary reverse moves
	if fromX != 5 {
		return rm, false
	}

	if !color {
		if diffx == 2 {
			rook_id := board.position[7][7]
			rm = JSONMove{PieceId: rook_id, CaptureId: 0, ToX: 6, ToY: 8}
		} else {
			rook_id := board.position[7][0]
			rm = JSONMove{PieceId: rook_id, CaptureId: 0, ToX: 4, ToY: 8}
		}
	} else {
		if diffx == 2 {
			rook_id := board.position[0][7]
			rm = JSONMove{PieceId: rook_id, CaptureId: 0, ToX: 6, ToY: 1}
		} else {
			rook_id := board.position[0][0]
			rm = JSONMove{PieceId: rook_id, CaptureId: 0, ToX: 4, ToY: 1}
		}
	}
	return rm, true
}

func (board *Board) getMoveFromLongAlgebraic(moveStr string) (JSONMove, error) {
	move := JSONMove{}
	if len(moveStr) != 5 {
		return move, fmt.Errorf("currently only algebraic notation with 5 chars is supported")
	}
	fromX := int(moveStr[0] - 'a' + 1)
	fromY := 9 - int(moveStr[1]-'0')
	toX := int(moveStr[3] - 'a' + 1)
	toY := 9 - int(moveStr[4]-'0')
	pieceId := board.position[fromY-1][fromX-1]
	if pieceId == 0 {
		return move, fmt.Errorf("there is no piece at that position")
	}
	if board.pieces[pieceId].color != board.color {
		return move, fmt.Errorf("the piece has the wrong color")
	}
	// Todo just check if legal move
	move = JSONMove{PieceId: pieceId, CaptureId: 0, ToX: toX, ToY: toY}
	board.fillMove(&move)
	return move, nil
}

func (board *Board) MoveLongAlgebraic(moveStr string) error {
	move, err := board.getMoveFromLongAlgebraic(moveStr)
	if err != nil {
		return err
	}
	board.move(&move)
	return nil
}

func Run() {
	websocketConns = make(map[int]*websocket.Conn)
	// Create a new engine
	engine := mustache.NewFileSystem(http.Dir("./../ghess/public/templates"), ".mustache")

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Static("/", "./../ghess/public")

	board := GetBoardFromFen("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
	board = GetBoardFromFen("rn1qkbnr/p1p1pppp/bp6/8/4p3/5NP1/PPPP1PBP/RNBQK2R w KQkq - 8 5")

	app.Get("/", func(c *fiber.Ctx) error {
		// Render index
		return c.Render("index", fiber.Map{
			"board": board.display(),
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
			if moveObj.PieceId != 0 && board.getPieceColor(moveObj.PieceId) == board.color {
				board.fillMove(&moveObj)

				fmt.Printf("move: %v\n", moveObj)
				legal := board.isLegal(&moveObj)

				if legal {
					_, castledMove := board.move(&moveObj)
					if castledMove.PieceId != 0 {
						err = c.WriteJSON(castledMove)
						if err != nil {
							log.Println("read:", err)
							break
						}
					}
					fmt.Println("legal move: ", moveObj)
					err = c.WriteJSON(moveObj)
					if err != nil {
						log.Println("read:", err)
						break
					}
					fmt.Println("turn: ", board.color)
				}

				/*
					moves := board.getPossibleMoves()
					fmt.Println(moves)
					fmt.Printf("#Possible moves %d\n", len(moves))
					fmt.Println("wk: ", board.white_castle_king)
					fmt.Println("wq: ", board.white_castle_queen)
				*/
				fmt.Println(board.getFen())
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
