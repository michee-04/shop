package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
	"github.com/michee/e-commerce/controller"
	"github.com/michee/e-commerce/provider"
)

const Port = ":8080"

var tokenAuth *jwtauth.JWTAuth

func main() {

	tokenAuth = jwtauth.New("HS256", []byte("ksQD5adHXZ-5SSJCupcHwBzDi6q5kfr5hdU7Eq5tMmo"), nil)

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.CleanPath)
	r.Use(middleware.RequestID)

	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", controller.CreateUser)
		r.Get("/verify", controller.VerifyHandler)
		r.Get("/google/login", controller.GoogleLogin)
		r.Get("/google/callback", controller.GoogleCallback)
		r.Post("/login", controller.LoginHandler)
		r.Post("/forgot-password", controller.ForgotPasswordHandler)
		r.Get("/reset-password-email", controller.ResetPasswordEmail)
		r.Post("/reset-password", controller.ResetPasswordHandler)
		r.Patch("/set-password", controller.SetPasswordHandler)
	})

	r.Route("/user", func(r chi.Router) {
		r.Get("/", controller.GetUser)
		r.Get("/{userId}", controller.GetUserById)

		r.Route("/{userId}", func(r chi.Router) {
			r.Use(jwtauth.Verifier(tokenAuth))
			r.Use(jwtauth.Authenticator(tokenAuth))

			r.Post("/", controller.LogoutUser)
			r.Patch("/", controller.UpddateUser)
			r.Delete("/", controller.DeleteUser)
		})
	})

	r.Route("/categorie", func(r chi.Router) {

		r.Get("/", controller.GetCategorie)
		r.Get("/{categorieId}", controller.GetCategorieId)

		r.Route("/", func(r chi.Router) {
			r.Use(jwtauth.Verifier(tokenAuth))
			r.Use(jwtauth.Authenticator(tokenAuth))

			r.With(provider.AdminOnly).Post("/", controller.CreateCategorie)
			r.With(provider.AdminOnly).Patch("/{categorieId}", controller.UpdateCategorie)
			r.With(provider.AdminOnly).Delete("/{categorieId}", controller.DeleteCategorie)
		})
	})

	r.Route("/article", func(r chi.Router) {

		r.Get("/", controller.GetArticle)
		r.Get("/{articleId}", controller.GetArticleId)

		r.Route("/", func(r chi.Router) {
			r.Use(jwtauth.Verifier(tokenAuth))
			r.Use(jwtauth.Authenticator(tokenAuth))
			r.With(provider.AdminOnly).Post("/", controller.CreateArticle)
			r.With(provider.AdminOnly).Patch("/{articleId}", controller.UpdateArticle)
			r.With(provider.AdminOnly).Delete("/{articleId}", controller.DeleteArticle)
		})
	})

	r.Route("/hero", func(r chi.Router) {

		r.Get("/", controller.GetHero)
		r.Get("/{heroId}", controller.GetHeroById)

		r.Route("/", func(r chi.Router) {
			r.Use(jwtauth.Verifier(tokenAuth))
			r.Use(jwtauth.Authenticator(tokenAuth))
			r.With(provider.AdminOnly).Post("/", controller.CreateHero)
			r.With(provider.AdminOnly).Patch("/{heroId}", controller.UpdateHero)
			r.With(provider.AdminOnly).Delete("/{heroId}", controller.DeleteHero)
		})
	})

	r.Route("/banner", func(r chi.Router) {

		r.Get("/", controller.GetBanner)
		r.Get("/{bannerId}", controller.GetBannerById)

		r.Route("/", func(r chi.Router) {
			r.Use(jwtauth.Verifier(tokenAuth))
			r.Use(jwtauth.Authenticator(tokenAuth))
			r.With(provider.AdminOnly).Post("/", controller.CreateBanner)

			r.With(provider.AdminOnly).Patch("/{bannerId}", controller.UpdateBanner)
			r.With(provider.AdminOnly).Delete("/{bannerId}", controller.DeleteBanner)
		})
	})

	fmt.Printf("le serveur fonctionne sur http://localhost%s", Port)

	http.ListenAndServe(Port, r)
}
