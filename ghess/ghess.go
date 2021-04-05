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
	id             int
	position       Position
	pieceType      rune
	isBlack        bool
	onBoard        bool
	vision         [8][8]bool
	movement       [8][8]bool
	moves          [32]Position
	numMoves       int
	updateMovement bool
}

func NewPiece(id int, position Position, pieceType rune, isBlack bool) Piece {
	vision := [8][8]bool{}
	movement := [8][8]bool{}
	moves := [32]Position{}
	return Piece{id: id, position: position, pieceType: pieceType, isBlack: isBlack, onBoard: true, vision: vision, movement: movement, moves: moves, numMoves: 0, updateMovement: false}
}

type Board struct {
	position             [8][8]int
	pieces               [33]Piece // we start with 1 as captureId = 0 should mean no capture
	isBlack              bool
	white_castle_king    bool
	white_castle_queen   bool
	black_castle_king    bool
	black_castle_queen   bool
	en_passant_position  Position // will be 0,0 if not possible
	halfMoves            int
	nextMove             int
	whiteKingId          int
	blackKingId          int
	lastMoveWasCheck     bool
	currentlyReverseMode bool
}

type BoardPrimitives struct {
	isBlack             bool
	white_castle_king   bool
	white_castle_queen  bool
	black_castle_king   bool
	black_castle_queen  bool
	en_passant_position Position // will be 0,0 if not possible
	halfMoves           int
	nextMove            int
	whiteKingId         int
	blackKingId         int
	lastMoveWasCheck    bool
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
	fromY     int
	fromX     int
	toY       int
	toX       int
}

func (board *Board) newMove(pieceId int, captureId int, toY int, toX int, isCapture bool) Move {
	fromY := board.pieces[pieceId].position.y
	fromX := board.pieces[pieceId].position.x
	if isCapture {
		toX = board.pieces[captureId].position.x
		toY = board.pieces[captureId].position.y
	} else if captureId == 0 {
		if board.position[toY][toX] != 0 { // fill capture if there is a piece on that position
			captureId = board.position[toY][toX]
		} else if toY == board.en_passant_position.y && toX == board.en_passant_position.x && board.pieces[pieceId].pieceType == PAWN {
			if board.isBlack {
				captureId = board.position[toY-1][toX]
			} else {
				captureId = board.position[toY+1][toX]
			}
		}
	}
	return Move{pieceId: pieceId, captureId: captureId, fromY: fromY, fromX: fromX, toY: toY, toX: toX}
}

