package provider

import (
	"log"
	"net/http"
	"strings"

	"github.com/michee/e-commerce/model"
)

func VerificationToken(r *http.Request, userId string) bool {
	// Extraire le token depuis l'en-tête
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		log.Println("Authorization header is missing")
		return false
	}

	// Supprimer le préfixe "Bearer " du token
	tokenFromHeader := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenFromHeader == authHeader {
		log.Println("Token does not have 'Bearer ' prefix")
		return false // le token n'a pas le préfixe "Bearer "
	}

	log.Printf("Extracted token from header: %s\n", tokenFromHeader)

	// Récupérer l'utilisateur depuis la base de données
	userDetail, _ := model.GetUserById(userId)
	if userDetail == nil {
		log.Println("User not found")
		return false
	}

	// Comparer le token
	log.Printf("Token from header: %s\n", tokenFromHeader)
	log.Printf("Token from database: %s\n", userDetail.TokenJwt)
	if tokenFromHeader != userDetail.TokenJwt {
		log.Println("Tokens do not match")
		return false
	}

	return true
}