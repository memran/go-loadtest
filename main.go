package main

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type Config struct {
	URL          string
	Requests     int
	Concurrency  int
	Timeout      int
	TotalTimeout int
	Method       string
	Body         string
	AuthUser     string
	AuthPass     string
	BearerToken  string
	KeepAlive    bool
	Insecure     bool
}

type Metrics struct {
	SuccessCount  int64
	FailCount     int64
	TotalDuration time.Duration
	StartTime     time.Time
	loggedCount   int64 // Track how many requests we've logged
}

func main() {
	config := parseFlags()

	fmt.Printf("🚀 Starting HTTP Load Test\n")
	fmt.Printf("📝 URL: %s\n", config.URL)
	fmt.Printf("🔢 Requests: %d, Concurrency: %d\n", config.Requests, config.Concurrency)
	fmt.Printf("⏱️  Timeout: %ds, Method: %s\n", config.Timeout, config.Method)
	if config.TotalTimeout > 0 {
		fmt.Printf("⏱️  Total Test Timeout: %ds\n", config.TotalTimeout)
	}
	if config.Body != "" {
		fmt.Printf("📦 Body: %s\n", config.Body)
	}
	if config.AuthUser != "" {
		fmt.Printf("🔐 Auth: Basic (%s)\n", config.AuthUser)
	}
	if config.BearerToken != "" {
		fmt.Printf("🔐 Auth: Bearer token\n")
	}
	fmt.Println()

	metrics := &Metrics{
		StartTime: time.Now(),
	}

	// Run the load test
	runLoadTest(config, metrics)

	// Print final results
	printResults(config, metrics)
}

func parseFlags() *Config {
	config := &Config{}

	flag.StringVar(&config.URL, "url", "https://httpbin.org/get", "URL to test")
	flag.IntVar(&config.Requests, "requests", 100, "Number of requests")
	flag.IntVar(&config.Concurrency, "concurrency", 10, "Number of concurrent requests")
	flag.IntVar(&config.Timeout, "timeout", 30, "Request timeout in seconds")
	flag.IntVar(&config.TotalTimeout, "total-timeout", 0, "Total test timeout in seconds (0 for no timeout)")
	flag.StringVar(&config.Method, "method", "GET", "HTTP method")
	flag.StringVar(&config.Body, "body", "", "Request body (for POST/PUT/PATCH)")
	flag.StringVar(&config.AuthUser, "auth-user", "", "Username for basic authentication")
	flag.StringVar(&config.AuthPass, "auth-pass", "", "Password for basic authentication")
	flag.StringVar(&config.BearerToken, "bearer-token", "", "Bearer token for authentication")
	flag.BoolVar(&config.KeepAlive, "keepalive", true, "Use HTTP keep-alive")
	flag.BoolVar(&config.Insecure, "insecure", false, "Skip TLS verification")

	flag.Parse()

	// Validate and normalize inputs
	config.Method = strings.ToUpper(config.Method)

	// Validate HTTP method
	validMethods := map[string]bool{
		"GET": true, "POST": true, "PUT": true, "DELETE": true, "PATCH": true, "HEAD": true,
	}
	if !validMethods[config.Method] {
		log.Fatalf("Invalid HTTP method: %s. Allowed methods: GET, POST, PUT, DELETE, PATCH, HEAD", config.Method)
	}

	// Validate URL
	if config.URL == "" {
		log.Fatal("URL is required")
	}
	if !strings.HasPrefix(config.URL, "http://") && !strings.HasPrefix(config.URL, "https://") {
		log.Fatalf("Invalid URL: %s. URL must start with http:// or https://", config.URL)
	}

	// Validate numeric parameters
	if config.Requests <= 0 {
		log.Fatal("Requests must be a positive integer")
	}
	if config.Concurrency <= 0 {
		log.Fatal("Concurrency must be a positive integer")
	}
	if config.Timeout <= 0 {
		log.Fatal("Request timeout must be a positive integer")
	}
	if config.TotalTimeout < 0 {
		log.Fatal("Total timeout must be non-negative")
	}

	// Validate body is only used with appropriate methods
	if config.Body != "" && (config.Method == "GET" || config.Method == "HEAD" || config.Method == "DELETE") {
		log.Fatalf("Request body is not typically used with %s method", config.Method)
	}

	// Validate authentication options
	if config.AuthUser != "" && config.BearerToken != "" {
		log.Fatal("Cannot use both basic auth and bearer token authentication")
	}

	return config
}

