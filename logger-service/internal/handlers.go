package internal

import (
	"log/slog"
	"net/http"
)

type payload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (a *App) Logger(w http.ResponseWriter, r *http.Request) {
	var response payload

	slog.Info("Calling auth service")

	err := a.readJSON(w, r, &response)
	if err != nil {
		_ = a.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	log := Log{
		Name: response.Name,
		Data: response.Data,
	}

	err = a.Insert(log)
	if err != nil {
		_ = a.errorJSON(w, err)
		return
	}

	resp := jsonResponse{
		Error:   false,
		Message: "logged",
	}

	_ = a.writeJSON(w, http.StatusAccepted, resp)
}
