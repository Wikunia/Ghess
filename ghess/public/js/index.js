var socket = new WebSocket("ws://localhost:3000/ws")

var PieceId = 0;
var captureId = 0;
var to = 0;

var endDialog = document.getElementById("endDialog")

socket.onopen = function(e) {
    socket.send(JSON.stringify({"NewConnection": true}));
};

socket.onmessage = function (event) {
    let jsonObj = JSON.parse(event.data)
    if (jsonObj.requestType == "surrounding") {
        showSurrounding(jsonObj)
    } else if (jsonObj.requestType == "move") {
        resetSurrounding()
        movePiece(jsonObj)
    } else if (jsonObj.requestType == "promotion") {
        let promotionDialog = document.getElementById("promotion");
        PieceId = jsonObj.PieceId
        captureId = jsonObj.captureId
        to = jsonObj.to
        promotionDialog.showModal();
    } else if (jsonObj.requestType == "end") {
        endDialog.innerHTML = jsonObj.msg;
        endDialog.showModal();
    }
}

function resetSurrounding() {
    for (let i = 0; i < 8; i++) {
        for (let j = 0; j < 8; j++) {
            let p = i*8 + j
            field = document.querySelector("#square_"+p+"_overlay") 
            field.style.display = "none"
        }
        
    }
}

function showSurrounding(jsonObj) {
    resetSurrounding()
    let surrounding = jsonObj.surrounding
    for (let i = 0; i < 8; i++) {
        for (let j = 0; j < 8; j++) {
            if (surrounding[i][j]) {
                let p = i*8 + j
                field = document.querySelector("#square_"+p+"_overlay") 
                field.style.display = "block"
            } 
        }
    }
}
  
function movePiece(move) {
    let piece = document.querySelector("#piece_"+move.PieceId);
    if (move.captureId !== 0) {
        let capturedPiece = document.querySelector("#piece_"+move.captureId);
        capturedPiece.style.display = "none";
    } 
    // can move to this position if en passant
    var y = Math.floor(move.to/8);
    var x = move.to % 8;

    field = document.querySelector("#square_"+move.to+"_overlay") 
    field.style.display = "block"

    piece.style.left = (x*10)+"vmin";
    piece.style.top = (y*10)+"vmin";
    if (move.promote != 0) {
        // 1 queen, 2 rook, 3 bishop, 4 knight
        let current_source = piece.src
        let parts = current_source.split("/")
        color = parts[parts.length-1].split("_")[0]
        switch (move.promote) {
            case 1:
                piece.src = "images/"+color+"_queen.png"
                break;
            case 2:
                piece.src = "images/"+color+"_rook.png"
                break;
            case 3:
                piece.src = "images/"+color+"_bishop.png"
                break;
            case 4:
                piece.src = "images/"+color+"_knight.png"
                break;
        }
    }
}

function onDragStart(event) {
    event
      .dataTransfer
      .setData('text/plain', event.target.id);
  
    // send click event first
    let PieceId = parseInt(event.target.id.split("_")[1])
    var data = JSON.stringify({"requestType": "movement", "PieceId": PieceId});
    socket.send(data)

    event
      .target
      .style
      .opacity = '0.0';
  }

function onDragEnd(event) {
    event.target
    .style
    .opacity = '1.0';
}

function onDragOver(event) {
    event.preventDefault();
}


function onDrop(event) {
    const id = event
      .dataTransfer
      .getData('text');

    let PieceId = id.split("_")[1]

    let pieceOrField = event.target.id.split("_")[0]
    console.log("event.target.id: ", event.target.id)

    if (pieceOrField == "piece") {
        sendCapturePiece(PieceId, event.target.id)
    } else if (pieceOrField == "square") {
        sendMovePiece(PieceId, event.target.id)
    }
}

function sendMovePiece(PieceId, toStr) {
    let [_,to] = toStr.split("_")
    console.log("move: ", PieceId, " -> ", to)

    var data = JSON.stringify({"requestType": "move", "PieceId": parseInt(PieceId), "to": parseInt(to)});
    socket.send(data);
}

function sendCapturePiece(PieceId, toStr) {
    let [_,captureId] = toStr.split("_")
    console.log("capture: ", PieceId, " -> ", captureId)

    var data = JSON.stringify({"requestType": "capture", "PieceId": parseInt(PieceId), "captureId": parseInt(captureId)});
    socket.send(data);
}

function onClick(event) {
    let PieceId = parseInt(event.target.id.split("_")[1])
    var data = JSON.stringify({"requestType": "movement", "PieceId": PieceId});
    socket.send(data)
}

let promotionQueen = document.getElementById("promotionQueen");
promotionQueen.addEventListener('click', function() {
    var data = JSON.stringify({"requestType": "move", "PieceId": PieceId, "to": to, "captureId": captureId,  "promote": 1});
    socket.send(data);
})

let promotionRook = document.getElementById("promotionRook");
promotionRook.addEventListener('click', function() {
    var data = JSON.stringify({"requestType": "move", "PieceId": PieceId, "to": to, "captureId": captureId, "promote": 2});
    socket.send(data);
})

let promotionBishop = document.getElementById("promotionBishop");
promotionBishop.addEventListener('click', function() {
    var data = JSON.stringify({"requestType": "move", "PieceId": PieceId, "to": to, "captureId": captureId, "promote": 3});
    socket.send(data);
})

let promotionKnight = document.getElementById("promotionKnight");
promotionKnight.addEventListener('click', function() {
    var data = JSON.stringify({"requestType": "move", "PieceId": PieceId, "to": to, "captureId": captureId, "promote": 4});
    socket.send(data);
})

let startButton = document.getElementById("start");
startButton.addEventListener('click', function() {
    var data = JSON.stringify({"requestType": "start"});
    socket.send(data);
})