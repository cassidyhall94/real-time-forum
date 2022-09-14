PROJECT: https://learn.01founders.co/git/root/public/src/branch/master/subjects/real-time-forum

AUDIT: https://learn.01founders.co/git/root/public/src/branch/master/subjects/real-time-forum/audit

- this all must be LIVE using WebSockets in backend/frontend
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
            - look at throttle/debounce
    - message format
        - date and username shown




