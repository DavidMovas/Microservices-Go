package internal

import (
	"bytes"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
)

const (
	AuthenticateServiceURL = "http://auth-service:8020/auth"
	LoggerServiceURL       = "http://logger-service:8030/log"
	MailServiceURL         = "http://mail-service:8040/send"

	AuthenticateAction = "auth"
	LoggingAction      = "log"
	MailAction         = "mail"
)

type RequestPayload struct {
	Action string        `json:"action"`
	Auth   AuthPayload   `json:"auth,omitempty"`
	Log    LoggerPayload `json:"log,omitempty"`
	Mail   MailPayload   `json:"mail,omitempty"`
}

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoggerPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

type MailPayload struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}

func (a *App) Broker(w http.ResponseWriter, _ *http.Request) {
	payload := jsonResponse{
		Error:   false,
		Message: "Hit the broker",
	}

	_ = a.writeJSON(w, http.StatusOK, payload)
}

func (a *App) HandleSubmission(w http.ResponseWriter, r *http.Request) {
	var requestPayload RequestPayload

	err := a.readJSON(w, r, &requestPayload)
	if err != nil {
		_ = a.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	switch requestPayload.Action {
	case AuthenticateAction:
		a.Login(w, requestPayload.Auth)
	case LoggingAction:
		a.Logger(w, requestPayload.Log)
	case MailAction:
		a.SendMail(w, requestPayload.Mail)
	default:
		_ = a.errorJSON(w, errors.New("unknown action"), http.StatusBadRequest)
	}
}

func (a *App) Login(w http.ResponseWriter, pl AuthPayload) {
	jsonData, _ := json.MarshalIndent(pl, "", "\t")

	slog.Info("Calling auth service")

	request, err := http.NewRequest(http.MethodPost, AuthenticateServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		_ = a.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		_ = a.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	defer func() {
		_ = response.Body.Close()
	}()

	if response.StatusCode == http.StatusUnauthorized {
		_ = a.errorJSON(w, errors.New("invalid credentials"), http.StatusBadRequest)
		return
	}

	if response.StatusCode != http.StatusAccepted {
		_ = a.errorJSON(w, errors.New("error calling auth service"), http.StatusBadRequest)
		return
	}

	var jsonResp jsonResponse
	err = json.NewDecoder(response.Body).Decode(&jsonResp)
	if err != nil {
		_ = a.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	if jsonResp.Error {
		_ = a.errorJSON(w, errors.New(jsonResp.Message), http.StatusUnauthorized)
		return
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = "Authenticated"
	payload.Data = jsonResp.Data

	_ = a.writeJSON(w, http.StatusAccepted, payload)
}

func (a *App) Logger(w http.ResponseWriter, pl LoggerPayload) {
	jsonData, _ := json.MarshalIndent(pl, "", "\t")

	request, err := http.NewRequest(http.MethodPost, LoggerServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		_ = a.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		_ = a.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	defer func() {
		_ = response.Body.Close()
	}()

	if response.StatusCode != http.StatusAccepted {
		_ = a.errorJSON(w, errors.New("error calling logger service"), http.StatusBadRequest)
		return
	}

	var jsonResp jsonResponse
	err = json.NewDecoder(response.Body).Decode(&jsonResp)
	if err != nil {
		_ = a.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	if jsonResp.Error {
		_ = a.errorJSON(w, errors.New(jsonResp.Message), http.StatusInternalServerError)
		return
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = "logged"

	_ = a.writeJSON(w, http.StatusAccepted, payload)
}

func (a *App) SendMail(w http.ResponseWriter, pl MailPayload) {
	jsonData, err := json.MarshalIndent(pl, "", "\t")
	if err != nil {
		_ = a.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	request, err := http.NewRequest(http.MethodPost, MailServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		_ = a.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		_ = a.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	defer func() {
		_ = response.Body.Close()
	}()

	if response.StatusCode != http.StatusAccepted {
		slog.Error("Error sending mail", "ERROR", err)
		_ = a.errorJSON(w, errors.New("error calling mail service"), http.StatusBadRequest)
		return
	}

	var jsonResp jsonResponse
	err = json.NewDecoder(response.Body).Decode(&jsonResp)
	if err != nil {
		_ = a.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	if jsonResp.Error {
		_ = a.errorJSON(w, errors.New(jsonResp.Message), http.StatusInternalServerError)
		return
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = "Mail sent"

	_ = a.writeJSON(w, http.StatusAccepted, payload)
}
