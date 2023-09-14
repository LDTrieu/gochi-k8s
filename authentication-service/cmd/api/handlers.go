package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
)

func (app *Config) Authenticate(w http.ResponseWriter, r *http.Request) {
	log.Println("auth...1")
	var requestPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)

		return
	}
	log.Println("auth...2")
	user, err := app.Models.User.GetByEmail(requestPayload.Email)
	if err != nil {
		app.errorJSON(w, errors.New("invalid creatials"), http.StatusBadRequest)
		return
	}
	log.Println("auth...3", user.Email)
	valid, err := user.PasswordMatches(requestPayload.Password)
	if err != nil || !valid {
		app.errorJSON(w, errors.New("user invalid"), http.StatusBadRequest)
		return
	}
	log.Println("auth...4")
	payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Logged in user %s", user.Email),
		Data:    user,
	}
	log.Println("auth...5", payload.Data)
	app.writeJSON(w, http.StatusAccepted, payload)
}
