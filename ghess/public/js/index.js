var socket = new WebSocket("ws://localhost:3000/ws")

var pieceId = 0;
var captureId = 0;
var to = 0;

socket.onopen = function(e) {
    socket.send(JSON.stringify({"NewConnection": true}));
};

socket.onmessage = function (event) {
    let jsonObj = JSON.parse(event.data)
    if (jsonObj.requestType == "surrounding") {
        showSurrounding(jsonObj)
    } else if (jsonObj.requestType == "move") {
        movePiece(jsonObj)
        resetSurrounding()
    } else if (jsonObj.requestType == "promotion") {
        let promotionDialog = document.getElementById("promotion");
        pieceId = jsonObj.pieceId
        captureId = jsonObj.captureId
        to = jsonObj.to
        promotionDialog.showModal();
    } else if (jsonObj.requestType == "end") {
        alert(jsonObj.type)
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
    let piece = document.querySelector("#piece_"+move.pieceId);
    if (move.captureId !== 0) {
        let capturedPiece = document.querySelector("#piece_"+move.captureId);
        capturedPiece.style.display = "none";
    } 
    // can move to this position if en passant
    var y = Math.floor(move.to/8);
    var x = move.to % 8;
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
    let pieceId = parseInt(event.target.id.split("_")[1])
    var data = JSON.stringify({"requestType": "movement", "pieceId": pieceId});
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

    let pieceId = id.split("_")[1]

    let pieceOrField = event.target.id.split("_")[0]

    if (pieceOrField == "piece") {
        sendCapturePiece(pieceId, event.target.id)
    } else if (pieceOrField == "square") {
        sendMovePiece(pieceId, event.target.id)
    }
}

function sendMovePiece(pieceId, toStr) {
    let [_,to] = toStr.split("_")
    console.log(pieceId, " -> ", to)

    var data = JSON.stringify({"requestType": "move", "pieceId": parseInt(pieceId), "to": parseInt(to)});
    socket.send(data);
}

function sendCapturePiece(pieceId, toStr) {
    let [_,captureId] = toStr.split("_")
    console.log(pieceId, " -> ", captureId)

    var data = JSON.stringify({"requestType": "capture", "pieceId": parseInt(pieceId), "captureId": parseInt(captureId)});
    socket.send(data);
}

function onClick(event) {
    let pieceId = parseInt(event.target.id.split("_")[1])
    var data = JSON.stringify({"requestType": "movement", "pieceId": pieceId});
    socket.send(data)
}

let promotionQueen = document.getElementById("promotionQueen");
promotionQueen.addEventListener('click', function() {
    var data = JSON.stringify({"requestType": "move", "pieceId": pieceId, "to": to, "promote": 1});
    socket.send(data);
})

let promotionRook = document.getElementById("promotionRook");
promotionRook.addEventListener('click', function() {
    var data = JSON.stringify({"requestType": "move", "pieceId": pieceId, "to": to, "promote": 2});
    socket.send(data);
})

let promotionBishop = document.getElementById("promotionBishop");
promotionBishop.addEventListener('click', function() {
    var data = JSON.stringify({"requestType": "move", "pieceId": pieceId, "to": to, "promote": 3});
    socket.send(data);
})

let promotionKnight = document.getElementById("promotionKnight");
promotionKnight.addEventListener('click', function() {
    var data = JSON.stringify({"requestType": "move", "pieceId": pieceId, "to": to, "promote": 4});
    socket.send(data);
})