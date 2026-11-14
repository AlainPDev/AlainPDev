package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
)

func main() {
	// 1. Argument Validation
	if len(os.Args) < 1 {
		fmt.Println("Usage: go run main.go <path_to_auth_or_syslog>")
		os.Exit(1)
	}
	logFilePath := os.Args[1]

	// 2. Open the log file safely
	file, err := os.Open(logFilePath)
	if err != nil {
		fmt.Printf("[!] Failed to open log file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close() // Ensure the file handles are released when main exits

	// 3. Define the Regex pattern to catch failed passwords
	failedPattern := regexp.MustCompile(`Failed password .* from ([0-9.]+)`)

	// 4. Instantiate a Map to hold our metrics
	failedAttempts := make(map[string]int)

	fmt.Printf("[*] Analyzing log file: %s ...\n", logFilePath)

	// 5. Use a buffered scanner to read the file line by line efficiently
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// Check if the current line matches our failed password pattern
		match := failedPattern.FindStringSubmatch(line)
		if len(match) > 1 {
			// match[0] is the entire matching string
			// match[1] is the data inside our first capture group
			attackerIP := match[1]
			failedAttempts[attackerIP]++
		}
	}

	// 6. Check if the scanner encountered any physical read errors
	if err := scanner.Err(); err != nil {
		fmt.Printf("[!] Error reading file contents: %v\n", err)
		os.Exit(1)
	}
	// 7. Print the final security findings summary
	fmt.Println("\n=============== Security Audit Report =============")
	fmt.Printf("Analysis complete. Found %d unique malicious IP addresses.\n\n", len(failedAttempts))

	if len(failedAttempts) == 0 {
		fmt.Println("[+] No failed login attempts detected. System baseline secure.")
	} else {
		fmt.Printf("%-20s | %s\n", "Attacker IP Address", "Failed Login Attempts")
		fmt.Println("------------------------------------------------------------")

		// Iterate through our map to display the keys and values
		for ip, count := range failedAttempts {
			// Highlight IP that have a high volume of attempts
			alertTag := ""
			if count > 5 {
				alertTag = " [!] High Brute Force Risk"
			}
			fmt.Printf("%-20s | %-10d%s\n", ip, count, alertTag)
		}
	}
	fmt.Println("====================================================================")
}
