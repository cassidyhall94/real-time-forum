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
        let time = new Date().toLocaleString();
        let line = "<b>" + time + " </b>" + "<br>" + "<b>You:</b> " + txt
        this.chatHandler(line, true);
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

    requestContent(text) {
        this.mysocket.send(text);
    }

    connectSocket(URI, handler) {
        if (URI === 'chat') {
            console.log("Chat Websocket Connected");
        }
        if (URI === 'content') {
            console.log("Content Websocket Connected");
        }
        var socket = new WebSocket("ws://localhost:8080/" + URI);
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
