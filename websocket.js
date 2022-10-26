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
    div.innerHTML = "<b>" + m.timestamp + " </b>" + "<br>" + "<b>" + m.nickname + ":</b> " + m.text;
    let cself = (myself) ? "self" : "";
    div.className = "msg " + cself;
    document.getElementById("msgcontainer").appendChild(div);
    div.after(msgContainer)
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
        contentSocket.sendContentRequest(event)
      });
      user.innerHTML = p.nickname
      user.className = "presence " + p.nickname
      document.getElementById("presencecontainer").appendChild(user)
    }
  }

  postHandler(text) {
    const m = JSON.parse(text)
    for (let p of m.posts) {
      let post = document.createElement("div");
      post.className = "submittedpost " + p.postid 
      post.innerHTML = "<b>Title: " + p.title + "</b>" + "<br>" + "Nickname: " + p.nickname + "<br>" + "Category/Categories: " + p.categories + "<br>" + p.body + "<br>";
      let button = document.createElement("button")
      button.classname = "addcomment"
      button.innerHTML = "Comments"
      button.setAttribute("data-postid", p.post_id)
      button.addEventListener('click', function (event) {
        event.target.id = "comment"
        contentSocket.sendContentRequest(event)
        for (let c of p.Comments) {
          let comment = document.createElement("div");
          comment.className = "submittedcomment" + c.commentid
          comment.innerHTML = "Nickname: " + c.nickname + "<br>" + c.body;
          document.getElementById("submittedpost " + c.postid).appendChild(comment)
        }
        if ((p.Comments).length > 0) {
          event.target.id = "post"
          postSocket.sendSubmittedCommentsRequest(p.post_id)
        }
      });
      post.appendChild(button)
      document.getElementById("submittedposts").appendChild(post)
    }
  }

  sendNewCommentRequest(e) {
    let m = {
      type: 'post',
      timestamp: "",
      posts: [
        {
          postid: e.target.post_id,
          username: e.target.username,
          title: document.getElementById('posttitle').value,
          categories: document.getElementById('category').value,
          body: document.getElementById('postbody').value,
          comments: [
            {
              commentid: "",
              postid: e.target.post_id,
              username: "",
              body: document.getElementById('commentbody').value,
            }
          ]
        }
      ]
    }
    this.mysocket.send(JSON.stringify(m));
    document.getElementById('commentbody').value = ""
  }

  // makes a call to the backend for comments saved in the database
  sendSubmittedCommentsRequest(postid) {
    this.mysocket.send(JSON.stringify({
      type: "post",
      return: postid,
    }));
  }
  //commentHandler(text) {
  //  const m = JSON.parse(text)
    //for (let p of m.posts) {
      //for (let c of p.comments) {

      //}
      //let comment = document.createElement("div");
      //comment.className = "submittedcomment" + c.commentid
      //comment.innerHTML = "Nickname: " + c.nickname + "<br>" + c.body;
      //document.getElementById("").appendChild(comment)
    //}
  //}

  // TODO: add timestamp
  sendNewPostRequest(e) {
    let m = {
      type: 'post',
      timestamp: "time",
      posts: [
        {
          postid: e.target.postid,
          username: e.target.nickname,
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
      return: "all posts",
    }));
  }

  // TODO: insert username variable
  sendContentRequest(e) {
    this.mysocket.send(JSON.stringify({
      type: "content",
      nickname: "?",
      resource: e.target.id,
    }));
  }

  // TODO: insert username variable
  // requestChat() {
  //   let m = {
  //     type: 'chat',
  //     text: document.getElementById("chatIPT").value,
  //     timestamp: time(),
  //     nickname: "?",
  //   }
  //   this.mysocket.send(JSON.stringify(m));
  //   document.getElementById("chatIPT").value = ""
  // }

  // keypress(e) {
  //   if (e.keyCode == 13) {
  //     this.wsType = e.target.id.slice(0, -3)
  //     switch (this.wsType) {
  //       case 'chat':
  //         this.requestChat()
  //         break;
  //       default:
  //         console.log("keypress registered for unknown wsType")
  //         break;
  //     }
  //   }
  // }

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
  //   fd.set("username", document.getElementById("reg-username").value);
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
