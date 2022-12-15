package authentication

import (
	"encoding/json"
	"fmt"
	"net/http"
	"real-time-forum/pkg/database"

	uuid "github.com/satori/go.uuid"
)

// RegistrationHanlder handles requests for registration from clients and responds with status codes
func RegistrationHandler(w http.ResponseWriter, r *http.Request) {
	reg, err := parseAuthForm(r)
	reg.Password = passwordHash(reg.Password)
	reg.ID = uuid.NewV4().String()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if reg.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	existingUsers, err := database.GetUsers()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// handle user conflicts, passwords do not conflict
	for _, u := range existingUsers {
		diff := u.Compare(reg)
		if len(diff) > 0 {
			w.WriteHeader(http.StatusConflict)
			b, err := json.Marshal(diff)
			if err != nil {
				fmt.Printf("error marshalling json for registration failure. diff: %+v, err: %+v\n", diff, err)
				return
			}
			w.Write(b)
			return
		}
	}

	if _, err := database.CreateUser(reg); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	return
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	auth, err := parseAuthForm(r)
	if err != nil {
		fmt.Printf("unable to parse auth form in loginHandler: %+v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err := database.GetUserByNickname(auth)
	if err != nil {
		fmt.Printf("unable to get user from auth request for '%s': %+v", user.Nickname, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !checkPwHash(auth.Password, user.Password) {
		fmt.Printf("unable to verify password for user '%s'\n", user.Nickname)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	cookie, _, err := database.CreateSession(user)
	if err != nil {
		fmt.Printf("unable to create session for user '%s' in db: %+v\n", user.Nickname, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, cookie)
}
