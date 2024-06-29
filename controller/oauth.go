package controller

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/michee/e-commerce/provider"
	"golang.org/x/oauth2"
)

func GoogleLogin(w http.ResponseWriter, r *http.Request) {
	googleConfig := provider.SetConfigAouth()
	url := googleConfig.AuthCodeURL("randomstate", oauth2.AccessTypeOffline)

	http.Redirect(w, r, url, http.StatusSeeOther)
}


func GoogleCallback(w http.ResponseWriter, r *http.Request) {

	state := r.URL.Query()["state"][0]
	if state != "randomstate" {
		fmt.Fprintln(w, "state don't match")
		return
	}

	code := r.URL.Query()["code"][0]

	googleConfig := provider.SetConfigAouth()

	token, err := googleConfig.Exchange(context.Background(), code)
	if err != nil {
		fmt.Fprintln(w, "Code-Token Exchange failed")
	}

	response, err := http.Get("https://www.googleapis.com/auth/userinfo?access_token" + token.AccessToken)
	if err != nil {
		fmt.Fprintln(w, "User Date fetch failed")
	}

	userDate, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Fprintln(w, "Json Parsing failed")
	}


	fmt.Fprintln(w, string(userDate))
}
