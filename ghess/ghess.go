package ghess

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
	"unicode"

	fiber "github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/mustache"
	websocket "github.com/gofiber/websocket/v2"
)

const KING = 'k'
const QUEEN = 'q'
const ROOK = 'r'
const KNIGHT = 'n'
const BISHOP = 'b'
const PAWN = 'p'

const MAX_ENGINE_TIME = 4

const ENGINE1 = "checkCaptureRandom"
const ENGINE2 = "alphaBeta"

const GAME_MODE = "human_vs_engine"

// const GAME_MODE = "human_vs_human"

// normal start
const START_FEN = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"

// KiwiPete
// const START_FEN = "r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1"

// Position 4
// const START_FEN = "r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1"

// const START_FEN = "8/7P/8/8/8/8/2K1k3/8 b  - 8 87"

var materialCountMap = map[rune]int{
	'k': 900,
	'q': 100,
	'r': 50,
	'b': 35,
	'n': 30,
	'p': 10,
}

type Piece struct {
	id          int
	pos         int    // position from 0 to 63
	posB        uint64 // binary position from 0...1 to 1...0
	pieceType   rune   // lower case 'k',...,'p' see KING, ... , PAWN constants
	isBlack     bool
	movementB   uint64 // possible movement of this piece encoded similarly to posB
	moves       [28]int
	numMoves    int
	pinnedMoveB uint64 // is set to all 1 by default but set to 1 only in the direction of pinn if pinned (including capturing pinned piece)
}

func NewPiece(id int, pos int, pieceType rune, isBlack bool) Piece {
	var moves [28]int
	var posB uint64 = 1 << pos
	return Piece{id: id, pos: pos, posB: posB, pieceType: pieceType, isBlack: isBlack, movementB: 0, moves: moves, numMoves: 0, pinnedMoveB: math.MaxUint64}
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
	movesTilEdge       [64][8]int
	check              bool
	doubleCheck        bool
	blockCheckSquaresB uint64
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
		if piece.id == 0 {
			continue
		}
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
		movesTilEdge:       getMovesTilEdge(),
		check:              false,
		doubleCheck:        false,
		blockCheckSquaresB: 0,
	}

	whitePiecePosB := board.combinePositionsOf(whiteIds)
	blackPiecePosB := board.combinePositionsOf(blackIds)
	board.whitePiecePosB = whitePiecePosB
	board.blackPiecePosB = blackPiecePosB
	board.setMovement()
	return board
}

type Move struct {
	pieceId   int
	captureId int
	from      int
	to        int
	promote   int
}

type JSONRequest struct {
	RequestType string `json:"requestType"`
	PieceId     int    `json:"pieceId"`
	CaptureId   int    `json:"captureId"`
	To          int    `json:"to"`
	Promote     int    `json:"promote"` // 0 -> no promotion, 1 -> queen, 2 -> rook, 3 -> bishop, 4 -> knight
}

type JSONSurrounding struct {
	RequestType string     `json:"requestType"`
	Surrounding [8][8]bool `json:"surrounding"`
}

