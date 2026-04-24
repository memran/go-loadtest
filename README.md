# HTTP Load Testing Tool

A simple and powerful HTTP load testing tool written in Go. Perform concurrent HTTP load tests against any web endpoint with configurable parameters.

## Features

- ⚡ **High Performance**: Concurrent request handling with connection pooling
- 🎯 **Configurable**: Adjustable request count, concurrency levels, timeouts, and overall test timeout
- 📊 **Real-time Metrics**: Live progress tracking and detailed results
- 🔒 **Secure**: Optional TLS verification with clear warnings for insecure mode
- 🚀 **Fast**: Single binary with no external dependencies
- 💻 **Cross-platform**: Works on Windows, macOS, and Linux
- 📦 **Request Body Support**: Test POST/PUT/PATCH requests with custom JSON bodies
- 🔐 **Authentication**: Basic Auth and Bearer token support
- ✅ **Input Validation**: Validates URLs, HTTP methods, and parameters
- 📝 **Smart Logging**: Throttled logging to avoid output spam

## Quick Start

```bash
# Clone the repository
git clone https://github.com/memran/go-loadtest.git
cd go-loadtest

# Run a basic load test
go run main.go -url https://httpbin.org/get -requests 1000 -concurrency 50

# View all options
go run main.go --help
```

## Usage Examples

### Basic GET Request Test

```bash
go run main.go -url https://httpbin.org/get -requests 100 -concurrency 10
```

### High Concurrency Test

```bash
go run main.go -url https://api.example.com -requests 10000 -concurrency 100 -timeout 60
```

### POST Request with JSON Body

```bash
go run main.go -url https://httpbin.org/post -method POST \
  -body '{"name":"test","value":123}' -requests 100 -concurrency 10
```

### Test with Basic Authentication

```bash
go run main.go -url https://httpbin.org/basic-auth/user/pass \
  -auth-user user -auth-pass pass -requests 50 -concurrency 5
```

### Test with Bearer Token

```bash
go run main.go -url https://httpbin.org/bearer \
  -bearer-token "your-token-here" -requests 100 -concurrency 10
```

### Test with Overall Timeout

```bash
# Stop the entire test after 30 seconds, regardless of pending requests
go run main.go -url https://api.example.com -requests 10000 \
  -concurrency 100 -total-timeout 30
```

### Test with Insecure SSL (Development Only)

```bash
go run main.go -url https://self-signed.example.com -insecure \
  -requests 100 -concurrency 10
```

### Disable HTTP Keep-Alive

```bash
go run main.go -url https://example.com -keepalive=false -requests 100 -concurrency 10
```

## Command Line Options

| Option | Description | Default |
|--------|-------------|---------|
| `-url` | URL to test | `https://httpbin.org/get` |
| `-requests` | Number of requests | `100` |
| `-concurrency` | Number of concurrent requests | `10` |
| `-timeout` | Request timeout in seconds | `30` |
| `-total-timeout` | Total test timeout in seconds (0 for no timeout) | `0` |
| `-method` | HTTP method (GET, POST, PUT, DELETE, PATCH, HEAD) | `GET` |
| `-body` | Request body (for POST/PUT/PATCH) | `` |
| `-auth-user` | Username for basic authentication | `` |
| `-auth-pass` | Password for basic authentication | `` |
| `-bearer-token` | Bearer token for authentication | `` |
| `-keepalive` | Use HTTP keep-alive | `true` |
| `-insecure` | Skip TLS verification | `false` |

## Building from Source

```bash
# Clone the repository
git clone https://github.com/memran/go-loadtest.git
cd go-loadtest

# Build the binary
go build -o loadtest main.go

# Run the built binary
./loadtest -url https://api.example.com -requests 1000
```

## Requirements

- Go 1.19 or higher

## Output Example

```
🚀 Starting HTTP Load Test
📝 URL: https://api.example.com
🔢 Requests: 1000, Concurrency: 50
⏱️  Timeout: 30s, Method: GET
⏱️  Total Test Timeout: 300s

📊 Progress: 200/1000 (20.0%) - Rate: 100.5 req/sec
📊 Progress: 400/1000 (40.0%) - Rate: 102.3 req/sec
📊 Progress: 600/1000 (60.0%) - Rate: 98.7 req/sec
📊 Progress: 800/1000 (80.0%) - Rate: 101.2 req/sec
📊 Progress: 1000/1000 (100.0%) - Rate: 99.8 req/sec

============================================================
🎉 LOAD TEST COMPLETED!
============================================================
📍 URL: https://api.example.com
⏱️  Total Time: 10.02 seconds
📊 Total Requests: 1000
✅ Successful: 980
❌ Failed: 20
📈 Success Rate: 98.00%
🚀 Requests/Second: 99.80
🔀 Concurrency: 50
⚡ Method: GET
🔒 Keep-Alive: true
⏱️  Total Timeout: 300s
🔓 Insecure: false
============================================================
```

### Example with Authentication and Body

```
🚀 Starting HTTP Load Test
📝 URL: https://httpbin.org/post
🔢 Requests: 100, Concurrency: 10
⏱️  Timeout: 30s, Method: POST
📦 Body: {"test":"data"}
🔐 Auth: Basic (user)

✅ Request 0: Status 200 (1191.80ms)
✅ Request 1: Status 200 (793.49ms)
📊 Progress: 30/100 (30.0%) - Rate: 15.0 req/sec
...

============================================================
🎉 LOAD TEST COMPLETED!
============================================================
📍 URL: https://httpbin.org/post
⏱️  Total Time: 6.52 seconds
📊 Total Requests: 100
✅ Successful: 100
❌ Failed: 0
📈 Success Rate: 100.00%
🚀 Requests/Second: 15.34
🔀 Concurrency: 10
⚡ Method: POST
🔒 Keep-Alive: true
🔓 Insecure: false
============================================================
```

## Input Validation

The tool validates all inputs before starting the test:

- **URL**: Must start with `http://` or `https://`
- **HTTP Method**: Only allows GET, POST, PUT, DELETE, PATCH, HEAD
- **Requests/Concurrency/Timeout**: Must be positive integers
- **Total Timeout**: Must be non-negative
- **Body**: Cannot be used with GET, HEAD, or DELETE methods
- **Authentication**: Cannot use both Basic Auth and Bearer Token simultaneously

## Security ⚠️

**Important**: When using the `-insecure` flag, be aware that you're bypassing TLS certificate validation. Only use this option in trusted development environments. Never use it for production traffic or with sensitive data.

## Logging Strategy

The tool uses smart logging to avoid output spam:

- First 10 requests (success or failure) are always logged
- Successful requests: logged every 100th request after the first 10
- Failed requests (non-2xx): logged every 50th request
- Network errors: logged every 100th error

## License

This project is licensed under the MIT License.

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## Roadmap

- [x] Request body support
- [x] Authentication (Basic Auth, Bearer tokens)
- [ ] Custom headers
- [ ] JSON/CSV output formats
- [ ] HTML report generation
- [ ] Rate limiting
- [ ] More detailed metrics (latency percentiles, throughput graphs)
- [ ] Proxy support
- [ ] Cookie/session support
