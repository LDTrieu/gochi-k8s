package main

import (
	"log"
	"net/http"
)

func (app *Config) SendMail(w http.ResponseWriter, r *http.Request) {
	type mailMessage struct {
		From    string `json:"from"`
		To      string `json:"to"`
		Subject string `json:"subject"`
		Message string `json:"message"`
	}

	var requestPayload mailMessage
	log.Println("email...1")
	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		log.Println("email...1.1", err)
		_ = app.errorJSON(w, err)
		return
	}
	log.Println("email...2")
	msg := Message{
		From:    requestPayload.From,
		To:      requestPayload.To,
		Subject: requestPayload.Subject,
		Data:    requestPayload.Message,
	}

	err = app.Mailer.SendSMTPMessage(msg)
	if err != nil {
		log.Println("email...2.1", err)
		_ = app.errorJSON(w, err)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: "sent to " + requestPayload.To,
	}
	log.Println("email...3", payload)
	_ = app.writeJSON(w, http.StatusAccepted, payload)
}
