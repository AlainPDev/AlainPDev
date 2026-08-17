package main

import (
	"fmt"
	"net"
	"os"
	"sort"
	"strconv"
	"sync"
	"time"
)

// scanPort attempts a TCP handshake. If successful, it sends the port to the openPorts channel.
func scanPort(ip string, port int, timeout time.Duration, wg *sync.WaitGroup, openPorts chan int) {
	// 1. Notify the WaitGroup that this specific goroutine has finished when it exits
	defer wg.Done()

	// Combine IP and port into standard dial string (e.g., "192.168.1.1:80")
	address := net.JoinHostPort(ip, strconv.Itoa(port))

	// Establish a raw network handshake with an explicit boundary timeout
	conn, err := net.DialTimeout("tcp", address, timeout)
	if err != nil {
		// Port is either closed, filtered by a firewall, or dropped
		return
	}

	// A connection succeeded! Close the socket cleanly to avoid resource leaking
	conn.Close()

	// Send the open port number into our thread-safe reporting channel
	openPorts <- port
}

func main() {
	// 2. Validate input parameters
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <Target-IP>")
		fmt.Println("Example: go run main.go 192.168.1.1")
		os.Exit(1)
	}

	targetIP := os.Args
	timeout := 2 * time.Second

	// 3. Initialize our synchronization and communication structures
	var wg sync.WaitGroup
	openPortsChan := make(chan int, 1024) // Buffered channel to safely collect results concurrently

	fmt.Printf("[*] Launching high-velocity parallel scan against: %s\n", targetIP)
	startTime := time.Now()

	// 4. Spin up the parallel scanning engine across ports 1 to 1024
	for port := 1; port <= 1024; port++ {
		wg.Add(1) // Increment the WaitGroup master thread counter

		// The "go" keyword spawns a lightweight background thread (Goroutine)
		go scanPort(targetIP, port, timeout, &wg, openPortsChan)
	}

	// 5. Manage the background thread lifecycle
	// We wait for all scanPort threads to complete in the background, then safely close the channel
	go func() {
		wg.Wait()
		close(openPortsChan) // Closing tells the upcoming range loop that no more data is arriving
	}()

	// 6. Collect results from the communication channel
	var discoveredPorts []int
	for port := range openPortsChan {
		fmt.Printf("[+] Discovered OPEN Port: %d\n", port)
		discoveredPorts = append(discoveredPorts, port)
	}

	// 7. Sort and display the final security posture summary
	sort.Ints(discoveredPorts)
	elapsedTime := time.Since(startTime)

	fmt.Println("\n=================== NETWORK AUDIT SUMMARY ===================")
	fmt.Printf("Scan Completed In : %v\n", elapsedTime)
	fmt.Printf("Total Open Ports   : %d\n", len(discoveredPorts))

	if len(discoveredPorts) == 0 {
		fmt.Println("[-] No core infrastructure entrypoints discovered.")
	} else {
		fmt.Printf("Identified Vector List: %v\n", discoveredPorts)
	}
	fmt.Println("=============================================================")
}
