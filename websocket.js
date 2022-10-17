const time = () => { new Date().toLocaleString() };

class MySocket {
  wsType = ""

  constructor() {
    this.mysocket = null;
  }

  chatHandler(text, myself) {
    const m = JSON.parse(text)
    let div = document.createElement("div");
    let msgContainer = document.getElementById('chatIPT')
    div.innerHTML = "<b>" + m.timestamp + " </b>" + "<br>" + "<b>" + m.username + ":</b> " + m.text;
    let cself = (myself) ? "self" : "";
    div.className = "msg " + cself;
    document.getElementById("msgcontainer").appendChild(div);
    div.after(msgContainer)
  }

  postHandler(text, myself) {
    let post = document.createElement("div");
    let postContainer = document.getElementById('postIPT')
    post.innerHTML = text;
    var cself = (myself) ? "self" : "";
    post.className = "post " + cself;
    document.getElementById("postcontainer").appendChild(post);
    post.after(postContainer)
  }

  contentHandler(text) {
    const c = JSON.parse(text)
    document.getElementById("content").innerHTML = c.body;
  }

  presenceHandler(text) {
    const m = JSON.parse(text)
    for (let p of m.presences) {
      
      let user = document.createElement("button");
      user.click = contentSocket.requestContent
      user.innerHTML = p.username
      user.className = "presence " + p.username
      document.getElementById("presencecontainer").appendChild(user)
    }
  }

  requestContent(e) {
    this.mysocket.send(JSON.stringify({
      type: "content",
      username: "?",
      resource: e.target.id,
    }));
  }

  requestChat() {
    let m = {
      type: 'chat',
      text: document.getElementById("chatIPT").value,
      timestamp: time(),
      username: "?",
    }
    this.mysocket.send(JSON.stringify(m));
    document.getElementById("chatIPT").value = ""
  }

  requestPost() {
    let txt = document.getElementById("postIPT").value;
    let line = "<b>" + time + " </b>" + "<br>" + "<b>You:</b> " + txt
    this.postHandler(line, true);
    this.mysocket.send(txt);
    document.getElementById("postIPT").value = ""
  }

  keypress(e) {
    if (e.keyCode == 13) {
      this.wsType = e.target.id.slice(0, -3)
      switch (this.wsType) {
        case 'post':
          this.requestPost()
          break;
        case 'chat':
          this.requestChat()
          break;
        default:
          console.log("keypress registered for unknown wsType")
          break;
      }
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
    if (URI === 'presence') {
      this.wsType = 'presence'
      console.log("Presence Websocket Connected");
    }
    var socket = new WebSocket("ws://localhost:8080/" + URI);
    this.mysocket = socket;

    socket.onmessage = (e) => {
      console.log("socket message")
      handler(e.data, false);
    };
    socket.onopen = () => {
      console.log("socket opened");
    };
    socket.onclose = () => {
      console.log("socket closed");
    };
  }

  getRegistrationDetails() {
    //AJAX html request
    httpRequest = new XMLHttpRequest();
    if (!httpRequest) {
      console.log("Giving up :( Cannot create an XMLHTTP instance'");
    }
    url = "ws://localhost:8080/";
    httpRequest.onreadystatechange = sendContents;
    httpRequest.open("POST", url);
    httpRequest.setRequestHeader(
      "Content-type",
      "application/x-www-form-urlencoded"
    );
    var fd = new FormData();
    fd.set("username", document.getElementById("reg-username").value);
    fd.set("email", document.getElementById("reg-email").value);
    fd.set("password", document.getElementsByClassName("reg-password").value);

    httpRequest.send(fd);

    function sendContents() {
      if (httpRequest.readyState === 4) {
        if (httpRequest.Status === 200) {
          alert(httpRequest.responseText);
        }
      }
    }
  }
}
