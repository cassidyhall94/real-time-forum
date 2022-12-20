package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	auth "real-time-forum/pkg/authentication"
	"real-time-forum/pkg/database"

	uuid "github.com/satori/go.uuid"
)

// RegistrationHanlder handles requests for registration from clients and responds with status codes
func RegistrationHandler(w http.ResponseWriter, r *http.Request) {
	reg, err := auth.ParseAuthForm(r)
	reg.Password = auth.PasswordHash(reg.Password)
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
	a, err := auth.ParseAuthForm(r)
	if err != nil {
		fmt.Printf("unable to parse auth form in loginHandler: %+v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err := database.GetUserByNickname(a)
	if err != nil {
		fmt.Printf("unable to get user from auth request for '%s': %+v\n", a.Nickname, err)
		ServeHomePage(w, r)
		return
	}

	if !auth.CheckPwHash(a.Password, user.Password) {
		fmt.Printf("unable to verify password for user '%s'\n", user.Nickname)
		ServeHomePage(w, r)
		return
	}

	cookie, _, err := database.CreateSession(user)
	if err != nil {
		fmt.Printf("unable to create session for user '%s' in db: %+v\n", user.Nickname, err)
		ServeHomePage(w, r)
		return
	}

	http.SetCookie(w, cookie)
	r.AddCookie(cookie)
	ServeHomePage(w, r)
}
