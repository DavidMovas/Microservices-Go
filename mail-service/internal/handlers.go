package internal

import (
	"log/slog"
	"net/http"
)

func (a *App) SendMail(w http.ResponseWriter, r *http.Request) {
	type mailPayload struct {
		From    string `json:"from"`
		To      string `json:"to"`
		Subject string `json:"subject"`
		Message string `json:"message"`
	}

	var requestPayload mailPayload

	err := a.readJSON(w, r, &requestPayload)
	if err != nil {
		_ = a.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	msg := Message{
		From:    requestPayload.From,
		To:      requestPayload.To,
		Subject: requestPayload.Subject,
		Data:    requestPayload.Message,
	}

	if err = a.mailer.SendSMTPMessage(msg); err != nil {
		slog.Error("Error sending mail", "ERROR", err)
		_ = a.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	slog.Info("Mail sent", "to", requestPayload.To)

	payload := jsonResponse{
		Error:   false,
		Message: "sent to " + requestPayload.To,
	}

	_ = a.writeJSON(w, http.StatusAccepted, payload)
}
