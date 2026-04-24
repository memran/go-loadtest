package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

type Config struct {
	URL         string
	Requests    int
	Concurrency int
	Timeout     int
	Method      string
	KeepAlive   bool
	Insecure    bool
}

type Metrics struct {
	TotalRequests int64
	SuccessCount  int64
	FailCount     int64
	TotalDuration time.Duration
	StartTime     time.Time
}

func main() {
	config := parseFlags()
	
	fmt.Printf("🚀 Starting HTTP Load Test\n")
	fmt.Printf("📝 URL: %s\n", config.URL)
	fmt.Printf("🔢 Requests: %d, Concurrency: %d\n", config.Requests, config.Concurrency)
	fmt.Printf("⏱️  Timeout: %ds, Method: %s\n\n", config.Timeout, config.Method)

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
	flag.StringVar(&config.Method, "method", "GET", "HTTP method")
	flag.BoolVar(&config.KeepAlive, "keepalive", true, "Use HTTP keep-alive")
	flag.BoolVar(&config.Insecure, "insecure", false, "Skip TLS verification")

	flag.Parse()
	return config
}

func runLoadTest(config *Config, metrics *Metrics) {
	// Create HTTP client
	client := createHTTPClient(config)

	// Semaphore for controlling concurrency
	semaphore := make(chan struct{}, config.Concurrency)

	// Wait group to wait for all requests
	var wg sync.WaitGroup

	// Progress tracking
	var counter int64

	// Start progress reporter
	go reportProgress(metrics, config.Requests, &counter)

	// Send requests
	for i := 0; i < config.Requests; i++ {
		wg.Add(1)
		semaphore <- struct{}{} // Acquire semaphore

		go func(requestID int) {
			defer wg.Done()
			defer func() { <-semaphore }() // Release semaphore

			atomic.AddInt64(&counter, 1)
			makeRequest(client, config, metrics, requestID)
		}(i)
	}

	// Wait for all requests to complete
	wg.Wait()
	metrics.TotalDuration = time.Since(metrics.StartTime)
}

func createHTTPClient(config *Config) *http.Client {
	transport := &http.Transport{
		MaxIdleConns:        config.Concurrency,
		MaxIdleConnsPerHost: config.Concurrency,
		IdleConnTimeout:     90 * time.Second,
		DisableKeepAlives:   !config.KeepAlive,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: config.Insecure,
		},
	}

	return &http.Client{
		Transport: transport,
		Timeout:   time.Duration(config.Timeout) * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse // Don't follow redirects
		},
	}
}

func makeRequest(client *http.Client, config *Config, metrics *Metrics, requestID int) {
	startTime := time.Now()

	// Create request
	req, err := http.NewRequest(config.Method, config.URL, nil)
	if err != nil {
		log.Printf("❌ Request %d: Failed to create request: %v", requestID, err)
		atomic.AddInt64(&metrics.FailCount, 1)
		return
	}

	// Set headers
	req.Header.Set("User-Agent", "Go-HTTP-Load-Tester/1.0")
	req.Header.Set("Accept", "*/*")

	// Make request
	resp, err := client.Do(req)
	duration := time.Since(startTime)

	if err != nil {
		if requestID%100 == 0 { // Log only some errors to avoid spam
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
		if requestID%100 == 0 { // Log only some successes
			log.Printf("✅ Request %d: Status %d (%.2fms)", requestID, resp.StatusCode, duration.Seconds()*1000)
		}
	} else {
		atomic.AddInt64(&metrics.FailCount, 1)
		if requestID%50 == 0 { // Log more failed requests
			log.Printf("⚠️ Request %d: Status %d (%.2fms)", requestID, resp.StatusCode, duration.Seconds()*1000)
		}
	}
}

func reportProgress(metrics *Metrics, total int, counter *int64) {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		current := atomic.LoadInt64(counter)
		if current >= int64(total) {
			break
		}

		elapsed := time.Since(metrics.StartTime).Seconds()
		rate := float64(current) / elapsed
		
		fmt.Printf("📊 Progress: %d/%d (%.1f%%) - Rate: %.1f req/sec\n", 
			current, total, float64(current)/float64(total)*100, rate)
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
	fmt.Printf("📈 Success Rate: %.2f%%\n", float64(success)/float64(total)*100)
	fmt.Printf("🚀 Requests/Second: %.2f\n", float64(total)/metrics.TotalDuration.Seconds())
	fmt.Printf("🔀 Concurrency: %d\n", config.Concurrency)
	fmt.Printf("⚡ Method: %s\n", config.Method)
	fmt.Printf("🔒 Keep-Alive: %v\n", config.KeepAlive)
	fmt.Printf("🔓 Insecure: %v\n", config.Insecure)
	fmt.Printf("%s\n", strings.Repeat("=", 60))
}

// Helper string repeat function (since we can't import strings in the basic example)
var strings = struct {
	Repeat func(string, int) string
}{
	Repeat: func(s string, count int) string {
		var result string
		for i := 0; i < count; i++ {
			result += s
		}
		return result
	},
}