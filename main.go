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
	logger             *slog.Logger
	mode               string
	logLevel           string
	healthcheckURL     string
	healthcheckMethod  string
	healthcheckTimeout int
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

func runHealthcheck(healthcheckMethod, helthcheckUrl string, timeout int) int {
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
		Timeout: time.Duration(timeout) * time.Second,
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

	if resp.StatusCode > 299 {
		logger.Error(fmt.Sprintf("Health check %s: Error - status code %d", helthcheckUrl, resp.StatusCode))
		return exitCodeError
	}

	logger.Info(fmt.Sprintf("Health check %s: HTTP status code %s", helthcheckUrl, resp.Status))
	return exitCodeSuccess
}

func init() {
	flag.StringVar(&mode, "mode", "http", "Healt check mode. At now only http available")
	flag.StringVar(&healthcheckURL, "url", "http://localhost:8080/healthz", "Url for check health of app")
	flag.StringVar(&healthcheckMethod, "method", "GET", "Method of http request for check health")
	flag.IntVar(&healthcheckTimeout, "timeout", 5, "Timeout (in seconds) for health check request to app. Minimum 0, maximun 360")
	flag.StringVar(&logLevel, "loglevel", "error", "Log level: debug, info, warn, error")
}

func main() {
	flag.Parse()
	logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: parseLogLevel(logLevel)}))

	if mode != "http" {
		log.Fatal("Only HTTP mode is available for the health check")
	}

	logger.Info(fmt.Sprintf("Starting health check at %v", healthcheckURL))

	exitcode := runHealthcheck(healthcheckMethod, healthcheckURL, healthcheckTimeout)
	os.Exit(exitcode)
}

func parseLogLevel(logLevel string) slog.Level {
	var logLevelOption slog.Level
	switch strings.ToLower(logLevel) {
	case "debug":
		logLevelOption = slog.LevelDebug
	case "info":
		logLevelOption = slog.LevelInfo
	case "warn":
		logLevelOption = slog.LevelWarn
	case "error":
		logLevelOption = slog.LevelError
	default:
		log.Fatalf("Invalid log level: %s", logLevel)
	}
	return logLevelOption
}