func runLoadTest(config *Config, metrics *Metrics) {
	// Create context for overall timeout
	ctx := context.Background()
	var cancel context.CancelFunc
	if config.TotalTimeout > 0 {
		ctx, cancel = context.WithTimeout(context.Background(), time.Duration(config.TotalTimeout)*time.Second)
	} else {
		ctx, cancel = context.WithCancel(context.Background())
	}
	defer cancel()

	// Create HTTP client and transport
	client, transport := createHTTPClient(config)
	defer transport.CloseIdleConnections()

	// Semaphore for controlling concurrency
	semaphore := make(chan struct{}, config.Concurrency)

	// Wait group to wait for all requests
	var wg sync.WaitGroup

	// Progress tracking
	var counter int64

	// Start progress reporter
	done := make(chan struct{})
	defer close(done)
	go reportProgress(metrics, config.Requests, &counter, ctx, done)

	// Send requests
	for i := 0; i < config.Requests; i++ {
		// Check if context is done before starting new request
		if ctx.Err() != nil {
			break
		}

		wg.Add(1)
		semaphore <- struct{}{} // Acquire semaphore

		go func(requestID int) {
			defer wg.Done()
			defer func() { <-semaphore }() // Release semaphore

			// Check if context is done before making request
			if ctx.Err() != nil {
				atomic.AddInt64(&metrics.FailCount, 1)
				return
			}

			atomic.AddInt64(&counter, 1)
			makeRequest(client, config, metrics, requestID)
		}(i)
	}

	// Wait for all requests to complete or context to timeout
	completed := make(chan struct{})
	go func() {
		wg.Wait()
		close(completed)
	}()

	select {
	case <-completed:
		// All requests completed
	case <-ctx.Done():
		log.Println("🛑 Overall test timeout reached")
	}

	metrics.TotalDuration = time.Since(metrics.StartTime)
}

func createHTTPClient(config *Config) (*http.Client, *http.Transport) {
	transport := &http.Transport{
		MaxIdleConns:        config.Concurrency,
		MaxIdleConnsPerHost: config.Concurrency,
		IdleConnTimeout:     90 * time.Second,
		DisableKeepAlives:   !config.KeepAlive,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: config.Insecure,
		},
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   time.Duration(config.Timeout) * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse // Don't follow redirects
		},
	}

	return client, transport
}

