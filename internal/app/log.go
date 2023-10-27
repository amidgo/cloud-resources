package app

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

type LogOutput struct {
	serverUrl string
	client    *http.Client
}

func NewLogOutput(client *http.Client, serverUrl string) io.Writer {
	return &LogOutput{serverUrl: serverUrl, client: client}
}

func (l *LogOutput) Write(data []byte) (int, error) {
	resp, err := l.client.Post(l.serverUrl, "application/json", bytes.NewReader(data))
	if err != nil {
		return 0, fmt.Errorf("failed do request, %w", err)
	}
	if resp.StatusCode > 299 || resp.StatusCode < 200 {
		return 0, fmt.Errorf("failed do request, wrong status code, %w", err)
	}
	return len(data), nil
}
