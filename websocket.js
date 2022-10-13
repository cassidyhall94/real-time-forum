let presentUsers = []

class MySocket {
  wsType = ""

  constructor() {
    this.mysocket = null;
  }

  chatHandler(text, myself) {
    const m = JSON.parse(text)
    if (m.hasOwnProperty('presenceList')) {
      presentUsers = m.PresenceList
    }
    var div = document.createElement("div");
    let msgContainer = document.getElementById('chatIPT')
    div.innerHTML = "<b>" + m.Timestamp + " </b>" + "<br>" + "<b>" + m.Username + ":</b> " + m.Text;
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
    document.getElementById("content").innerHTML = text;
  }

  requestContent(e) {
    // console.log(e.target.id)
    // console.log(this)
    this.mysocket.send(JSON.stringify({
      type: 'content',
      username: "?",
      resource: e.target.id,
    }));
  }

  send() {
    let time = new Date().toLocaleString();
    let txt
    if (this.wsType === 'chat') {
      txt = document.getElementById("chatIPT").value;
      let line = "<b>" + time + " </b>" + "<br>" + "<b>You:</b> " + txt
      let m = {
        type: 'chat',
        text: txt,
        timestamp: time,
        username: "?",
      }
      this.mysocket.send(JSON.stringify(m));
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
