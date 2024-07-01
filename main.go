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
	})

	r.Use(jwtauth.Verifier(tokenAuth))
	r.Use(jwtauth.Authenticator(tokenAuth))

	r.Route("/user", func(r chi.Router) {

		r.Get("/", controller.GetUser)

		r.Route("/{userId}", func(r chi.Router) {
			r.Post("/", controller.LogoutUser)
			r.Get("/", controller.GetUserById)
			r.Patch("/", controller.UpddateUser)
			r.Delete("/", controller.DeleteUser)
		})
	})

	r.Route("/hero", func(r chi.Router) {

		r.With(provider.AdminOnly).Post("/", controller.CreateHero)
		r.Get("/", controller.GetHero)

		r.Route("/{heroId}", func(r chi.Router) {
			r.Get("/", controller.GetHeroById)
			r.With(provider.AdminOnly).Patch("/", controller.UpdateHero)
			r.With(provider.AdminOnly).Delete("/", controller.DeleteHero)
		})
	})

	r.Route("/banner", func(r chi.Router) {
		r.With(provider.AdminOnly).Post("/", controller.CreateBanner)
		r.Get("/", controller.GetBanner)

		r.Route("/{bannerId}", func(r chi.Router) {
			r.Get("/", controller.GetBannerById)
			r.With(provider.AdminOnly).Patch("/", controller.UpdateBanner)
			r.With(provider.AdminOnly).Delete("/", controller.DeleteBanner)
		})
	})

	fmt.Printf("le serveur fonctionne sur http://localhost%s", Port)

	http.ListenAndServe(Port, r)
}
