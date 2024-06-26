package main

import (
	"io/ioutil"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	logger = slog.New(slog.NewJSONHandler(ioutil.Discard, nil))
	exitCode := m.Run()
	os.Exit(exitCode)
}

func TestRunHealthcheck(t *testing.T) {
	tests := []struct {
		name       string
		method     string
		statusCode int
		timeout    time.Duration
		expected   int
	}{
		{"Valid GET request", "GET", http.StatusOK, 5, exitCodeSuccess},
		{"Valid POST request", "POST", http.StatusOK, 5, exitCodeSuccess},
		{"Invalid method", "INVALID", http.StatusMethodNotAllowed, 5, exitCodeError},
		{"Timeout exceeded", "GET", http.StatusOK, 1, exitCodeError},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != test.method {
					w.WriteHeader(http.StatusMethodNotAllowed)
				} else {
					time.Sleep(2 * time.Second)
					w.WriteHeader(test.statusCode)
				}
			}))
			defer ts.Close()
			result := runHealthcheck(test.method, ts.URL, test.timeout)
			if result != test.expected {
				t.Errorf("expected %v, got %v", test.expected, result)
			}
		})
	}
}
