//TODO: fix const time as it is not formatted correctly and add where time/date is needed
const time = () => { return new Date().toLocaleString() };

class MySocket {
  wsType = ""

  constructor() {
    this.mysocket = null;
  }

  // TODO: insert user ID variable, participants needs to be filled
  sendNewChatRequest() {
    console.log("new chat request")
    let m = {
      type: 'chat',
      timestamp: time(),
      conversations: [
        {
          participants: [
            //sender: bar userID
            {
              id: "975496ca-9bfc-4d71-8736-da4b6383a575",
            },
            //other participants (receiver): foo userID
            {
              id: "6d01e668-2642-4e55-af73-46f057b731f9",
            }
          ],
          chats: [
            {
              sender: {
                // TODO: this is just the first placeholder above, once the user is logged in and their ID is stored client side this ID should represent the logged in user
                // bar userID
                id: "975496ca-9bfc-4d71-8736-da4b6383a575",
              },
              body: document.getElementById('chatIPT').value,
            }
          ]
        }
      ]
    }
    this.mysocket.send(JSON.stringify(m));
    document.getElementById('chatIPT').value = ""
  }

  keypress(e) {
    if (e.keyCode == 13) {
      this.wsType = e.target.id.slice(0, -3)
      if (this.wsType = 'chat') {
        this.sendNewChatRequest()
      }
    }
  }

  chatHandler(text) {
    const m = JSON.parse(text)
    for (let c of m.conversations) {
      for (let p of c.chats) {
        let chat = document.createElement("div");
        chat.className = "submittedchat"
        chat.id = p.chat_id
        chat.innerHTML = "<b>Me: " + p.sender.nickname + "</b>" + "<br>" + "<b>Date: " + "</b>" + p.date + "<br>" + p.body + "<br>";
        document.getElementById("chatcontainer").appendChild(chat)
      }
    }
  }

  contentHandler(text) {
    const c = JSON.parse(text)
    document.getElementById("content").innerHTML = c.body;
  }

  presenceHandler(text) {
    const m = JSON.parse(text)
    for (let p of m.presences) {
      const consp = p
      let user = document.createElement("button");
      user.addEventListener('click', function (event, chat = consp) {
        event.target.id = "chat"
        contentSocket.sendChatContentRequest(event, chat.chat_id)
      });
      user.id = p.id
      user.innerHTML = p.nickname
      user.style.color = 'white'
      user.className = "presence " + p.nickname
      document.getElementById("presencecontainer").appendChild(user)
    }
  }

  postHandler(text) {
    const m = JSON.parse(text)
    for (let p of m.posts) {
      const consp = p
      let post = document.createElement("div");
      post.className = "submittedpost"
      post.id = p.post_id
      post.innerHTML = "<b>Title: " + p.title + "</b>" + "<br>" + "<b>Nickname: " + "</b>" + p.nickname + "<br>" + "<b>Category: " + p.categories + "</b>" + "<br>" + p.body + "<br>";
      let button = document.createElement("button")
      button.classname = "addcomment"
      button.innerHTML = "Comments"
      button.addEventListener('click', function (event, post = consp) {
        event.target.id = "comment"
        contentSocket.sendContentRequest(event, post.post_id)
      });
      post.appendChild(button)
      console.log(document.getElementById("submittedposts"))
      document.getElementById("submittedposts").appendChild(post)
    }
  }

  sendNewCommentRequest(e) {
    let post = document.getElementById('postcontainerforcomments')
    for (const child of post.children) {
      if (containsNumber(child.id)) {
        let m = {
          type: 'post',
          timestamp: time(),
          posts: [
            {
              post_id: child.id,
              comments: [
                {
                  post_id: child.id,
                  body: document.getElementById('commentbody').value,
                }
              ]
            }
          ]
        }
        this.mysocket.send(JSON.stringify(m));
        document.getElementById('commentbody').value = ""
      }
    }
  }

  sendNewPostRequest(e) {
    let m = {
      type: 'post',
      timestamp: time(),
      posts: [
        {
          nickname: e.target.nickname,
          title: document.getElementById('posttitle').value,
          categories: document.getElementById('category').value,
          body: document.getElementById('postbody').value,
        }
      ]
    }
    this.mysocket.send(JSON.stringify(m));
    document.getElementById('posttitle').value = ""
    document.getElementById('category').value = ""
    document.getElementById('postbody').value = ""
  }

  sendSubmittedPostsRequest() {
    this.mysocket.send(JSON.stringify({
      type: "post",
    }));
  }

  sendContentRequest(e, post_id = "") {
    this.mysocket.send(JSON.stringify({
      type: "content",
      resource: e.target.id,
      post_id: post_id,
    }));
  }

  sendChatContentRequest(e, chat_id = "") {
    this.mysocket.send(JSON.stringify({
      type: "content",
      resource: e.target.id,
      chat_id: chat_id,
    }));
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
      // console.log("socket message")
      handler(e.data, false);
    };
    socket.onopen = () => {
      // console.log("socket opened");
    };
    socket.onclose = () => {
      // console.log("socket closed");
    };
  }

  // getRegistrationDetails() {
  //   //AJAX html request
  //   httpRequest = new XMLHttpRequest();
  //   if (!httpRequest) {
  //     console.log("Giving up :( Cannot create an XMLHTTP instance'");
  //   }
  //   url = "ws://localhost:8080/";
  //   httpRequest.onreadystatechange = sendContents;
  //   httpRequest.open("POST", url);
  //   httpRequest.setRequestHeader(
  //     "Content-type",
  //     "application/x-www-form-urlencoded"
  //   );
  //   var fd = new FormData();
  //   fd.set("nickname", document.getElementById("reg-nickname").value);
  //   fd.set("email", document.getElementById("reg-email").value);
  //   fd.set("password", document.getElementsByClassName("reg-password").value);

  //   httpRequest.send(fd);

  //   function sendContents() {
  //     if (httpRequest.readyState === 4) {
  //       if (httpRequest.Status === 200) {
  //         alert(httpRequest.responseText);
  //       }
  //     }
  //   }
  // }
}

function containsNumber(str) {
  return /[0-9]/.test(str);
}