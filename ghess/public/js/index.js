var socket = new WebSocket("ws://localhost:3000/ws")

socket.onopen = function(e) {
    socket.send(JSON.stringify({"NewConnection": true}));
};

socket.onmessage = function (event) {
    move = JSON.parse(event.data)
    if (move.pieceId !== undefined) {
        movePiece(move)
    }
}
  
function movePiece(move) {
    let piece = document.querySelector("#piece_"+move.pieceId);
    if (move.captureId !== 0) {
        let capturedPiece = document.querySelector("#piece_"+move.captureId);
        capturedPiece.style.display = "none";

        piece.style.left = capturedPiece.style.left;
        piece.style.top = capturedPiece.style.top;
    } else {
        piece.style.left = ((move.toX-1)*10)+"vmin";
        piece.style.top = ((move.toY-1)*10)+"vmin";
    }
}

function onDragStart(event) {
    console.log("event", event)
    event
      .dataTransfer
      .setData('text/plain', event.target.id);
  
    event
      .currentTarget
      .style
      .backgroundColor = 'yellow';
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