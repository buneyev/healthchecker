# Health Check Application

This Go application performs a health check on a specified URL using a specified HTTP method. It uses the `slog` package for logging and allows for setting the log level through a flag.

## Features

- Perform HTTP health checks with various methods (GET, POST, etc.)
- Configurable timeout for the health check requests
- Logging with different log levels (debug, info, warn, error)

## Usage

### Flags

- `-mode`: Health check mode. Currently, only "http" is available.
- `-url`: URL to check the health of the app. Default is `http://localhost:8080/healthz`.
- `-method`: HTTP method to use for the health check. Default is `GET`.
- `-timeout`: Timeout (in seconds) for the health check request. Minimum 0, maximum 360. Default is 5.
- `-loglevel`: Log level for the application. Options are `debug`, `info`, `warn`, `error`. Default is `error`.

### Example

```sh
go run main.go -url http://example.com/health -method GET -timeout 10 -loglevel info
```