package provider

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)


func SetConfigAouth() *oauth2.Config {

	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	conf := &oauth2.Config{
		ClientID: os.Getenv("Oauth_Client_Id"),
		ClientSecret: os.Getenv("Oauth_Secret_Key"),
		RedirectURL: "http://localhost:8080/auth/google/callback",
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
			"openid",
		},
		Endpoint: google.Endpoint,
	}


	return conf
}