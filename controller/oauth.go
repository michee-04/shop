package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/michee/e-commerce/model"
	"github.com/michee/e-commerce/provider"
	"github.com/michee/e-commerce/utils"
	"golang.org/x/oauth2"
)

func GoogleLogin(w http.ResponseWriter, r *http.Request) {
	googleConfig := provider.SetConfigAouth()
	url := googleConfig.AuthCodeURL("randomstate", oauth2.AccessTypeOffline)

	http.Redirect(w, r, url, http.StatusSeeOther)
}



func GoogleCallback(w http.ResponseWriter, r *http.Request) {
	state := r.URL.Query().Get("state")
	if state != "randomstate" {
		fmt.Fprintln(w, "State don't match")
		return
	}

	code := r.URL.Query().Get("code")
	googleConfig := provider.SetConfigAouth()

	// Échange du code d'autorisation pour un jeton d'accès
	token, err := googleConfig.Exchange(context.Background(), code)
	if err != nil {
		fmt.Fprintln(w, "Code-Token Exchange failed:", err)
		return
	}

	// Création d'un client HTTP avec le jeton d'accès
	client := googleConfig.Client(context.Background(), token)

	// Requête pour obtenir les informations utilisateur
	response, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		fmt.Fprintln(w, "User Data fetch failed:", err)
		return
	}
	defer response.Body.Close()

	// Lecture de la réponse
	userData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Fprintln(w, "Json Parsing failed:", err)
		return
	}

	// Parsing des données utilisateur
	var googleUser struct {
		ID            string `json:"id"`
		Email         string `json:"email"`
		VerifiedEmail bool   `json:"verified_email"`
		Name          string `json:"name"`
		GivenName     string `json:"given_name"`
		FamilyName    string `json:"family_name"`
		Picture       string `json:"picture"`
	}
	if err := json.Unmarshal(userData, &googleUser); err != nil {
		fmt.Fprintln(w, "Json Unmarshal failed:", err)
		return
	}

	// Vérifiez si l'utilisateur existe déjà dans la base de données
	user, err := model.GetUserByEmail(googleUser.Email)
	if err != nil && err.Error() != "user not found" {
		fmt.Fprintln(w, "Database query failed:", err)
		return
	}

	// Si l'utilisateur n'existe pas, créez un nouvel utilisateur
	if user == nil {
		user = &model.User{
			Email:         googleUser.Email,
			EmailVerified: googleUser.VerifiedEmail,
			Name:          googleUser.Name,
			Username:      googleUser.GivenName,
		}
		user.CreateUser()
	}

	// Générer un jeton JWT pour l'utilisateur
	jwtToken, err := provider.GenerateJWT(user.UserId, user.IsAdmin)
	if err != nil {
		fmt.Fprintln(w, "JWT generation failed:", err)
		return
	}

	// Mettre à jour le champ token_jwt de l'utilisateur
	user.TokenJwt = jwtToken
	model.Db.Save(user)

	// Répondre avec le jeton JWT
	utils.RespondWithJSON(w, http.StatusOK, "Login successful", map[string]string{"token": jwtToken})
}