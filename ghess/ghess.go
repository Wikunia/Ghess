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

const KING = 'k'
const QUEEN = 'q'
const ROOK = 'r'
const KNIGHT = 'n'
const BISHOP = 'b'
const PAWN = 'p'

type Piece struct {
	id        int
	pos       int    // position from 0 to 63
	posB      uint64 // binary position from 0...1 to 1...0
	pieceType rune   // lower case 'k',...,'p' see KING, ... , PAWN constants
	isBlack   bool
	movementB uint64 // possible movement of this piece encoded similarly to posB
	moves     [28]int
	numMoves  int
}

func NewPiece(id int, pos int, pieceType rune, isBlack bool) Piece {
	var moves [28]int
	var posB uint64 = 1 << pos
	return Piece{id: id, pos: pos, posB: posB, pieceType: pieceType, isBlack: isBlack, movementB: 0, moves: moves, numMoves: 0}
}

type Board struct {
	pos2PieceId        [64]int
	pieces             [33]Piece // piece with id 0 should point to empty piece
	whiteIds           [16]int   // list all white piece ids
	whitePiecePosB     uint64    // the | operator applied to all whitePieces
	whitePieceMovB     uint64
	blackIds           [16]int // list all black piece ids
	blackPiecePosB     uint64  // the | operator applied to all blackPieces
	blackPieceMovB     uint64
	isBlacksTurn       bool
	white_castle_king  bool
	white_castle_queen bool
	black_castle_king  bool
	black_castle_queen bool
	en_passant_pos     int // will be -1 if not possible
	halfMoves          int
	nextMove           int
	whiteKingId        int
	blackKingId        int
	checkForChecks     bool // set to false if one wants to get movements that can capture the king even though the move itself would lead to check
	movesTilEdge       [64][8]int
}

type BoardPrimitives struct {
	isBlacksTurn       bool
	white_castle_king  bool
	white_castle_queen bool
	black_castle_king  bool
	black_castle_queen bool
	en_passant_pos     int // will be -1 if not possible
	halfMoves          int
	nextMove           int
	whiteKingId        int
	blackKingId        int
}

func NewBoard(pieces [33]Piece, whiteIds [16]int, blackIds [16]int, isBlack bool,
	white_castle_king bool,
	white_castle_queen bool,
	black_castle_king bool,
	black_castle_queen bool,
	en_passant_pos int, halfMoves int, nextMove int, whiteKingId int, blackKingId int) Board {

	var pos2PieceId [64]int
	for _, piece := range pieces {
		pos2PieceId[piece.pos] = piece.id
	}
	board := Board{
		pos2PieceId:        pos2PieceId,
		pieces:             pieces,
		whiteIds:           whiteIds,
		whitePiecePosB:     0,
		whitePieceMovB:     0,
		blackIds:           blackIds,
		blackPiecePosB:     0,
		blackPieceMovB:     0,
		isBlacksTurn:       isBlack,
		white_castle_king:  white_castle_king,
		white_castle_queen: white_castle_queen,
		black_castle_king:  black_castle_king,
		black_castle_queen: black_castle_queen,
		en_passant_pos:     en_passant_pos,
		halfMoves:          halfMoves,
		nextMove:           nextMove,
		whiteKingId:        whiteKingId,
		blackKingId:        blackKingId,
		checkForChecks:     true,
		movesTilEdge:       getMovesTilEdge(),
	}

	whitePiecePosB := board.combinePositionsOf(whiteIds)
	blackPiecePosB := board.combinePositionsOf(blackIds)
	board.whitePiecePosB = whitePiecePosB
	board.blackPiecePosB = blackPiecePosB
	board.setMovement()
	return board
}

type JSONMove struct {
	PieceId   int `json:"pieceId"`
	CaptureId int `json:"captureId"`
	ToY       int `json:"toY"`
	ToX       int `json:"toX"`
}

type Move struct {
	pieceId   int
	captureId int
	from      int
	to        int
}

