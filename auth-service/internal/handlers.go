package internal

import (
	"errors"
	"net/http"
)

func (a *App) Authenticate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var requestPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

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

	_ = a.writeJSON(w, http.StatusAccepted, payload)
}
