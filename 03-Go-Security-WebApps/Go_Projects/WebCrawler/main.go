package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

// scanWorker processes URLs at a controlled pace mandated by the ticker channel
func scanWorker(baseURL string, jobs <-chan string, throttle <-chan time.Time, wg *sync.WaitGroup) {
	defer wg.Done()

	// Initialize an advanced HTTP client with strict timeout controls
	client := &http.Client{
		Timeout: 4 * time.Second,
		// Prevent the client from automatically following redirects (301/302).
		// In cybersecurity assessments, we want to know if a folder redirects us elsewhere.
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	for word := range jobs {
		// Wait for the master time ticker to grant execution clearance
		<-throttle

		// Construct the absolute path target
		targetURL := fmt.Sprintf("%s/%s", baseURL, word)

		// Initialize a new custom request object
		req, err := http.NewRequest("GET", targetURL, nil)
		if err != nil {
			continue
		}

		// Inject realistic browser headers to bypass primitive signature blocks
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
		req.Header.Set("Accept", "*/*")

		// Execute the network connection request
		resp, err := client.Do(req)
		if err != nil {
			// Remote server dropped connection or timed out
			continue
		}

		// Clean up system sockets immediately after reading the response status
		resp.Body.Close()

		// Analyze status classifications
		switch {
		case resp.StatusCode == http.StatusOK: // 200 OK
			fmt.Printf("[+] [200 OK] Discovered Entrypoint: %s\n", targetURL)
		case resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden: // 401 / 403
			fmt.Printf("[!] [%d Restricted] Guarded Assets Found: %s\n", resp.StatusCode, targetURL)
		case resp.StatusCode >= 300 && resp.StatusCode < 400: // Redirects
			fmt.Printf("[*] [%d Redirect] Points Elsewhere: %s\n", resp.StatusCode, targetURL)
		}
	}
}

func main() {
	// 1. Argument and Input Verification
	if len(os.Args) < 5 {
		fmt.Println("Usage: go run main.go <target_url> <wordlist_path> <threads> <delay_in_ms>")
		fmt.Println("Example: go run main.go https://example.com words.txt 5 200")
		os.Exit(1)
	}

	targetURL := strings.TrimSuffix(os.Args[1], "/")
	wordlistPath := os.Args[2]

	var threads int
	var delayMs int
	_, _ = fmt.Sscanf(os.Args[3], "%d", &threads)
	_, _ = fmt.Sscanf(os.Args[4], "%d", &delayMs)

	if threads <= 0 || delayMs < 0 {
		fmt.Println("[!] Error: Threads and delays must be valid positive calculations.")
		os.Exit(1)
	}

	// 2. Open wordlist file safely
	file, err := os.Open(wordlistPath)
	if err != nil {
		fmt.Printf("[!] Wordlist extraction failure: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	// 3. Set up the structural rate limiter (The Pacing Engine)
	// time.NewTicker pulses a signal onto a channel at an absolute interval
	ticker := time.NewTicker(time.Duration(delayMs) * time.Millisecond)
	defer ticker.Stop()

	jobs := make(chan string, 500)
	var wg sync.WaitGroup

	fmt.Printf("[*] Initializing Paced Web Scan against: %s\n", targetURL)
	fmt.Printf("[*] Throttling Metric: %d active threads running with a mandatory %dms buffer interval...\n\n", threads, delayMs)

	// 4. Spin up our synchronized worker pool
	for i := 0; i < threads; i++ {
		wg.Add(1)
		// We pass the ticker's communication channel (ticker.C) directly into our workers
		go scanWorker(targetURL, jobs, ticker.C, &wg)
	}

	// 5. Asynchronously feed words from our text file into the processing pipeline
	go func() {
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			word := strings.TrimSpace(scanner.Text())
			if word != "" && !strings.HasPrefix(word, "#") {
				jobs <- word
			}
		}
		close(jobs)
	}()

	// 6. Complete lifecycle block management
	wg.Wait()
	fmt.Println("\n[*] Directory traversal sequence complete.")
}