type JSONRequest struct {
	RequestType string `json:"requestType"`
	PieceId     int    `json:"pieceId"`
	CaptureId   int    `json:"captureId"`
	To          int    `json:"to"`
}

type JSONSurrounding struct {
	RequestType string     `json:"requestType"`
	Surrounding [8][8]bool `json:"surrounding"`
}

var websocketConns map[int]*websocket.Conn
var nextConnectionId = 1

func getPieceName(piece *Piece) string {
	color := "white"
	if piece.isBlack {
		color = "black"
	}
	switch piece.pieceType {
	case 'k':
		return color + "_king"
	case 'q':
		return color + "_queen"
	case 'b':
		return color + "_bishop"
	case 'n':
		return color + "_knight"
	case 'r':
		return color + "_rook"
	case 'p':
		return color + "_pawn"
	}
	return "NONAME"
}

func (board *Board) getFen() string {
	fen := ""
	n := 0
	pieceId := 0
	for y := 0; y < 8; y++ {
		n = 0
		for x := 0; x < 8; x++ {
			p := y*8 + x
			pieceId = board.pos2PieceId[p]
			if pieceId == 0 {
				n += 1
			} else {
				if n != 0 {
					fen += strconv.Itoa(n)
					n = 0
				}
				pieceType := board.pieces[pieceId].pieceType
				if !board.pieces[pieceId].isBlack {
					fen += string(unicode.ToUpper(pieceType))
				} else {
					fen += string(pieceType)
				}
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
	isBlackInitial := "w"
	if board.isBlacksTurn {
		isBlackInitial = "b"
	}
	fen += isBlackInitial

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
	if board.en_passant_pos >= 0 {
		x, y := xy(board.en_passant_pos)
		fen += string(rune('a' + x))
		fen += strconv.Itoa(8 - y)
	} else {
		fen += "-"
	}

	fen += " "
	fen += strconv.Itoa(board.halfMoves)
	fen += " "
	fen += strconv.Itoa(board.nextMove)

	return fen
}

/*
func (board *Board) moveToToLongAlgebraic(fromY, fromX, toY, toX int) string {
	res := string(rune('a' + fromX))
	res += strconv.Itoa(8 - fromY)
	res += "-" + string(rune('a'+toX))
	res += strconv.Itoa(8 - toY)
	return res
}
*/

func displayFen(fen string) string {
	board := GetBoardFromFen(fen)
	return board.display()
}

func (board *Board) display() string {
	result := displayGround()
	for pieceId, piece := range board.pieces {
		if piece.posB == 0 {
			continue
		}
		pieceName := getPieceName(&board.pieces[pieceId])
		position := piece.pos
		x, y := xy(position)
		left := strconv.Itoa(x * 10)
		top := strconv.Itoa(y * 10)
		result += `<div class="piece" draggable="true" onclick="onClick(event);" ondragstart="onDragStart(event);" ondragend="onDragEnd(event);" ondrop="onDrop(event);" ondragover="onDragOver(event);">
				<img id="piece_` + strconv.Itoa(pieceId) + `" src="images/` + pieceName + `.png" style="left: ` + left + `vmin; top: ` + top + `vmin;"/>
			</div>`
	}
	return result
}

func displayGround() string {
	result := ""
	colors := []string{"white", "black"}
	id := 0
	for i := 0; i < 8; i++ {
		result += `<div class="board_row">`
		for j := 0; j < 8; j++ {
			color := colors[(i+j)%2]
			left := strconv.Itoa(j * 10)
			top := strconv.Itoa(i * 10)
			result += `<div id="square_` + strconv.Itoa(id) + `" class="square square_` + color + `"> </div>`
			result += `<div id="square_` + strconv.Itoa(id) + `_overlay" class="square_overlay" ondrop="onDrop(event);" ondragover="onDragOver(event);" style="left: ` + left + `vmin; top: ` + top + `vmin;"> </div>`
			id++
		}
		result += `</div>`
	}
	return result
}

func GetBoardFromFen(fen string) Board {
	parts := strings.Split(fen, " ")
	fen_pieces := parts[0]
	rows := strings.Split(fen_pieces, "/")
	var pieces [33]Piece
	pieceId := 1
	whiteKingId := 0
	blackKingId := 0
	var blackIds [16]int
	var whiteIds [16]int
	numBlack := 0
	numWhite := 0

	for r, row := range rows {
		cpos := 0
		for _, p := range row {
			if (p > 'a' && p < 'z') || (p > 'A' && p < 'Z') {
				isBlack := p == unicode.ToLower(p)
				pieces[pieceId] = NewPiece(pieceId, cpos+r*8, unicode.ToLower(p), isBlack)
				if isBlack {
					blackIds[numBlack] = pieceId
					numBlack++
				} else {
					whiteIds[numWhite] = pieceId
					numWhite++
				}
				if p == 'K' {
					whiteKingId = pieceId
				}
				if p == 'k' {
					blackKingId = pieceId
				}
				pieceId += 1
				cpos += 1
			} else {
				n, _ := strconv.Atoi(string(p))
				cpos += n
			}
		}
	}
	isBlack := false
	if parts[1][0] == 'b' {
		isBlack = true
	}
	en_passant_pos := -1
	if parts[3][0] != '-' {
		y := 8 - int(parts[3][1]-'0')
		x := int(parts[3][0] - 'a')
		en_passant_pos = x + 8*y
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
	white_castle_king := strings.ContainsRune(parts[2], 'K')
	white_castle_queen := strings.ContainsRune(parts[2], 'Q')
	black_castle_king := strings.ContainsRune(parts[2], 'k')
	black_castle_queen := strings.ContainsRune(parts[2], 'q')
	board := NewBoard(pieces, whiteIds, blackIds, isBlack,
		white_castle_king,
		white_castle_queen,
		black_castle_king,
		black_castle_queen,
		en_passant_pos, halfMoves, nextMove, whiteKingId, blackKingId)
	return board
}

func engineMove() JSONMove {
	return JSONMove{PieceId: 13, CaptureId: 0, ToY: 4, ToX: 5}
}

func (board *Board) getPieceisBlack(piece int) bool {
	return board.pieces[piece].isBlack
}

/*
func (board *Board) getRookMoveIfCastle(m *JSONMove) (JSONMove, bool) {
	piece := board.pieces[m.PieceId]
	rm := JSONMove{PieceId: 0, CaptureId: 0, ToX: 0, ToY: 0}
	if piece.pieceType != KING {
		return rm, false
	}
	isBlack := piece.isBlack
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

	if !isBlack {
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

func (board *Board) getMoveFromLongAlgebraic(moveStr string) (Move, error) {
	move := Move{}
	if len(moveStr) != 5 {
		return move, fmt.Errorf("currently only algebraic notation with 5 chars is supported")
	}
	fromX := int(moveStr[0] - 'a')
	fromY := 8 - int(moveStr[1]-'0')
	toX := int(moveStr[3] - 'a')
	toY := 8 - int(moveStr[4]-'0')
	pieceId := board.position[fromY][fromX]
	if pieceId == 0 {
		return move, fmt.Errorf("there is no piece at that position")
	}
	if board.pieces[pieceId].isBlack != board.isBlack {
		return move, fmt.Errorf("the piece has the wrong color")
	}
	move = board.NewMove(pieceId, 0, toY, toX, false)
	if board.isLegal(&move) {
		// capture will be filled automatically
		return move, nil
	}
	return move, fmt.Errorf("the move is not legal")
}

func (board *Board) MoveLongAlgebraic(moveStr string) error {
	move, err := board.getMoveFromLongAlgebraic(moveStr)
	if err != nil {
		return err
	}
	board.move(&board.pieces[move.pieceId], move.toY, move.toX)
	return nil
}

// check if the move is really legal based on the movement array
func (board *Board) isLegal(move *Move) bool {
	piece := board.pieces[move.pieceId]
	return piece.movement[move.toY][move.toX] && piece.isBlack == board.isBlack
}
*/

func Run() {
	websocketConns = make(map[int]*websocket.Conn)
	// Create a new engine
	engine := mustache.NewFileSystem(http.Dir("./../ghess/public/templates"), ".mustache")

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Static("/", "./../ghess/public")

	board := GetBoardFromFen("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
	// board := GetBoardFromFen("4k2r/5pp1/8/6Pp/8/8/6PP/4K2R w K h6 0 1")
	// board := GetBoardFromFen("8/5r2/8/8/2B5/8/4Q3/8 w - - 0 1")
	// board := GetBoardFromFen("rnbqkbnr/pppp1ppp/8/4p3/8/5N2/PPPP1PPP/4K3 b KQkq - 0 1")
	// board := GetBoardFromFen("r3k2r/p1ppqpb1/bn2pnp1/3PN3/Pp2P3/2N2Q1p/1PPB1PPP/R3K2R w KQkq a3 0 0")
	// board := GetBoardFromFen("r3k2r/p1ppqpb1/1n2pnp1/1b1PN3/Pp2P3/5Q1p/1PPB1PPP/R3K2R w KQkq - 0 0")

	// board.GetNumberOfMoves(2, false)
	// fmt.Println("white castle king: ", board.white_castle_king)
	// fmt.Println("isBlacksTurn: ", board.isBlacksTurn)

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

			var jsonObj JSONRequest
			err := c.ReadJSON(&jsonObj)
			if err != nil {
				log.Println("read:", err)
				break
			}
			isMove := false
			move := Move{}
			switch jsonObj.RequestType {
			case "movement":
				c.WriteJSON(JSONSurrounding{RequestType: "surrounding", Surrounding: bits2array(board.pieces[jsonObj.PieceId].movementB)})
			case "move", "capture":
				isMove = true
				move = board.NewMove(jsonObj.PieceId, jsonObj.CaptureId, jsonObj.To)
			}

			if isMove && board.isLegal(&move) {
				rookMove := board.Move(&move)
				err = c.WriteJSON(JSONRequest{RequestType: "move", PieceId: move.pieceId, CaptureId: move.captureId, To: move.to})
				if err != nil {
					log.Println("write:", err)
					break
				}
				if rookMove.pieceId != 0 {
					err = c.WriteJSON(JSONRequest{RequestType: "move", PieceId: rookMove.pieceId, CaptureId: 0, To: rookMove.to})
					if err != nil {
						log.Println("rookMove write:", err)
						break
					}
				}
			}

			/*
				if isMove && board.isLegal(&move) {
					fmt.Println("move: ", move)
					_, rookMove := board.move(&board.pieces[move.pieceId], move.toY, move.toX)
					err = c.WriteJSON(JSONRequest{RequestType: "move", PieceId: move.pieceId, CaptureId: move.captureId, ToY: move.toY, ToX: move.toX})
					if err != nil {
						log.Println("write:", err)
						break
					}
					fmt.Println("rookMove: ", rookMove)
					if rookMove.pieceId != 0 {
						err = c.WriteJSON(JSONRequest{RequestType: "move", PieceId: rookMove.pieceId, CaptureId: rookMove.captureId, ToY: rookMove.toY, ToX: rookMove.toX})
						if err != nil {
							log.Println("rookMove write:", err)
							break
						}
					}
				}
				fmt.Println(board.getFen())
			*/

			/*
				log.Printf("recv: %v\n", moveObj)
				if moveObj.PieceId != 0 && board.getPieceisBlack(moveObj.PieceId) == board.isBlack {
					board.fillMove(&moveObj)

						fmt.Printf("move: %v\n", moveObj)
						legal := board.isLegal(&moveObj)
						fmt.Printf("legal: %v\n", legal)

						if legal {
							_, castledMove := board.move(&moveObj)
							fmt.Println("a4: ", board.position[4][0])
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
							fmt.Println("turn: ", board.isBlack)
							fmt.Println("a4: ", board.position[4][0])
						}
					}
			*/
		}
	}))

	log.Fatal(app.Listen(":3000"))

}
