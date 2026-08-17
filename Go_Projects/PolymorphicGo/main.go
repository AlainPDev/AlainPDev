package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <output_file_path>")
		os.Exit(1)
	}

	outputPath := os.Args[1]

	// 1. Define the core execution logic/data of our application
	coreLogic := "FUNCTION: ExecuteSystemAudit(); STATUS: Active;"

	// 2. Generate a unique runtime marker (a timestamp)
	// This acts as a structural modifier that changes on every single execution
	uniqueMarker := time.Now().Format(time.RFC3339Nano)

	// 3. Combine the core logic with the unique marker
	// The structural behavior remains identical, but the raw data array changes
	combinedPayload := fmt.Sprintf("%s\n// Build Marker: %s\n", coreLogic, uniqueMarker)

	// 4. Calculate the unique cryptographic signature of this specific iteration
	hasher := sha256.New()
	hasher.Write([]byte(combinedPayload))
	computedHash := hex.EncodeToString(hasher.Sum(nil))

	// 5. Commit the payload to disk
	err := os.WriteFile(outputPath, []byte(combinedPayload), 0644)
	if err != nil {
		fmt.Printf("[!] Error saving data structure: %v\n", err)
		os.Exit(1)
	}

	// 6. Display the metrics to the analyst
	fmt.Println("================= COMPILATION STRUCT WATCH =================")
	fmt.Printf("Target File Out   : %s\n", outputPath)
	fmt.Printf("Generated Marker  : %s\n", uniqueMarker)
	fmt.Printf("Resulting Hash    : %s\n", computedHash)
	fmt.Println("[*] Status: File signature successfully modified.")
	fmt.Println("============================================================")
}
