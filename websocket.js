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
  // registerHandler(text){
  //   console.log("register handler")
  // }
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

}
function containsNumber(str) {
  return /[0-9]/.test(str);
}
  
//object to store form data
let registerForm = {
  nickname : "",
  age: "",
  gender: "",
  fName: "",
  lName: "",
  email: "",
  password: "",
}

let loginForm ={
  nickname:"",
  password:"",
}

//gets registration form details
function getRegDetails(){

    //creates array of gender radio buttons 
  let genderRadios = Array.from (document.getElementsByName('gender'))
  for(let i=0; i <genderRadios.length; i ++){
    // console.log(genderRadios[i].checked)
    if(genderRadios[i].checked){ //stores checked value
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
    
// SEND DATA TO BACKEND USING FETCH
    fetch("/register",{
      headers:{
        'Accept':'application/json',
        'Content-Type': 'application/json'
      },
      method: "POST",
      body:jsonRegForm
    }).then((response)=>{
      response.text().then(function (jsonRegForm){
        let result = JSON.parse(jsonRegForm)
        console.log(result)
      })
    }).catch((error)=>{
      console.log(error)
    })

    document.getElementById('register').reset()
    alert("successfully registered")
}

function loginFormData(){
  loginForm.nickname = document.getElementById('nickname-login').value
  loginForm.password = document.getElementById('password-login').value

  // console.log(loginForm)

  let loginFormJSON = JSON.stringify(loginForm)
  // console.log(loginFormJSON)
  let logindata = {nickname:"",
  password:"",}

  fetch("/login",{
      headers:{
        'Accept':'application/json',
        'Content-Type': 'application/json'
      },
      method: "POST",
      body:loginFormJSON
    }).then((response)=>{
      response.text().then(function (loginFormJSON){
        let result = JSON.parse(loginFormJSON)
        console.log("parse",result)
        if (result == null){
          alert("incorrect username or password")
        } else{
          logindata.nickname = result[0].Nickname
          logindata.password = result[0].Password
          alert("you are logged in ")
        }
      })
    }).catch((error)=>{
      console.log(error)
    })
    console.log("logindata",logindata)

  document.getElementById('login-form').reset()
  // console.log(loginForm)


  //  fetch("/login",{
  //     headers:{
  //       'Accept':'application/json',
  //       'Content-Type': 'application/json'
  //     },
  //     method: "GET",
     
  //   }).then((response)=>{
  //     response.text().then(function (t){
  //       let result = JSON.parse(t)
  //       console.log(result)
  //     })
  //     // return response.json()
  //   }).then((data)=>{
  //     console.log("backend",data)
  //   }).catch((error)=>{
  //     console.log(error)
  //   })

  document.getElementById('login-form').reset()
  // console.log(t)

}
