var socket = new WebSocket("ws://localhost:3000/ws")

socket.onopen = function(e) {
    socket.send(JSON.stringify({"NewConnection": true}));
};

socket.onmessage = function (event) {
    let jsonObj = JSON.parse(event.data)
    if (jsonObj.requestType == "vision") {
        showVision(jsonObj)
    } else if (jsonObj.requestType == "move") {
        movePiece(jsonObj)
    }
}

function resetVision() {
    for (let i = 0; i < 8; i++) {
        for (let j = 0; j < 8; j++) {
            field = document.querySelector("#square_"+i+"_"+j+"_overlay") 
            field.style.display = "none"
        }
        
    }
}

function showVision(jsonObj) {
    resetVision()
    let vision = jsonObj.vision
    for (let i = 0; i < 8; i++) {
        for (let j = 0; j < 8; j++) {
            if (vision[i][j]) {
                console.log("#square_"+i+"_"+j)
                field = document.querySelector("#square_"+i+"_"+j+"_overlay") 
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

        piece.style.left = capturedPiece.style.left;
        piece.style.top = capturedPiece.style.top;
    } 
    // can move to this position if en passant
    if (move.toX != 0) {
        piece.style.left = (move.toX*10)+"vmin";
        piece.style.top = (move.toY*10)+"vmin";
    }
}

function onDragStart(event) {
    console.log("event", event)
    event
      .dataTransfer
      .setData('text/plain', event.target.id);
  
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

function sendMovePiece(pieceId, to) {
    let [_,to_y,to_x] = to.split("_")
    console.log(pieceId, " -> ", to_y, " ", to_x)

    var data = JSON.stringify({"pieceId": parseInt(pieceId), "toY": parseInt(to_y), "toX": parseInt(to_x)});
    socket.send(data);
}

function sendCapturePiece(pieceId, to) {
    let [_,captureId] = to.split("_")
    console.log(pieceId, " -> ", captureId)

    var data = JSON.stringify({"pieceId": parseInt(pieceId), "captureId": parseInt(captureId)});
    socket.send(data);
}

function onClick(event) {
    let pieceId = parseInt(event.target.id.split("_")[1])
    var data = JSON.stringify({"requestType": "vision", "pieceId": pieceId});
    socket.send(data)
}