package internal

import (
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type payload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (a *App) Logger(w http.ResponseWriter, r *http.Request) {
	var response payload

	logrus.Info("Calling logger service")

	err := a.readJSON(w, r, &response)
	if err != nil {
		_ = a.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	log := Log{
		Name: response.Name,
		Data: response.Data,
	}

	logrus.WithFields(logrus.Fields{
		"name": log.Name,
		"data": log.Data,
		"time": time.Now().String(),
	}).Info("logged message")

	resp := jsonResponse{
		Error:   false,
		Message: "logged",
	}

	_ = a.writeJSON(w, http.StatusAccepted, resp)
}
