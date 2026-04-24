# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a simple HTTP load testing tool written in Go. It allows you to perform concurrent HTTP load tests against any web endpoint with configurable parameters including request count, concurrency level, timeout, and HTTP method.

## Build and Run

The project is a single-file Go application with no external dependencies.

### Build and Run
```bash
# Build the application
go build -o loadtest main.go

# Run the load test
go run main.go -url https://api.example.com -requests 1000 -concurrency 50 -timeout 30

# Run with help
go run main.go --help
```

### Common Commands

```bash
# Basic load test
go run main.go -url https://httpbin.org/get -requests 100 -concurrency 10

# High concurrency test
go run main.go -url https://api.example.com -requests 10000 -concurrency 100 -timeout 60

# POST request test
go run main.go -url https://api.example.com/submit -method POST -requests 500 -concurrency 20

# Test with insecure SSL (for testing HTTPS endpoints with self-signed certs)
go run main.go -url https://self-signed.example.com -insecure -requests 100 -concurrency 10

# Disable keep-alive
go run main.go -url https://example.com -keepalive=false -requests 100 -concurrency 10
```

## Architecture

The application follows a simple structure:

1. **Configuration** (`Config` struct): Parses command-line flags to configure the test
2. **Metrics collection** (`Metrics` struct): Tracks request success/failure counts and timing
3. **HTTP client management** (`createHTTPClient`): Configures a reusable HTTP client with proper connection pooling
4. **Concurrency control**: Uses a semaphore pattern to limit concurrent requests
5. **Progress reporting**: Real-time progress updates every 2 seconds
6. **Request execution** (`makeRequest`): Individual HTTP request handling with response status tracking

### Key Design Patterns

- **Semaphore pattern**: Controls concurrency using buffered channels
- **Atomic counters**: Thread-safe metrics collection using `sync/atomic`
- **Progress reporting**: Separate goroutine for periodic progress updates
- **Connection pooling**: Reuses HTTP connections for performance
- **Graceful shutdown**: Proper resource cleanup with `defer` and `WaitGroup`

### Important Implementation Details

- The application uses a custom string repeat function to avoid importing the `strings` package (demonstrating minimal dependency approach)
- HTTP redirects are disabled by default (`http.ErrUseLastResponse`)
- Response bodies are discarded but fully read to ensure proper connection reuse
- Error logging is throttled to prevent spam (every 100th error, every 50th failure, every 100th success)
- The application tracks both success (2xx) and failure (non-2xx) status codes

## Testing

No unit tests are currently included. When adding tests, consider:
- Testing the flag parsing logic
- Mocking HTTP requests for request testing
- Testing metrics collection accuracy
- Testing concurrency control mechanisms