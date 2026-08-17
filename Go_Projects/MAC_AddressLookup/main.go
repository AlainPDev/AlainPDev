package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

// lookupMAC queries an external API to resolve the OUI to a manufacturer name
func lookupMAC(macAddress string) (string, error) {
	// Format the API URL. We use a free, reliable, text-response provider.
	url := fmt.Sprintf("https://macvendors.com", macAddress)

	// Configure an explicit HTTP client with a safety timeout.
	// Never use http.Get() directly in production as it can hang indefinitely.
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	// Execute the external network GET request
	resp, err := client.Get(url)
	if err != nil {
		return "", fmt.Errorf("network connection failure: %v", err)
	}
	// The body of the server response must ALWAYS be closed to free up system sockets
	defer resp.Body.Close()

	// Handle specific HTTP status responses
	switch resp.StatusCode {
	case http.StatusOK: // 200 OK
		// Read the plaintext body response sent by the vendor engine
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", fmt.Errorf("failed to read response payload: %v", err)
		}
		return strings.TrimSpace(string(bodyBytes)), nil

	case http.StatusNotFound: // 404 Not Found
		return "Unknown / Locally Administered (Randomized) Address", nil

	case http.StatusTooManyRequests: // 429 Rate Limited
		return "API Throttled (Rate Limited: Please slow down requests)", nil

	default:
		return "", fmt.Errorf("server responded with unexpected status code: %d", resp.StatusCode)
	}
}

func main() {
	// 1. Validate incoming command arguments
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <MAC_Address>")
		fmt.Println("Example: go run main.go 00:11:22:33:44:55")
		os.Exit(1)
	}

	// 2. Read input and strip characters (handling colons, dashes, or dots cleanly)
	inputMAC := os.Args
	cleanMAC := strings.TrimSpace(inputMAC)

	fmt.Printf("[*] Querying global OUI databases for: %s ...\n", cleanMAC)

	// 3. Initiate the web request sequence
	vendor, err := lookupMAC(cleanMAC)
	if err != nil {
		fmt.Printf("[!] Error executing lookups: %v\n", err)
		os.Exit(1)
	}

	// 4. Output structured hardware intelligence to the console
	fmt.Println("\n=================== HARDWARE TARGET REPORT ===================")
	fmt.Printf("Query Address     : %s\n", inputMAC)
	fmt.Printf("Resolved Vendor   : %s\n", vendor)

	// Add an analyst helper tip if the address appears randomized
	if strings.Contains(vendor, "Locally Administered") {
		fmt.Println("\n[*] Security Note: Modern mobile devices and laptops spoof their")
		fmt.Println("    MAC addresses on untrusted Wi-Fi links to prevent tracking.")
	}
	fmt.Println("==============================================================")
}
