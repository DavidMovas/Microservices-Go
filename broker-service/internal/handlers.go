package internal

import (
	"bytes"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
)

type RequestPayload struct {
	Action string      `json:"action"`
	Auth   AuthPayload `json:"auth,omitempty"`
}

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (a *App) Broker(w http.ResponseWriter, r *http.Request) {
	payload := jsonResponse{
		Error:   false,
		Message: "Hit the broker",
	}

	_ = a.writeJSON(w, http.StatusOK, payload)
}

func (a *App) Authenticate(w http.ResponseWriter, r *http.Request) {
	var requestPayload RequestPayload

	err := a.readJSON(w, r, &requestPayload)
	if err != nil {
		_ = a.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	switch requestPayload.Action {
	case "login":
		a.Login(w, requestPayload.Auth)
	default:
		_ = a.errorJSON(w, errors.New("unknown action"), http.StatusBadRequest)
	}
}

func (a *App) Login(w http.ResponseWriter, pl AuthPayload) {
	jsonData, _ := json.MarshalIndent(pl, "", "\t")

	slog.Info("Calling auth service")

	request, err := http.NewRequest("POST", "http://auth-service:8020/auth", bytes.NewBuffer(jsonData))
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
