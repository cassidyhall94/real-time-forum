package handlers

import (
	"fmt"
	"net/http"
	"real-time-forum/pkg/authentication"
	"text/template"
)

func ServeHomePage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" || r.URL.Path == "index.html" {
		tpl, err := template.ParseGlob("templates/*")
		if err != nil {
			fmt.Printf("unable to parse templates: %+v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if !authentication.RequestIsLoggedIn(r) {
			if err := tpl.ExecuteTemplate(w, "indexWithLogin.template", nil); err != nil {
				fmt.Printf("unable to render indexWithLogin.template: %+v\n", err)
				w.WriteHeader(http.StatusInternalServerError)
			}
		} else if err := tpl.ExecuteTemplate(w, "index.template", nil); err != nil {
			fmt.Printf("unable to render index.template: %+v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}
