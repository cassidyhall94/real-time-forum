class MySocket {
    constructor() {
        this.mysocket = null;
    }

    chatHandler(text, myself) {
        var div = document.createElement("div");
        div.innerHTML = text;
        var cself = (myself) ? "self" : "";
        div.className = "msg " + cself;
        document.getElementById("msgcontainer").appendChild(div);
    }

    send() {
        var txt = document.getElementById("ipt").value;
        this.chatHandler("<b>You:</b> " + txt, true);
        this.mysocket.send(txt);
        document.getElementById("ipt").value = ""
    }

    keypress(e) {
        if (e.keyCode == 13) {
            this.send();
        }
    }

    contentHandler(text) {
        document.getElementById("content").innerHTML = text;
    }

    connectSocket(URI, handler) {
        console.log("websocket connected");
        var socket = new WebSocket("ws://localhost:8080/" + URI); //make sure the port matches with your golang code
        this.mysocket = socket;

        socket.onmessage = (e) => {
            handler(e.data, false);
        }
        socket.onopen = () => {
            console.log("socket opened")
        };
        socket.onclose = () => {
            console.log("socket closed")
        }
    }
}