package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
)

func (app *Config) routes() http.Handler {
	mux := chi.NewRouter()
	mux.Use(middleware.Recoverer)
	mux.Use(middleware.Heartbeat("/ping"))
	mux.Get("/check", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("logger_exist"))
	})

	mux.Mount("/", app.webRouter())
	mux.Mount("/api", app.apiRouter())

	return mux
}

// apiRouter is for api routes (no session load)
func (app *Config) apiRouter() http.Handler {
	mux := chi.NewRouter()

	// specify who is allowed to connect to our API service
	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	mux.Post("/log", app.WriteLog)
	return mux
}

func (app *Config) webRouter() http.Handler {
	mux := chi.NewRouter()
	mux.Use(app.SessionLoad)

	mux.Get("/", app.LoginPage)
	mux.Get("/login", app.LoginPage)
	mux.Post("/login", app.LoginPagePost)
	mux.Get("/logout", app.Logout)
	
	mux.Get("/check", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("check_exist"))
	})

	mux.Route("/admin", func(mux chi.Router) {
		mux.Use(app.Auth)
		mux.Get("/dasboard", app.Dashboard)
	})

	return mux

}