type JSONEnd struct {
	RequestType string `json:"requestType"`
	Msg         string `json:"msg"`
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
			result += `<div id="square_` + strconv.Itoa(id) + `" class="square square_` + color + `" ondrop="onDrop(event);" ondragover="onDragOver(event);" > </div>`
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

func (board *Board) checkGameEnded() (bool, string, string) {
	numMoves := board.GetNumberOfMoves(1)
	if numMoves == 0 {
		if board.check {
			msg := "Checkmate!<br>Good job "
			if board.isBlacksTurn {
				msg += "White!"
			} else {
				msg += "Black!"
			}
			return true, "checkmate", msg
		} else {
			return true, "draw", "Stalemate..."
		}
	}
	// 50 move rule
	if board.halfMoves == 100 {
		return true, "draw", "Draw: Come on you had 50 moves!"
	}
	// draw by insufficent material
	hasEnoughMaterial := false
	numBishop := 0
	numKnight := 0
	for _, pieceId := range board.whiteIds {
		piece := board.pieces[pieceId]
		if piece.posB != 0 {
			if piece.pieceType == PAWN || piece.pieceType == ROOK || piece.pieceType == QUEEN {
				hasEnoughMaterial = true
				break
			} else if piece.pieceType == KNIGHT {
				numKnight += 1
			} else if piece.pieceType == BISHOP {
				numBishop += 1
			}
		}
	}
	if !hasEnoughMaterial && numKnight <= 1 && numBishop <= 1 {
		hasEnoughMaterial = false
		numBishop := 0
		numKnight := 0
		for _, pieceId := range board.blackIds {
			piece := board.pieces[pieceId]
			if piece.posB != 0 {
				if piece.pieceType == PAWN || piece.pieceType == ROOK || piece.pieceType == QUEEN {
					hasEnoughMaterial = true
					break
				} else if piece.pieceType == KNIGHT {
					numKnight += 1
				} else if piece.pieceType == BISHOP {
					numBishop += 1
				}
			}
		}
		if !hasEnoughMaterial && numKnight <= 1 && numBishop <= 1 {
			return true, "draw", "Draw: Not enough material..."
		}
	}

	return false, "", ""
}

func (board *Board) makeEngineMove() (Move, Move) {
	rand.Seed(time.Now().UnixNano())
	engineMove := Move{}
	engine := ENGINE1
	if board.isBlacksTurn {
		engine = ENGINE2
	}

	switch engine {
	case "random":
		engineMove = board.randomEngineMove()
	case "captureRandom":
		engineMove = board.captureEngineMove()
	case "checkCaptureRandom":
		engineMove = board.checkCaptureEngineMove()
	case "alphaBeta":
		engineMove = board.alphaBetaEngineMove()
	}
	// time.Sleep(time.Duration((rand.Intn(3) + 1)) * time.Second)
	// time.Sleep(500 * time.Millisecond)
	engineRookMove := board.Move(&engineMove)
	return engineMove, engineRookMove
}

func (board *Board) makeHumanMove(c *websocket.Conn) (bool, Move, Move) {
	jsonObj := JSONRequest{}
	err := c.ReadJSON(&jsonObj)
	if err != nil {
		log.Println("read:", err)
		return false, Move{}, Move{}
	}
	isMove := false
	needsPromotionType := false
	move := Move{}
	switch jsonObj.RequestType {
	case "movement":
		c.WriteJSON(JSONSurrounding{RequestType: "surrounding", Surrounding: bits2array(board.pieces[jsonObj.PieceId].movementB)})
		// c.WriteJSON(JSONSurrounding{RequestType: "surrounding", Surrounding: bits2array(board.blackPiecePosB)})
	case "move", "capture":
		isMove = true
		move, needsPromotionType = board.NewMove(jsonObj.PieceId, jsonObj.CaptureId, jsonObj.To, jsonObj.Promote)
		if needsPromotionType {
			if board.isLegal(&move) {
				c.WriteJSON(JSONRequest{RequestType: "promotion", PieceId: jsonObj.PieceId, CaptureId: jsonObj.CaptureId, To: jsonObj.To})
			}
			isMove = false
		}
	}

	if isMove && board.isLegal(&move) {
		rookMove := board.Move(&move)

		return true, move, rookMove
	}
	return false, Move{}, Move{}
}

func Run() {
	websocketConns = make(map[int]*websocket.Conn)
	// Create a new engine
	engine := mustache.NewFileSystem(http.Dir("./../ghess/public/templates"), ".mustache")

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Static("/", "./../ghess/public")

	board := GetBoardFromFen(START_FEN)
	// board := GetBoardFromFen("r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1") // Kiwipete
	// board := GetBoardFromFen("r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1") // position 4
	// board := GetBoardFromFen("rnbqkbnr/pppp1ppp/8/4p3/8/5N2/PPPP1PPP/4K3 b KQkq - 0 1")
	// board := GetBoardFromFen("r3k2r/p1ppqpb1/bn2pnp1/3PN3/Pp2P3/2N2Q1p/1PPB1PPP/R3K2R w KQkq a3 0 0")
	// board := GetBoardFromFen("r3k2r/p1ppqpb1/1n2pnp1/1b1PN3/Pp2P3/5Q1p/1PPB1PPP/R3K2R w KQkq - 0 0")

	// board.GetNumberOfMoves(2)
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

		isStarted := false

		var jsonObj JSONRequest
		err := c.ReadJSON(&jsonObj)
		move := Move{}
		rookMove := Move{}
		isMove := false
		playedMoves := []Move{}
		for {
			fmt.Println(board.getFen())
			ended, _, msg := board.checkGameEnded()
			if ended {
				writePGNFile(playedMoves)
				err = c.WriteJSON(JSONEnd{RequestType: "end", Msg: msg})
				if err != nil {
					log.Println("Couldn't send end message:", err)
					break
				}
				break
			}

			if GAME_MODE == "engine_vs_engine" {
				if !isStarted {
					jsonObj = JSONRequest{}
					err = c.ReadJSON(&jsonObj)
					if err != nil {
						log.Println("Couldn't read message:", err)
						break
					}
					if jsonObj.RequestType == "start" {
						isStarted = true
					}
				}

				if !isStarted {
					continue
				} else {
					isStarted = false
				}

				move, rookMove = board.makeEngineMove()
				isMove = true
			} else if GAME_MODE == "human_vs_engine" {
				if board.isBlacksTurn {
					move, rookMove = board.makeEngineMove()
					isMove = true
				} else {
					isMove, move, rookMove = board.makeHumanMove(c)
				}
			} else if GAME_MODE == "human_vs_human" {
				isMove, move, rookMove = board.makeHumanMove(c)
			}
			if isMove {
				playedMoves = append(playedMoves, move)
				err = c.WriteJSON(JSONRequest{RequestType: "move", PieceId: move.pieceId, CaptureId: move.captureId, To: move.to, Promote: move.promote})
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
				isMove = false
			}
		}
	}))

	log.Fatal(app.Listen(":3000"))

}
