package main

import (
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

const (
	exitCodeError   = 1
	exitCodeSuccess = 0
)

var (
	logger             = slog.New(slog.NewJSONHandler(os.Stdout, nil))
	mode               string
	healthcheckURL     string
	healthcheckMethod  string
	healthcheckTimeout time.Duration
	validMethods       = map[string]struct{}{
		"GET":     {},
		"POST":    {},
		"PUT":     {},
		"DELETE":  {},
		"PATCH":   {},
		"HEAD":    {},
		"OPTIONS": {},
		"CONNECT": {},
		"TRACE":   {}}
)

func runHealthcheck(healthcheckMethod, helthcheckUrl string, timeout time.Duration) int {
	logger.Info(fmt.Sprintf("Querying Endpoint %v", helthcheckUrl))

	if healthcheckTimeout < 0 || healthcheckTimeout > 360 {
		logger.Error("The timeout must be between 0 and 360 seconds")
		return exitCodeError
	}

	_, err := url.ParseRequestURI(helthcheckUrl)
	if err != nil {
		logger.Error("The health check URL must be valid")
		return exitCodeError
	}

	if _, ok := validMethods[strings.ToUpper(healthcheckMethod)]; !ok {
		logger.Error("The method must be a valid HTTP method")
		return exitCodeError
	}

	client := &http.Client{
		Timeout: timeout * time.Second,
	}
	r, err := http.NewRequest(healthcheckMethod, helthcheckUrl, nil)
	if err != nil {
		logger.Error(fmt.Sprintf("error creating healthcheck request: %v", err))
		return exitCodeError
	}

	resp, err := client.Do(r)
	if err != nil {
		logger.Error(fmt.Sprintf("Health check %s: Error %s", helthcheckUrl, err))
		return exitCodeError
	}
	defer resp.Body.Close()

	logger.Info(fmt.Sprintf("Health check %s: HTTP status code %s", helthcheckUrl, resp.Status))
	return exitCodeSuccess
}

func init() {
	flag.StringVar(&mode, "mode", "http", "Healt check mode. At now only http available")
	flag.StringVar(&healthcheckURL, "url", "http://localhost:8080/health", "Url for check health of app")
	flag.StringVar(&healthcheckMethod, "method", "HEAD", "Method of http request for check health")
	flag.DurationVar(&healthcheckTimeout, "timeout", 5, "Timeout (in seconds) for health check request to app. Minimum 0, maximun 360")
}

func main() {
	flag.Parse()
	if mode != "http" {
		log.Fatal("Only HTTP mode is available for the health check")
	}

	logger.Info(fmt.Sprintf("Starting health check at %v", healthcheckURL))

	exitcode := runHealthcheck(healthcheckMethod, healthcheckURL, healthcheckTimeout)
	os.Exit(exitcode)
}
