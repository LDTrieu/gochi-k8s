package main

import (
	"fmt"
	"log"
	"net/http"
)

const authServiceURL = "http://authentication-service/authenticate"

// JSONPayload is the type for JSON posted to this API
type JSONPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

// WriteLog is the handler to accept a post request consisting of json payload,
// and then write it to Mongo
func (app *Config) WriteLog(w http.ResponseWriter, r *http.Request) {
	// read json into var
	var requestPayload JSONPayload
	_ = app.readJSON(w, r, &requestPayload)
	log.Println("requestPayload", requestPayload)
	// insert the data
	err := app.logEvent(requestPayload.Name, requestPayload.Data)
	if err != nil {
		log.Println("err", err)
		_ = app.errorJSON(w, err, http.StatusBadRequest)
		return
	}
	log.Println("log....2")
	// create the response we'll send back as JSON
	resp := jsonResponse{
		Error:   false,
		Message: "logged",
	}

	// write the response back as JSON
	_ = app.writeJSON(w, http.StatusAccepted, resp)
}

// Logout logs the user out and redirects them to the login page
func (app *Config) Logout(w http.ResponseWriter, r *http.Request) {
	// log the event
	_ = app.logEvent("authentication", fmt.Sprintf("%s logged out of the logger service", app.Session.GetString(r.Context(), "email")))

	// clean up session
	_ = app.Session.Destroy(r.Context())
	_ = app.Session.RenewToken(r.Context())

	// redirect to login page
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (app *Config) LoginPage(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "login.page.gohtml", nil)
}

func (app *Config) LoginPagePost(w http.ResponseWriter, r *http.Request) {
	_ = app.Session.RenewToken(r.Context())
}

func (app *Config) SampleAPI(w http.ResponseWriter, r *http.Request) {
	var requestPayload JSONPayload
	_ = app.readJSON(w, r, requestPayload)

	resp := jsonResponse{
		Error:   false,
		Message: requestPayload.Name,
		Data:    requestPayload.Data,
	}

	_ = app.writeJSON(w, http.StatusAccepted, resp)

}

// Dashboard displays the dashboard page
func (app *Config) Dashboard(w http.ResponseWriter, r *http.Request) {
	// get the list of all log entries from mongo
	logs, err := app.Models.LogEntry.All()
	if err != nil {
		log.Println("Error getting all log entries")
		app.clientError(w, http.StatusBadRequest)
	}

	templateData := make(map[string]any)
	templateData["logs"] = logs

	app.render(w, r, "dashboard.page.gohtml", &TemplateData{
		Data: templateData,
	})
}
