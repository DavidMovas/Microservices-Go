package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type LokiHook struct {
	Host     string
	Labels   map[string]string
	Username string
	Password string
}

func NewLokiHook(cfg *LokiConfig) logrus.Hook {
	return &LokiHook{
		Host:     cfg.Host,
		Labels:   cfg.Labels,
		Username: cfg.Username,
		Password: cfg.Password,
	}
}

func (l *LokiHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (l *LokiHook) Fire(entry *logrus.Entry) error {
	timestamp := fmt.Sprintf("%d", time.Now().Unix())

	logEntry := map[string]any{
		"streams": []map[string]any{
			{
				"stream": l.Labels,
				"values": [][]string{
					{timestamp, entry.Message},
				},
			},
		},
	}

	jsonPayload, err := json.Marshal(logEntry)
	if err != nil {
		return err
	}

	resp, err := http.Post(l.Host, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return err
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf("error sending log to loki: %s", resp.Status)
	}

	return nil
}
