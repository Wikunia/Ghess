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
        capturePiece(pieceId, event.target.id)
    } else if (pieceOrField == "square") {
        movePiece(pieceId, event.target.id)
    }
}

function movePiece(pieceId, to) {
    let [_,to_y,to_x] = to.split("_")
    console.log(pieceId, " -> ", to_y, " ", to_x)

    var xhr = new XMLHttpRequest();
    var url = "/api/move";
    xhr.open("POST", url, true);
    xhr.setRequestHeader("Content-Type", "application/json");
    xhr.onreadystatechange = function () {
        if (xhr.readyState === 4 && xhr.status === 200) {
            console.log(xhr.responseText)
            if (xhr.responseText == "success") {
                let piece = document.querySelector("#piece_"+pieceId);
                piece.style.left = (parseInt(to_x)*10)+"vmin";
                piece.style.top = (parseInt(to_y)*10)+"vmin";
            }
        }
    };
    var data = JSON.stringify({"pieceId": parseInt(pieceId), "to_y": parseInt(to_y), "to_x": parseInt(to_x)});
    xhr.send(data);
}

function capturePiece(pieceId, to) {
    let [_,captureId] = to.split("_")
    console.log(pieceId, " -> ", captureId)

    var xhr = new XMLHttpRequest();
    var url = "/api/capture";
    xhr.open("POST", url, true);
    xhr.setRequestHeader("Content-Type", "application/json");
    xhr.onreadystatechange = function () {
        if (xhr.readyState === 4 && xhr.status === 200) {
            console.log(xhr.responseText)
            if (xhr.responseText == "success") {
                let capturedPiece = document.querySelector("#piece_"+captureId);
                capturedPiece.style.display = "none";

                let piece = document.querySelector("#piece_"+pieceId);
                piece.style.left = capturedPiece.style.left;
                piece.style.top = capturedPiece.style.top;
            }
        }
    };
    var data = JSON.stringify({"pieceId": parseInt(pieceId), "captureId": parseInt(captureId)});
    xhr.send(data);
}