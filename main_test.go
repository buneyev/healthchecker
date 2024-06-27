package main

import (
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestRunHealthcheck(t *testing.T) {
	logger = slog.New(slog.NewJSONHandler(io.Discard, nil))

	tests := []struct {
		name              string
		method            string
		url               string
		timeout           int
		serverStatusCode  int
		expectedExitCode  int
		expectedLogPrefix string
	}{
		{
			name:              "Valid GET request",
			method:            "GET",
			url:               "/health",
			timeout:           5,
			serverStatusCode:  http.StatusOK,
			expectedExitCode:  exitCodeSuccess,
			expectedLogPrefix: "INFO: Health check /health: HTTP status code",
		},
		{
			name:              "Invalid timeout",
			method:            "GET",
			url:               "/health",
			timeout:           1,
			serverStatusCode:  http.StatusOK,
			expectedExitCode:  exitCodeError,
			expectedLogPrefix: "ERROR: The timeout must be between 0 and 360 seconds",
		},
		{
			name:              "Invalid URL",
			method:            "GET",
			url:               "invalid-url",
			timeout:           5,
			serverStatusCode:  http.StatusOK,
			expectedExitCode:  exitCodeError,
			expectedLogPrefix: "ERROR: The health check URL must be valid",
		},
		{
			name:              "Invalid HTTP method",
			method:            "INVALID",
			url:               "/health",
			timeout:           5,
			serverStatusCode:  http.StatusOK,
			expectedExitCode:  exitCodeError,
			expectedLogPrefix: "ERROR: The method must be a valid HTTP method",
		},
		{
			name:              "Server error response",
			method:            "GET",
			url:               "/health",
			timeout:           5,
			serverStatusCode:  http.StatusInternalServerError,
			expectedExitCode:  exitCodeError,
			expectedLogPrefix: "ERROR: Health check /health: Error - status code 500",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				time.Sleep(2 * time.Second)
				w.WriteHeader(tt.serverStatusCode)
			}))
			defer server.Close()

			if tt.url != "invalid-url" {
				tt.url = server.URL + tt.url
			}

			exitCode := runHealthcheck(tt.method, tt.url, tt.timeout)
			if exitCode != tt.expectedExitCode {
				t.Errorf("expected exit code %d, got %d", tt.expectedExitCode, exitCode)
			}
		})
	}
}
