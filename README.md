PROJECT: https://learn.01founders.co/git/root/public/src/branch/master/subjects/real-time-forum

AUDIT: https://learn.01founders.co/git/root/public/src/branch/master/subjects/real-time-forum/audit

- this all must be LIVE using WebSockets in backend/frontend
    - https://pkg.go.dev/github.com/gorilla/websocket
    - https://medium.com/@bootcampmillionaire/what-i-learned-about-websockets-by-building-a-real-time-chat-application-using-socket-io-3d9e163e504
    - websockets depend on two hosts having a connection, so they run on top of a TCP layer and modify the TCP layer so that the client and server agree for the socket to stay open
        - https://en.wikipedia.org/wiki/Transmission_Control_Protocol
        - https://en.wikipedia.org/wiki/Connection-oriented_communication
    - websocket handshake: how the data being exchanged should be interpreted by both client and server 
    - user does not have to refresh the page to see


- one HTML file only: https://en.wikipedia.org/wiki/Single-page_application

Registration and login
    - form data must include: nickname, age, gender, first and last name, email, password
    - login using either nickname or email with password
    - logout from any page

Creation of posts and comments
    - posts will have categories

Private messaging
    - online/offline section
        - organised by last message sent
            - if user is new, then contacts are ordered alphabetical
        - send private messages to online users
        - must be visible at all times
    - when user clicked on, it loads previous messages with that user
        - reloads the last 10 messages and when user scrolls up, 10 more messages should be shown (without spamming the scroll event)
            - look at [throttle](https://css-tricks.com/debouncing-throttling-explained-examples/#throttle)/debounce for not spamming the scroll event to see messages
    - message format
        - date and username shown

Bonus
- user profiles
- send images through messaging
- code is synchronicity (promises and go routines/channels) to increase performance




