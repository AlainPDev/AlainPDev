package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

// Result structural layout holds our discovered asset data
type Result struct {
	Subdomain string
	IPs       []string
}

// worker acts as an independent execution thread.
// It reads from the jobs channel, resolves the DNS records, and pushes findings to the results channel.
func worker(domain string, jobs <-chan string, results chan<- Result, wg *sync.WaitGroup) {
	defer wg.Done()

	for sub := range jobs {
		// Construct the absolute Fully Qualified Domain Name (FQDN)
		fqdn := fmt.Sprintf("%s.%s", sub, domain)

		// Perform an actual system network lookup against the global DNS hierarchy
		ips, err := net.LookupHost(fqdn)
		if err != nil {
			// DNS entry does not exist or failed to resolve
			continue
		}

		// If resolution succeeds, ship the structured findings to the collector
		results <- Result{
			Subdomain: fqdn,
			IPs:       ips,
		}
	}
}

func main() {
	// 1. Argument Validation
	if len(os.Args) < 4 {
		fmt.Println("Usage: go run main.go <target_domain> <wordlist_path> <concurrency_threads>")
		fmt.Println("Example: go run main.go example.com subdomains.txt 50")
		os.Exit(1)
	}

	targetDomain := os.Args[1]
	wordlistPath := os.Args[2]

	// Convert the third argument from a string into an integer
	workerCount, err := fmt.Sscanf(os.Args[3], "%d")
	if err != nil || workerCount <= 0 {
		fmt.Println("[!] Error: Concurrency threads must be a valid positive integer.")
		os.Exit(1)
	}

	// 2. Open the wordlist file safely
	file, err := os.Open(wordlistPath)
	if err != nil {
		fmt.Printf("[!] Failed to open wordlist: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	// 3. Initialize Worker Pool Channels and Synchronization Groups
	jobs := make(chan string, 1000)    // Buffered pipeline feeding words to workers
	results := make(chan Result, 1000) // Buffered pipeline collecting found assets
	var wg sync.WaitGroup

	fmt.Printf("[*] Launching Worker Pool with %d active threads...\n", workerCount)
	fmt.Printf("[*] Targeted Domain Space: %s\n\n", targetDomain)
	startTime := time.Now()

	// 4. Spin up our fixed pool of worker threads
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go worker(targetDomain, jobs, results, &wg)
	}

	// 5. Spawn an asymmetric background supervisor to close the results stream
	// once all active worker tasks report they are completely finished.
	go func() {
		wg.Wait()
		close(results)
	}()

	// 6. Spawn a dedicated execution thread to read the wordlist and populate the jobs channel.
	// This prevents the main loop from blocking if the channels fill up.
	go func() {
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			word := strings.TrimSpace(scanner.Text())
			if word != "" {
				jobs <- word // Injects work into the pipeline
			}
		}
		close(jobs) // Closing jobs tells the workers there are no more tasks coming
	}()

	// 7. The main thread acts as the centralized result processing engine.
	// It reads from the results channel as findings pop up and prints them in real time.
	discoveredCount := 0
	for res := range results {
		discoveredCount++
		fmt.Printf("[+] Discovered: %-30s -> IPs: %v\n", res.Subdomain, res.IPs)
	}

	elapsed := time.Since(startTime)
	fmt.Println("\n=================== INFRASTRUCTURE ENUMERATION COMPLETE ===================")
	fmt.Printf("Scan Execution Duration : %v\n", elapsed)
	fmt.Printf("Total Verified Subdomains: %d\n", discoveredCount)
	fmt.Println("===========================================================================")
}
