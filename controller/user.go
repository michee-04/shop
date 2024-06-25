package controller

import (
	"encoding/json"
	"errors"
	"html/template"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/michee/e-commerce/email"
	"github.com/michee/e-commerce/model"
	"github.com/michee/e-commerce/provider"
	"github.com/michee/e-commerce/utils"
	"gorm.io/gorm"
)

func CreateUser(w http.ResponseWriter, r *http.Request) {
	user := model.User{}
	utils.ParseBody(r, &user)
	u := user.CreateUser()

	email.SendVerificationAccount(&user)

	res, _ := json.Marshal(u)
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	u := model.GetUser()
	res, _ := json.Marshal(u)
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func GetUserById(w http.ResponseWriter, r *http.Request) {
	userId := chi.URLParam(r, "userId")
	u, _ := model.GetUserById(userId)
	res, _ := json.Marshal(u)
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	userId := chi.URLParam(r, "userId")
	u := model.DeleteUser(userId)
	res, _ := json.Marshal(u)
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)

	utils.RespondWithJSON(w, http.StatusOK, "Delete user successful", map[string]string{"token": string(res)})
}

func UpddateUser(w http.ResponseWriter, r *http.Request) {
	userUpdate := &model.User{}
	utils.ParseBody(r, userUpdate)

	userId := chi.URLParam(r, "userId")

	if !provider.VerificationToken(r, userId) {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	var u model.User
	err := model.Db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("user_id=?", userId).First(&u).Error; err != nil {
			return err
		}

		if userUpdate.Name != "" {
			u.Name = userUpdate.Name
		}
		if userUpdate.Username != "" {
			u.Username = userUpdate.Username
		}
		if userUpdate.Email != "" {
			u.Email = userUpdate.Email
			email.SendVerificationAccount(&u)
		}

		return tx.Save(&u).Error
	})

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "User not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to update user: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	res, err := json.Marshal(u)
	if err != nil {
		http.Error(w, "Failed update user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)

}

// fonction pour le login
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var loginReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	err := json.NewDecoder(r.Body).Decode(&loginReq)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	user, err := model.GetUserByEmail(loginReq.Email)
	if err != nil || !utils.CheckPassordHash(loginReq.Password, user.Password) {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	token, err := provider.GenerateJWT(user.UserId, user.IsAdmin)
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	user.TokenJwt = token
	// Mettre à jour le modèle User dans la base de données avec le nouveau token
	model.Db.Save(&user)

	utils.RespondWithJSON(w, http.StatusOK, "Login successful", map[string]string{"token": token})
}

func LogoutUser(w http.ResponseWriter, r *http.Request) {
	userId := chi.URLParam(r, "userId")
	u, _ := model.GetUserById(userId)

	if u == nil {
		http.Error(w, "user not found", http.StatusInternalServerError)
		return
	}

	if !provider.VerificationToken(r, userId) {
		http.Error(w, "Invalid token", http.StatusInternalServerError)
		return
	}

	u.Logout()

	utils.RespondWithJSON(w, http.StatusOK, "Logout successful", nil)
}

func VerifyHandler(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		http.Error(w, "Missing token", http.StatusBadRequest)
		return
	}

	user, err := model.FindUserByToken(token)
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	err = user.Verify()
	if err != nil {
		log.Printf("Erreur lors de la vérification de l'utilisateur: %v\n", err)
		http.Error(w, "Unable to verify user", http.StatusInternalServerError)
		return
	}

	tmpl := template.Must(template.ParseFiles("/home/michee/go_project/E-commerce/template/emailverification.tmpl"))

	w.Header().Set("content-type", "text/html")
	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, nil)
}
