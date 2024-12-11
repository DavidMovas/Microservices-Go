package api

import (
	"net/http"
)

func (a *App) Broker(w http.ResponseWriter, r *http.Request) {
	payload := jsonResponse{
		Error:   false,
		Message: "Hit the broker",
	}

	_ = a.writeJSON(w, http.StatusOK, payload)
}
