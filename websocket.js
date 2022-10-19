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
    let div = document.createElement("div");
    div.innerHTML = "<b>" + m.timestamp + " </b>" + "<br>" + "<b>" + m.username + ":</b> " + m.text;
    let cself = (myself) ? "self" : "";
    div.className = "msg " + cself;
    document.getElementById("submittedposts").appendChild(div);
  }

  contentHandler(text) {
    const c = JSON.parse(text)
    document.getElementById("content").innerHTML = c.body;
  }

  presenceHandler(text) {
    const m = JSON.parse(text)
    for (let p of m.presences) {
      let user = document.createElement("button");
      user.addEventListener('click', function (event) {
        event.target.id = "presence"
        contentSocket.requestContent(event)
      });
      user.innerHTML = p.username
      user.className = "presence " + p.username
      document.getElementById("presencecontainer").appendChild(user)
    }
  }

  requestContent(e) {
    console.log(e.target.id)
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

  requestPost(e) {
    let m = {
      type: 'post',
      text: e.value,
      timestamp: time(),
      username: "?",
    }
    this.mysocket.send(JSON.stringify(m));
    // document.getElementById("chatIPT").value = ""
  }

  keypress(e) {
    if (e.keyCode == 13) {
      this.wsType = e.target.id.slice(0, -3)
      switch (this.wsType) {
        // case 'post':
        //   this.requestPost()
        //   break;
        case 'chat':
          this.requestChat()
          break;
        // case 'comment':
        //   this.requestComment()
        //   break;
        default:
          console.log("keypress registered for unknown wsType")
          break;
      }
    }
  }

  mouseclick(e) {
    console.log(this.wsType)
    switch (this.wsType) {
      case 'post':
        this.requestPost(e)
        break;
      // case 'chat':
      // this.requestChat()
      // break;
      // case 'comment':
      //   this.requestComment()
      // break;
      default:
        console.log("mouseclick registered for unknown wsType")
        break;
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
    if (URI === 'comment') {
      this.wsType = 'comment'
      console.log("comment Websocket Connected");
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
