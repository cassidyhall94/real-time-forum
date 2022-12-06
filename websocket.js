//TODO: fix const time as it is not formatted correctly and add where time/date is needed
const time = () => { return new Date().toLocaleString() };
// let loggedInUserID

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
          //TODO CHAT: replace sender with getCookieName() once the userID is retrieved
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

  sendChatContentRequest(e, participant_ids) {
    this.mysocket.send(JSON.stringify({
      type: "content",
      resource: e.target.id,
      participant_ids: participant_ids,
    }));
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

  presenceHandler(text) {
    const m = JSON.parse(text)
    for (let p of m.presences) {
      const consp = p
      let user = document.createElement("button");
      user.addEventListener('click', function (event, user = consp) {
        event.target.id = "chat"
        participant_ids = [GetLoggedInUserID(text), user.id]
        contentSocket.sendChatContentRequest(event, participant_ids)
      });
      user.id = p.id
      user.innerHTML = p.nickname
      user.style.color = 'white'
      user.className = "presence " + p.nickname
      document.getElementById("presencecontainer").appendChild(user)
    }
  }

  GetLoggedInUserID(text) {
    // match nickname with correct userID
    const m = JSON.parse(text)
    for (let user of m.presences) {
      if (user.nickname === getCookieName()) {
        return user.ID
      }
    }
  }

  keypress(e) {
    if (e.keyCode == 13) {
      this.wsType = e.target.id.slice(0, -3)
      if (this.wsType = 'chat') {
        this.sendNewChatRequest()
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
      // console.log(document.getElementById("submittedposts"))
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
}
function containsNumber(str) {
  return /[0-9]/.test(str);
}

//object to store form data
let registerForm = {
  nickname: "",
  age: "",
  gender: "",
  fName: "",
  lName: "",
  email: "",
  password: "",
  loggedin: "false",
}

let loginForm = {
  nickname: "",
  password: "",
}

//******************* */gets registration form details*******************************
function getRegDetails() {
  //creates array of gender radio buttons 
  let genderRadios = Array.from(document.getElementsByName('gender'))
  for (let i = 0; i < genderRadios.length; i++) {
    // console.log(genderRadios[i].checked)
    if (genderRadios[i].checked) { //stores checked value
      registerForm.gender = genderRadios[i].value
    }
  }
  // POPULATE REGISTER FORM WITH FORM VALUES
  registerForm.nickname = document.getElementById('nickname').value
  registerForm.age = document.getElementById('age').value
  registerForm.firstname = document.getElementById('fname').value
  registerForm.lastname = document.getElementById('lname').value
  registerForm.email = document.getElementById('email').value
  registerForm.password = document.getElementById('password').value
  //convert data to JSON
  let jsonRegForm = JSON.stringify(registerForm)
  // console.log(jsonRegForm)
  if (registerForm.password.length < 5) {
    registerForm.password = ""
  }

  // SEND DATA TO BACKEND USING FETCH
  console.log(registerForm)
  if (registerForm.nickname != "" && registerForm.email != "" && registerForm.password != "") {

    fetch("/register", {
      headers: {
        'Accept': 'application/json',
        'Content-Type': 'application/json'
      },
      method: "POST",
      body: jsonRegForm
    }).then((response) => {
      response.text().then(function (jsonRegForm) {
        let result = JSON.parse(jsonRegForm)
        // console.log("register", result)
        //cheks result value and only clears form when registered
        if (result == "registered") {
          document.getElementById('register').reset()
          alert("successfully registered")
        } else {
          alert(result)
        }
      })
    }).catch((error) => {
      console.log(error)
    })
  }
}
// **********************************LOGIN*******************************************
function loginFormData() {
  loginForm.nickname = document.getElementById('nickname-login').value
  loginForm.password = document.getElementById('password-login').value
  let loginFormJSON = JSON.stringify(loginForm)
  // console.log(loginFormJSON)
  let logindata = {
    nickname: "",
    password: "",
  }
  // let id = ""
  getCookieName()
  fetch("/login", {
    headers: {
      'Accept': 'application/json',
      'Content-Type': 'application/json'
    },
    method: "POST",
    body: loginFormJSON
  }).then((response) => {
    response.text().then(function (loginFormJSON) {
      let result = JSON.parse(loginFormJSON)
      console.log("parse", result)

      if (result == null) {
        alert("incorrect username or password")
      } else {
        logindata.nickname = result[0].nickname
        // logindata.password = result[0].password
        user.innerText = `Hello ${document.cookie.match(logindata.nickname)}`
        alert("you are logged in ")
      }
    })
  }).catch((error) => {
    console.log(error)
  })
  // console.log("logindata",logindata, "hi")
  // console.log( Object.keys(logindata).length)
  // console.log(JSON.stringify(logindata))
  document.getElementById('login-form').reset()
  let user = document.getElementById('welcome')
  // document.getElementById('login-form').reset()
  // console.log(t)
}
// ********************************LOGOUT******************************
function Logout() {
  let logout = {
    nickname: "",
  }

  logout.nickname = document.getElementById('welcome')
  logout.nickname = logout.nickname.textContent.replace("Hello", '')
  let user = document.getElementById('welcome')
  user.innerText = ""
  // console.log("logout", user.textContent.replace("Hello", ''))

  let parseUser = JSON.stringify(logout)
  fetch("/logout", {
    headers: {
      'Accept': 'application/json',
      'Content-Type': 'application/json'
    },
    method: "POST",
    body: parseUser
  }).then((response) => {
    response.text().then(function (parseUser) {
      // let result = JSON.parse(parseUser)
      // console.log("parse",result)
    })
  }).catch((error) => {
    console.log(error)
  })
  alert("you are now logged out")
}
function getCookieName() {
  let cookies = document.cookie.split(";")
  let lastCookieName = cookies[cookies.length - 1].split("=")[0].replace(" ", '')
  // console.log("cookie",cookies, "length", cookies.length)
  return lastCookieName
  // console.log("h",lastCookieName)
}
function getIdValue() {
  //     let cookies = document.cookie.split(";")
  // let lastCookieName = cookies[cookies.length -1].split("=")[1].replace(" ", '')
  // return lastCookieName
  return document.cookie.split(";")[0].split("=")[1]
}
console.log("test,", getIdValue())

console.log("hello", document.cookie.split(";")[0].split("=")[1])
// console.log(cookies[cookies.length-1].split("=")[0])