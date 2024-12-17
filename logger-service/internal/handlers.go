package internal

import (
	"github.com/sirupsen/logrus"
	"net/http"
)

type payload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (a *App) Logger(w http.ResponseWriter, r *http.Request) {
	var response payload

	logrus.Info("Calling auth service")

	err := a.readJSON(w, r, &response)
	if err != nil {
		_ = a.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	log := Log{
		Name: response.Name,
		Data: response.Data,
	}

	logrus.Info("Logging", "log", log)

	resp := jsonResponse{
		Error:   false,
		Message: "logged",
	}

	_ = a.writeJSON(w, http.StatusAccepted, resp)
}