func makeRequest(client *http.Client, config *Config, metrics *Metrics, requestID int) {
	startTime := time.Now()

	// Create request body if provided
	var body io.Reader
	if config.Body != "" {
		body = strings.NewReader(config.Body)
	}

	// Create request
	req, err := http.NewRequest(config.Method, config.URL, body)
	if err != nil {
		// Log first 10 request creation errors
		if atomic.AddInt64(&metrics.loggedCount, 1) <= 10 {
			log.Printf("❌ Request %d: Failed to create request: %v", requestID, err)
		}
		atomic.AddInt64(&metrics.FailCount, 1)
		return
	}

	// Set headers
	req.Header.Set("User-Agent", "Go-HTTP-Load-Tester/1.0")
	req.Header.Set("Accept", "*/*")
	if config.Body != "" && req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", "application/json")
	}

	// Add authentication
	if config.AuthUser != "" {
		req.SetBasicAuth(config.AuthUser, config.AuthPass)
	}
	if config.BearerToken != "" {
		req.Header.Set("Authorization", "Bearer "+config.BearerToken)
	}

	// Make request
	resp, err := client.Do(req)
	duration := time.Since(startTime)

	if err != nil {
		// Log errors: first 10, then every 100th
		logCount := atomic.AddInt64(&metrics.loggedCount, 1)
		if logCount <= 10 || logCount%100 == 0 {
			log.Printf("❌ Request %d: Error: %v (%.2fms)", requestID, err, duration.Seconds()*1000)
		}
		atomic.AddInt64(&metrics.FailCount, 1)
		return
	}
	defer resp.Body.Close()

	// Read response body (but discard it)
	_, _ = io.Copy(io.Discard, resp.Body)

	// Check if successful
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		atomic.AddInt64(&metrics.SuccessCount, 1)
		// Log successes: first 10, then every 100th
		logCount := atomic.LoadInt64(&metrics.SuccessCount)
		if logCount <= 10 || logCount%100 == 0 {
			log.Printf("✅ Request %d: Status %d (%.2fms)", requestID, resp.StatusCode, duration.Seconds()*1000)
		}
	} else {
		atomic.AddInt64(&metrics.FailCount, 1)
		// Log failures: first 10, then every 50th
		failCount := atomic.LoadInt64(&metrics.FailCount)
		if failCount <= 10 || failCount%50 == 0 {
			log.Printf("⚠️ Request %d: Status %d (%.2fms)", requestID, resp.StatusCode, duration.Seconds()*1000)
		}
	}
}

func reportProgress(metrics *Metrics, total int, counter *int64, ctx context.Context, done <-chan struct{}) {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return
		case <-ctx.Done():
			return
		case <-ticker.C:
			current := atomic.LoadInt64(counter)
			if current >= int64(total) {
				return
			}

			elapsed := time.Since(metrics.StartTime).Seconds()
			rate := float64(current) / elapsed

			fmt.Printf("📊 Progress: %d/%d (%.1f%%) - Rate: %.1f req/sec\n",
				current, total, float64(current)/float64(total)*100, rate)
		}
	}
}

func printResults(config *Config, metrics *Metrics) {
	total := atomic.LoadInt64(&metrics.SuccessCount) + atomic.LoadInt64(&metrics.FailCount)
	success := atomic.LoadInt64(&metrics.SuccessCount)
	fail := atomic.LoadInt64(&metrics.FailCount)

	fmt.Printf("\n%s\n", strings.Repeat("=", 60))
	fmt.Println("🎉 LOAD TEST COMPLETED!")
	fmt.Printf("%s\n", strings.Repeat("=", 60))
	fmt.Printf("📍 URL: %s\n", config.URL)
	fmt.Printf("⏱️  Total Time: %.2f seconds\n", metrics.TotalDuration.Seconds())
	fmt.Printf("📊 Total Requests: %d\n", total)
	fmt.Printf("✅ Successful: %d\n", success)
	fmt.Printf("❌ Failed: %d\n", fail)

	// Calculate success rate safely
	var successRate float64
	if total > 0 {
		successRate = float64(success) / float64(total) * 100
	}
	fmt.Printf("📈 Success Rate: %.2f%%\n", successRate)

	fmt.Printf("🚀 Requests/Second: %.2f\n", float64(total)/metrics.TotalDuration.Seconds())
	fmt.Printf("🔀 Concurrency: %d\n", config.Concurrency)
	fmt.Printf("⚡ Method: %s\n", config.Method)
	fmt.Printf("🔒 Keep-Alive: %v\n", config.KeepAlive)
	if config.TotalTimeout > 0 {
		fmt.Printf("⏱️  Total Timeout: %ds\n", config.TotalTimeout)
	}
	fmt.Printf("🔓 Insecure: %v\n", config.Insecure)
	fmt.Printf("%s\n", strings.Repeat("=", 60))
}

