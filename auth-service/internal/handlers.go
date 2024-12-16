package internal

import (
	"bytes"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
)

func (a *App) Authenticate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var requestPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	slog.Info("Calling auth service")

	err := a.readJSON(w, r, &requestPayload)
	if err != nil {
		_ = a.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	user, err := a.Models.GetByEmail(ctx, requestPayload.Email)
	if err != nil {
		_ = a.errorJSON(w, errors.New("invalid credentials"), http.StatusBadRequest)
		return
	}

	slog.Info("User found", "user", user)

	a.Models.User = *user
	valid, err := a.Models.PasswordMatches(requestPayload.Password)
	if err != nil || !valid {
		_ = a.errorJSON(w, errors.New("invalid credentials"), http.StatusBadRequest)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: "Authenticated",
		Data:    user,
	}

	slog.Info("Authenticated user", "user", user)
	if err = a.logRequest("Authenticated user", "Email:", user.Email); err != nil {
		slog.Error("Error calling logger service", "ERROR", err)
	}

	_ = a.writeJSON(w, http.StatusAccepted, payload)
}

func (a *App) logRequest(name string, args ...any) error {
	var entry = struct {
		Name string `json:"name"`
		Args []any  `json:"data"`
	}{
		Name: name,
		Args: args,
	}

	jsonData, _ := json.Marshal(entry)
	logServiceURL := "http://logger-service:8030/log"

	request, err := http.NewRequest("POST", logServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return err
	}

	defer func() {
		_ = response.Body.Close()
	}()

	if response.StatusCode != http.StatusAccepted {
		return errors.New("error calling logger service")
	}

	return nil
}
