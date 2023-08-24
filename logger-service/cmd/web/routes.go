package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (app *Config) routes() http.Handler {
	mux := chi.NewRouter()
	mux.Use(middleware.Recoverer)
	mux.Use(middleware.Heartbeat("/ping"))

	mux.Mount("/", app.webRouter())
	//mux.Mount("/api",app.apiRouter)
	return mux
}

func (app *Config) webRouter() http.Handler {
	mux := chi.NewRouter()
	mux.Use(app.SessionLoad)

	mux.Get("/", app.LoginPage)
	mux.Get("login", app.LoginPage)
	//mux.Post("/login", app.LoginPagePost)
	//mux.Get("logout", app.LogoutPage)

	mux.Route("/admin", func(mux chi.Router) {
		mux.Use(app.Auth)
		//mux.Get("/dasboard",app.Dasboard)
	})

	return mux

}
