package authentication

import (
	"fmt"
	"net/http"
	"net/url"
	"real-time-forum/pkg/database"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// parseRegForm attempts to parse a form from the request and returns a database.User containing the form data
func ParseAuthForm(r *http.Request) (database.User, error) {
	// Parse the form data
	err := r.ParseForm()
	if err != nil {
		return database.User{}, fmt.Errorf("unable to parse reg form: %w", err)
	}

	getFormVal := func(f url.Values, names ...string) string {
		for _, name := range names {
			c := f.Get(name)
			if c != "" {
				return c
			}
		}
		return ""
	}

	return database.User{
		Nickname:  getFormVal(r.Form, "nickname", "login-nickname"),
		Age:       r.Form.Get("age"),
		Gender:    r.Form.Get("gender"),
		FirstName: r.Form.Get("fname"),
		LastName:  r.Form.Get("lname"),
		Email:     r.Form.Get("email"),
		Password:  getFormVal(r.Form, "password", "login-password"),
	}, nil
}

// requestIsLoggedIn validates sessions from a http.Request
func RequestIsLoggedIn(r *http.Request) bool {
	return cookieIsLoggedIn(extractSessionCookie(r))
}

// cookieIsLoggedIn validates sessions from a http.Cookie
func cookieIsLoggedIn(c *http.Cookie) bool {
	if c == nil {
		return false
	}

	sess, err := database.GetSessions()
	if err != nil {
		return false
	}

	if validateSessCookie(c, sess) {
		return true
	}
	return false
}

func extractSessionCookie(r *http.Request) *http.Cookie {
	for _, c := range r.Cookies() {
		if c.Name == database.SessionCookieName {
			return c
		}
	}
	return nil
}

func validateSessCookie(c *http.Cookie, sess []*database.Session) bool {
	if c == nil || len(sess) == 0 {
		return false
	}
	for _, s := range sess {
		if s.SessionID == c.Value {
			return !isCookieExpired(s)
		}
	}

	return false
}

func isCookieExpired(sess *database.Session) bool {
	expirationTime, err := time.Parse("2006-01-02 15:04:05-07:00", sess.ExpiryTime)
	if err != nil {
		fmt.Printf("unable to validate expiry time on session cookie: %+v\n", sess)
		return true
	}
	if expirationTime.Before(time.Now()) {
		return true
	}
	return false
}

func PasswordHash(str string) string {
	hashedPw, err := bcrypt.GenerateFromPassword([]byte(str), 8)
	if err != nil {
		fmt.Println("unable to hash password")
		return ""
	}
	return string(hashedPw)
}

func CheckPwHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GetUserFromRequest(r *http.Request) (*database.User, error) {
	cookie := extractSessionCookie(r)
	if cookie == nil {
		return nil, fmt.Errorf("no session cookie in request")
	}
	return database.GetUserFromSessionID(&database.Session{
		SessionID: cookie.Value,
	})
}

