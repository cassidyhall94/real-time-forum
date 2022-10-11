

class MySocket {
    wsType = ""
    constructor() {
        this.mysocket = null;
    }

    chatHandler(text, myself) {
        var div = document.createElement("div");
        let msgContainer = document.getElementById('chatIPT')
        div.innerHTML = text;
        var cself = (myself) ? "self" : "";
        div.className = "msg " + cself;
        document.getElementById("msgcontainer").appendChild(div);
        div.after(msgContainer)
    }

    postHandler(text, myself) {
        var post = document.createElement("div");
        let postContainer = document.getElementById('postIPT')
        post.innerHTML = text;
        var cself = (myself) ? "self" : "";
        post.className = "post " + cself;
        document.getElementById("postcontainer").appendChild(post);
        post.after(postContainer)
    }

    contentHandler(text) {
        console.log('----------------------')
        console.log(text)
        console.log('----------------------')
        document.getElementById("content").innerHTML = text;
    }

    requestContent(e) {
        console.log(e.target.id)
        console.log(this)
        this.mysocket.send(e.target.id);
    }

    send() {
        let time = new Date().toLocaleString();
        let txt
        if (this.wsType === 'chat') {
            txt = document.getElementById("chatIPT").value;
            let line = "<b>" + time + " </b>" + "<br>" + "<b>You:</b> " + txt
            this.chatHandler(line, true);
            this.mysocket.send(txt);
            document.getElementById("chatIPT").value = ""
        }
        if (this.wsType === 'post') {
            txt = document.getElementById("postIPT").value;
            let line = "<b>" + time + " </b>" + "<br>" + "<b>You:</b> " + txt
            this.postHandler(line, true);
            this.mysocket.send(txt);
            document.getElementById("postIPT").value = ""
        }
    }

    keypress(e) {
        if (e.keyCode == 13) {
            this.wsType = e.target.id.slice(0, -3)
            this.send();
        }
    }

    connectSocket(URI, handler) {
        if (URI === 'chat') {
            this.wsType = 'chat'
            console.log("Chat Websocket Connected");
        }
        if (URI === 'content') {
            this.wsType = 'content'
            console.log("Content Websocket Connected");
        }
        if (URI === 'post') {
            this.wsType = 'post'
            console.log("Post Websocket Connected");
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
