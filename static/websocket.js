//TODO: fix const time as it is not formatted correctly and add where time/date is needed
const time = () => { return new Date().toLocaleString() };
let clickedParticipantID
// let participant_ids

class MySocket {
  wsType;
  handler;
  constructor() {
    this.mysocket = null;
  }

  presenceHandler(text) {
    const m = JSON.parse(text)
    let presences = document.getElementById("presencecontainer")
    let arr = Array.from(presences.childNodes)
    for (let p of m.presences) {
      if (p.id !== getIdValue()) {
        let user = document.createElement("button");
        user.addEventListener('click', async function (event) {
          clickedParticipantID = user.id
          event.target.id = "chat"
          let participant_ids = [getIdValue(), clickedParticipantID]
          // TODO: pull participant_ids out of this function
          await contentSocket.sendChatContentRequest(event)
          chatSocket.sendNewChatRequest(event, "", participant_ids)
        });

        let existingPresences = (arr.filter(item => item.textContent === p.user.nickname))

        if (existingPresences.length < 1) {

          user.id = p.user.id
          user.innerHTML = p.user.nickname
          user.style.color = 'white'
          user.className = "presence " + p.user.nickname

          presences.appendChild(user)
        }
      }
    }
  }

  sendNewChatRequest(event, inputText = "", participant_ids = []) {
    console.log("new chat request " + inputText)
    let chats = []
    if (inputText !== "") {
      chats = [
        {
          sender: {
            id: getIdValue(),
          },
          date: time(),
          body: inputText,
        }
      ]
    }
    let participants = []
    for (let p of participant_ids) {
      participants.push({ id: p })
    }
    console.log("sendNewChatRequest participants: ", participants)
    let m = {
      type: 'chat',
      timestamp: time(),
      conversations: [
        {
          participants: participants,
          chats: chats
        }
      ]
    }
    this.mysocket.send(JSON.stringify(m));
    document.getElementById('chatIPT').value = ""
  }

  sendChatContentRequest(e) {
    this.mysocket.send(JSON.stringify({
      type: "content",
      resource: e.target.id,
    }));
  }

  chatHandler(text) {
    const m = JSON.parse(text);
    const submittedchats = document.getElementById("submittedchats");
    submittedchats.textContent = '';
    for (let c of m.conversations) {
      for (let p of c.chats) {
        let chat = document.createElement("div");
        chat.className = "submittedchat";
        chat.id = p.chat_id;
        chat.innerHTML = "<b>Me: " + p.sender.nickname + "</b>" + "<br>" + "<b>Date: " + "</b>" + p.date + "<br>" + p.body + "<br><br>";
        submittedchats.appendChild(chat);
      }
    }
  }

  keypress(e) {
    if (e.keyCode == 13) {
      this.wsType = e.target.id.slice(0, -3)
      if (this.wsType = 'chat') {
        this.sendNewChatRequest(e, e.target.value, [getIdValue(), clickedParticipantID])
      }
    }
  }

  contentHandler(text) {
    const c = JSON.parse(text)
    document.getElementById("content").innerHTML = c.body;
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
      checkAndAppendDiv("submittedposts", "content")
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

  sendNewPresenceUpdate(e) {
    let m = {
      type: 'presence',
      timestamp: time(),
      presences: [
        {
          // nickname: e.target.nickname,
          id: getCookieName(),
          nickname: getCookieName(),
          online: "hello",
          last_contacted_time: "0",
        }
      ]
    }
    // this.mysocket.send(JSON.stringify(m))
    console.log("asking for update")
    this.mysocket.send(JSON.stringify(m))
    console.log(m)
    console.log("asked for update")
  }

  sendNewPostRequest(e) {
    let m = {
      type: 'post',
      timestamp: time(),
      posts: [
        {
          // nickname: e.target.nickname,
          nickname: getCookieName(),
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
    const socket = new WebSocket("ws://localhost:8080/" + URI);
    this.mysocket = socket;
    this.handler = handler
    socket.onmessage = (e) => {
      // console.log("socket message", e)
      handler(e.data, false);
    };
    socket.onopen = () => {
      // console.log("socket opened");
    };
    socket.onclose = () => {
      // TODO: make this not so nuts
      // console.log("socket closed");
      // this.connectSocket(this.wsType, this.handler)
    };
    socket.onerror = (event) => {
      console.error(event)
      this.connectSocket(this.wsType, this.handler)
    }
  }
}

function containsNumber(str) {
  return /[0-9]/.test(str);
}

function getIdValue() {
  return document.cookie.split(";")[0].split("=")[1]
}

function checkAndAppendDiv(divId, targetId) {
  // Check if the div element with the specified id exists in the DOM
  let div = document.getElementById(divId);

  // If the div element does not exist, create it and append it to the target element
  if (!div) {
    div = document.createElement("div");
    div.id = divId;
    document.getElementById(targetId).appendChild(div);
  }
}