type JSONRequest struct {
	RequestType string `json:"requestType"`
	PieceId     int    `json:"pieceId"`
	CaptureId   int    `json:"captureId"`
	ToY         int    `json:"toY"`
	ToX         int    `json:"toX"`
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
			pieceId = board.position[y][x]
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
	if board.isBlack {
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
	if board.en_passant_position.x != 0 {
		fen += string(rune('a' + board.en_passant_position.x))
		fen += strconv.Itoa(8 - board.en_passant_position.y)
	} else {
		fen += "-"
	}

	fen += " "
	fen += strconv.Itoa(board.halfMoves)
	fen += " "
	fen += strconv.Itoa(board.nextMove)

	return fen
}

func (board *Board) moveToToLongAlgebraic(fromY, fromX, toY, toX int) string {
	res := string(rune('a' + fromX))
	res += strconv.Itoa(8 - fromY)
	res += "-" + string(rune('a'+toX))
	res += strconv.Itoa(8 - toY)
	return res
}

func displayFen(fen string) string {
	board := GetBoardFromFen(fen)
	return board.display()
}

func (board *Board) display() string {
	result := displayGround()
	for pieceId, piece := range board.pieces {
		if !piece.onBoard {
			continue
		}
		pieceName := getPieceName(&board.pieces[pieceId])
		position := piece.position
		left := strconv.Itoa(position.x * 10)
		top := strconv.Itoa(position.y * 10)
		result += `<div class="piece" draggable="true" onclick="onClick(event);" ondragstart="onDragStart(event);" ondragend="onDragEnd(event);" ondrop="onDrop(event);" ondragover="onDragOver(event);"> 
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
			left := strconv.Itoa(j * 10)
			top := strconv.Itoa(i * 10)
			result += `<div id="square_` + strconv.Itoa(i) + `_` + strconv.Itoa(j) + `" class="square square_` + color + `"> </div>`
			result += `<div id="square_` + strconv.Itoa(i) + `_` + strconv.Itoa(j) + `_overlay" class="square_overlay" ondrop="onDrop(event);" ondragover="onDragOver(event);" style="left: ` + left + `vmin; top: ` + top + `vmin;"> </div>`
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
	var pieces [33]Piece
	pieceId := 1
	whiteKingId := 0
	blackKingId := 0
	for r, row := range rows {
		cpos := 0
		for _, p := range row {
			if (p > 'a' && p < 'z') || (p > 'A' && p < 'Z') {
				position[r][cpos] = pieceId
				pieces[pieceId] = NewPiece(pieceId, Position{x: cpos, y: r}, unicode.ToLower(p), p == unicode.ToLower(p))
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
	isBlack := false
	if parts[1][0] == 'b' {
		isBlack = true
	}
	en_passant_position := Position{x: 0, y: 0}
	if parts[3][0] != '-' {
		en_passant_position.y = 8 - int(parts[3][1]-'0')
		en_passant_position.x = int(parts[3][0] - 'a')
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
	board := Board{
		position:             position,
		pieces:               pieces,
		isBlack:              isBlack,
		white_castle_king:    strings.ContainsRune(parts[2], 'K'),
		white_castle_queen:   strings.ContainsRune(parts[2], 'Q'),
		black_castle_king:    strings.ContainsRune(parts[2], 'k'),
		black_castle_queen:   strings.ContainsRune(parts[2], 'q'),
		en_passant_position:  en_passant_position,
		halfMoves:            halfMoves,
		nextMove:             nextMove,
		whiteKingId:          whiteKingId,
		blackKingId:          blackKingId,
		lastMoveWasCheck:     false,
		currentlyReverseMode: false,
	}
	board.updateVision()
	board.updateMovement()
	return board
}

func (board *Board) fillMove(m *JSONMove) {
	if m.ToX == 0 {
		captureId := m.CaptureId
		m.ToX = board.pieces[captureId].position.x
		m.ToY = board.pieces[captureId].position.y
	} else if m.CaptureId == 0 && board.position[m.ToY-1][m.ToX-1] != 0 { // fill capture if there is a piece on that position
		fmt.Printf("fill capture as there is %d on y,x: %d, %d\n", board.position[m.ToY-1][m.ToX-1], m.ToY, m.ToX)
		m.CaptureId = board.position[m.ToY-1][m.ToX-1]
	}
}

func (board *Board) updateCastlePrivileges() {
	for _, piece := range board.pieces {
		if piece.pieceType == KING {
			if piece.isBlack && (piece.position.x != 4 || piece.position.y != 0) {
				board.black_castle_king = false
				board.black_castle_queen = false
			} else if !piece.isBlack && (piece.position.x != 4 || piece.position.y != 7) {
				board.white_castle_king = false
				board.white_castle_queen = false
			}
			// if rook moved or captured
			if piece.isBlack {
				if board.pieces[board.position[0][7]].pieceType != ROOK || !board.pieces[board.position[0][7]].isBlack {
					board.black_castle_king = false
				}
				if board.pieces[board.position[0][0]].pieceType != ROOK || !board.pieces[board.position[0][0]].isBlack {
					board.black_castle_queen = false
				}
			}
			if !piece.isBlack {
				if board.pieces[board.position[7][7]].pieceType != ROOK || board.pieces[board.position[7][7]].isBlack {
					board.white_castle_king = false
				}
				if board.pieces[board.position[7][0]].pieceType != ROOK || board.pieces[board.position[7][0]].isBlack {
					board.white_castle_queen = false
				}
			}
		}
	}
}

/*
func (board *Board) move(m *JSONMove) (int, JSONMove) {
	capturedId, castledMove := board.moveTemp(m)
	fmt.Println("a4 inside move: ", board.position[4][0])
	// only count once for castling
	if board.pieces[m.PieceId].pieceType == PAWN || m.CaptureId != 0 {
		board.halfMoves = 0
	} else {
		board.halfMoves += 1
	}
	if !board.isBlack {
		board.nextMove += 1
	}
	board.isBlack = !board.isBlack
	return capturedId, castledMove
}
*/

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
	board.pieces[m.PieceId].position.x = m.ToX
	board.pieces[m.PieceId].position.y = m.ToY
	if m.CaptureId != 0 {
		board.pieces[m.CaptureId].onBoard = false
	}

	// disallow castling
	if piece.pieceType == KING {
		if !piece.isBlack {
			board.white_castle_king = false
			board.white_castle_queen = false
		} else {
			board.black_castle_king = false
			board.black_castle_queen = false
		}
	}

	if piece.pieceType == ROOK {
		if !piece.isBlack {
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
	if piece.pieceType == PAWN {
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

/*
func (board *Board) reverseMove(m *JSONMove, fromY int, fromX int, capturedId int, castledMove *JSONMove, boardPrimitives BoardPrimitives) {
	fmt.Println("reverse a move")
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
*/

func engineMove() JSONMove {
	return JSONMove{PieceId: 13, CaptureId: 0, ToY: 4, ToX: 5}
}

func (board *Board) getPieceisBlack(piece int) bool {
	return board.pieces[piece].isBlack
}

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
	move = board.newMove(pieceId, 0, toY, toX, false)
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

func Run() {
	websocketConns = make(map[int]*websocket.Conn)
	// Create a new engine
	engine := mustache.NewFileSystem(http.Dir("./../ghess/public/templates"), ".mustache")

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Static("/", "./../ghess/public")

	board := GetBoardFromFen("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
	// board := GetBoardFromFen("rnbqkbnr/pppp1ppp/8/4p3/8/5N2/PPPP1PPP/4K3 b KQkq - 0 1")
	// board := GetBoardFromFen("r3k2r/p1ppqpb1/bn2pnp1/3PN3/Pp2P3/2N2Q1p/1PPB1PPP/R3K2R w KQkq a3 0 0")
	// board := GetBoardFromFen("r3k2r/p1ppqpb1/1n2pnp1/1b1PN3/Pp2P3/5Q1p/1PPB1PPP/R3K2R w KQkq - 0 0")

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
			case "vision":
				c.WriteJSON(JSONSurrounding{RequestType: "surrounding", Surrounding: board.pieces[jsonObj.PieceId].vision})
			case "movement":
				c.WriteJSON(JSONSurrounding{RequestType: "surrounding", Surrounding: board.pieces[jsonObj.PieceId].movement})
			case "move":
				isMove = true
				move = board.newMove(jsonObj.PieceId, jsonObj.CaptureId, jsonObj.ToY, jsonObj.ToX, false)
			case "capture":
				isMove = true
				move = board.newMove(jsonObj.PieceId, jsonObj.CaptureId, jsonObj.ToY, jsonObj.ToX, true)
			}
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
