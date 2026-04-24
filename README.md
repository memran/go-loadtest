# HTTP Load Testing Tool

A simple and powerful HTTP load testing tool written in Go. Perform concurrent HTTP load tests against any web endpoint with configurable parameters.

## Features

- ⚡ **High Performance**: Concurrent request handling with connection pooling
- 🎯 **Configurable**: Adjustable request count, concurrency levels, and timeouts
- 📊 **Real-time Metrics**: Live progress tracking and detailed results
- 🔒 **Secure**: Optional TLS verification with clear warnings for insecure mode
- 🚀 **Fast**: Single binary with no external dependencies
- 💻 **Cross-platform**: Works on Windows, macOS, and Linux

## Quick Start

```bash
# Clone the repository
git clone https://github.com/yourusername/go-http-loadtest.git
cd go-http-loadtest

# Run a basic load test
go run main.go -url https://api.example.com -requests 1000 -concurrency 50

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

### POST Request Test

```bash
go run main.go -url https://api.example.com/submit -method POST -requests 500 -concurrency 20
```

### Test with Insecure SSL (Development Only)

```bash
go run main.go -url https://self-signed.example.com -insecure -requests 100 -concurrency 10
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
| `-method` | HTTP method | `GET` |
| `-keepalive` | Use HTTP keep-alive | `true` |
| `-insecure` | Skip TLS verification | `false` |

## Building from Source

```bash
# Clone the repository
git clone https://github.com/yourusername/go-http-loadtest.git
cd go-http-loadtest

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
🔓 Insecure: false
============================================================
```

## Security ⚠️

**Important**: When using the `-insecure` flag, be aware that you're bypassing TLS certificate validation. Only use this option in trusted development environments. Never use it for production traffic or with sensitive data.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## Roadmap

- [ ] Request body support
- [ ] Authentication (Basic Auth, API keys, Bearer tokens)
- [ ] Custom headers
- [ ] JSON/CSV output formats
- [ ] HTML report generation
- [ ] Rate limiting
- [ ] More detailed metrics (latency percentiles, throughput graphs)