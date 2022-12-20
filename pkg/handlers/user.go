package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"real-time-forum/pkg/authentication"
)

func Me(w http.ResponseWriter, r *http.Request) {
	if !authentication.RequestIsLoggedIn(r) {
		fmt.Println("request to /me was unauthorised")
		w.WriteHeader(http.StatusUnauthorized)
		ServeHomePage(w, r)
		return
	}

	user, err := authentication.GetUserFromRequest(r)
	if err != nil {
		fmt.Println("request to /me failed when getting the user")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	b, err := json.Marshal(user)
	if err != nil {
		fmt.Println("request to /me failed when marshalling json")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if n, err := w.Write(b); err != nil {
		panic(err)
	} else if n != len(b) {
		panic("request to /me resulted in an incorrect number of bytes being written to the cient")
	}
}
