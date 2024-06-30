package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io"
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

	hashedPassword, _ := utils.HashPassword(user.Password)
	emailToken := utils.GenerateVerificationToken()
	user.Password = hashedPassword
	user.Email = emailToken
	
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
		if userUpdate.Password != "" {
			hashedPassword, _ := utils.HashPassword(userUpdate.Password)
			u.Password = hashedPassword
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

	if !user.EmailVerified {
		http.Error(w, "Email not verified", http.StatusUnauthorized)
		return
	}

	if user.LoginGoogle {
		http.Error(w, "Please log in using Google", http.StatusUnauthorized)
		return
	}

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

func ForgotPasswordHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	var user model.User
	if err := model.Db.Where("email = ?", req.Email).First(&user).Error; err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	if err := user.GeneratePasswordToken(); err != nil {
		http.Error(w, "Failed to generate reset token", http.StatusInternalServerError)
		return
	}

	email.SendResetPasswordAccoubt(&user)

	utils.RespondWithJSON(w, http.StatusOK, "Veuillez verifier votre email pour la reinitialisation du mot de passe", nil)

}

func ResetPasswordEmail(w http.ResponseWriter, r *http.Request) {
	// Parse le fichier HTML
	tmpl := template.Must(template.ParseFiles("/home/michee/go_project/E-commerce/template/emailPassword.tmpl"))

	// Définir le content-type à text/html
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)

	// Exécuter le template avec une structure de données vide
	tmpl.Execute(w, nil)
}

func ResetPasswordHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
			TokenPassword string `json:"tokenPassword"`
			Password      string `json:"password"`
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
	}
	fmt.Println("Received body:", string(body))

	if err := json.Unmarshal(body, &req); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
	}

	fmt.Println("Received tokenPassword:", req.TokenPassword)

	user, err := model.FindUserPasswordToken(req.TokenPassword)
	if err != nil {
			if err.Error() == "reset token has expired" {
					http.Error(w, "Reset token has expired", http.StatusUnauthorized)
			} else {
					http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			}
			return
	}

	if user.TokenPassword != req.TokenPassword {
			http.Error(w, "Invalid reset token", http.StatusUnauthorized)
			return
	}

	err = user.UpdatePassword(req.Password)
	if err != nil {
			http.Error(w, "Failed to update password", http.StatusInternalServerError)
			return
	}

	utils.RespondWithJSON(w, http.StatusOK, "Password updated successfully", nil)
}


