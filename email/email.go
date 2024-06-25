package email

import (
	"fmt"

	"github.com/michee/e-commerce/model"
	"gopkg.in/gomail.v2"
)

const port = ":8080"

// Fonction pour l'envoi de l'email sur l'adresse de l'utilisateur pour activer son compte
func SendVerificationAccount(u *model.User) {
	m := gomail.NewMessage()
	m.SetHeader("From", "voteprojet@gmail.com")
	m.SetHeader("To", u.Email)
	m.SetHeader("Subject", "Veuillez activer votre compte d'utilisateur")
	m.SetBody("text/html", fmt.Sprintf("Cliquer sur <a href=\"http://localhost%s/auth/verify?token=%s\">Ici</a> pour verifier votre adresse email", port, u.VerificationToken))

	d := gomail.NewDialer("smtp.gmail.com", 587, "voteprojet@gmail.com", "jmbd aicq hdov mvyq")
	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}
}

// Fonction pour la reinitialisation du mot de passe de l'utilisateur par email
func SendResetPasswordAccoubt(u *model.User) {
	m := gomail.NewMessage()
	m.SetHeader("From", "voteprojet@gmail.com")
	m.SetHeader("To", u.Email)
	m.SetHeader("Subject", "Reinitialisation du mot de passe")
	m.SetBody("text/html", fmt.Sprintf("Cliquez <a href=\"http://localhost%s/auth/reset-password-email?token=%s\">ici</a> pour reinitialiser le mot de passe. Ce lien est valide pour une duree d'une heure.", port, u.TokenPassword))

	d := gomail.NewDialer("smtp.gmail.com", 587, "voteprojet@gmail.com", "jmbd aicq hdov mvyq")
	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}
}
