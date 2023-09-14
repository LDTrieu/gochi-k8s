package main

import (
	"broker/event"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/rpc"
)

type RequestPayload struct {
	Action string      `json:"action"`
	Auth   AuthPayload `json:"auth,omitempty"`
	Mail   MailPayload `json:"mail,omitempty"`
	Log    LogPayload  `json:"log,omitempty"`
}

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type MailPayload struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}

type LogPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {

	payload := jsonResponse{
		Error:   false,
		Message: "Hit the broker",
	}
	_ = app.writeJSON(w, http.StatusOK, payload)
}

func (app *Config) HandleSubmission(w http.ResponseWriter, r *http.Request) {

	var requestPayload RequestPayload

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	switch requestPayload.Action {
	case "auth":
		app.authenticate(w, requestPayload.Auth)
		// case "log":
		// 	app.logEventViarabbit(w, requestPayload.Log)
	// case "log":
	// 	app.logItemViaRPC(w, requestPayload.Log)
	case "log":
		log.Println("LOG")
		app.logItem(w, requestPayload.Log)
	case "mail":
		app.sendMail(w, requestPayload.Mail)
	default:
		app.errorJSON(w, errors.New("unknow action"))
	}
}

func (app *Config) authenticate(w http.ResponseWriter, a AuthPayload) {

	// creare some json we 'll sent to auth service
	jsonData, _ := json.MarshalIndent(a, "", "\t")

	// call the service
	//request, err := http.NewRequest("POST", "http://localhost:8080/authenticate", bytes.NewBuffer(jsonData))
	request, err := http.NewRequest("POST", "http://authentication-service/authenticate", bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	defer response.Body.Close()

	// make sure we get back the correct status code
	if response.StatusCode == http.StatusUnauthorized {

		app.errorJSON(w, errors.New("invalid creadentials"))
		return
	} else if response.StatusCode != http.StatusAccepted {

		app.errorJSON(w, errors.New("error calling auth service"))
		return
	}

	// create a varabirel wee 'll read response.body into

	var jsonFromService jsonResponse

	// decode the json from the auth service

	err = json.NewDecoder(response.Body).Decode(&jsonFromService)
	if err != nil {

		app.errorJSON(w, err)
		return
	}

	if jsonFromService.Error {
		app.errorJSON(w, err, http.StatusUnauthorized)
		return
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = "Authenticated!"
	payload.Data = jsonFromService.Data

	app.writeJSON(w, http.StatusAccepted, payload)
}

func (app *Config) sendMail(w http.ResponseWriter, msg MailPayload) {
	jsonData, _ := json.MarshalIndent(msg, "", "\t")

	//mailServiceURL := fmt.Sprintf("http://%s/send", app.GetServiceURL("mail"))
	mailServiceURL := fmt.Sprintf("http://%s/send", "mail-service")

	// now post to the mail service
	request, err := http.NewRequest("POST", mailServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {

		app.errorJSON(w, err)
		return
	}

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {

		_ = app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	defer response.Body.Close()

	// make sure we get back the right status code
	if response.StatusCode != http.StatusAccepted {

		_ = app.errorJSON(w, errors.New("error calling mail service"), http.StatusBadRequest)
		return
	}

	// send json back to our end user
	var payload jsonResponse
	payload.Error = false
	payload.Message = "Message sent to " + msg.To

	// app.writeJSON(w, http.StatusAccepted, payload)
	out, _ := json.MarshalIndent(payload, "", "\t")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)

	_, _ = w.Write(out)

}

func (app *Config) logItem(w http.ResponseWriter, l LogPayload) {
	jsonData, _ := json.MarshalIndent(l, "", "\t")
	logServiceURL := "http://logger-service/api/log"
	request, err := http.NewRequest("POST", logServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	log.Println("request.....1")
	request.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Println("request.....2", err)
		app.errorJSON(w, err)
		return
	}
	log.Println("request.....3", err)
	defer response.Body.Close()
	if response.StatusCode != http.StatusAccepted {
		log.Println("request.....3.1", err)
		app.errorJSON(w, err)
		return
	}

	var payload jsonResponse
	payload.Error = false

	payload.Message = "logged"
	log.Println("request.....4", err)
	app.writeJSON(w, http.StatusAccepted, payload)
}

func (app *Config) logEventViarabbit(w http.ResponseWriter, l LogPayload) {
	err := app.pushToQueue(l.Name, l.Data)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	var payload jsonResponse
	payload.Error = false
	payload.Message = "logged via RabbitMQ"
	app.writeJSON(w, http.StatusAccepted, payload)

}

func (app *Config) pushToQueue(name, msg string) error {
	emitter, err := event.NewEventEmitter(app.Rabbit)
	if err != nil {
		return err
	}

	payload := LogPayload{
		Name: name,
		Data: msg,
	}
	j, _ := json.MarshalIndent(&payload, "", "\t")
	err = emitter.Push(string(j), "log.INFO")
	if err != nil {
		return err
	}
	return nil
}

type RPCPayload struct {
	Name string
	Data string
}

func (app *Config) logItemViaRPC(w http.ResponseWriter, l LogPayload) {
	client, err := rpc.Dial("tcp", "logger-service:5001")
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	rpcPayload := RPCPayload{
		Name: l.Name,
		Data: l.Data,
	}
	var result string

	err = client.Call("RPCServer.LogInfo", rpcPayload, &result)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	payload := jsonResponse{
		Error:   false,
		Message: result,
	}

	app.writeJSON(w, http.StatusAccepted, payload)
}
