package main

import (
        "crypto/sha256"
        "encoding/hex"
        "encoding/json"
        "fmt"
        "io"
        "os"
        "path/filepath"
)

// FileInfo holds the structural details we want to save for each file
type FileInfo struct {
        Path string `json:"path"`
        Hash string `json:"hash"`
}

// hashFile safely reads a file in chunks and computes its SHA-256 hash
func hashFile(filePath string) (string, error) {
        file, err := os.Open(filePath)
        if err != nil {
                return "", err
        }
        defer file.Close() // Guarantees the OS releases this file handle

        hasher := sha256.New()

        // Copy data in small chuncks directly to the hasher
        if _, err := io.Copy(hasher, file); err != nil {
                return "", err

        }

        return hex.EncodeToString(hasher.Sum(nil)), nil
}

func main() {
        // Basic Argument Validation
        if len(os.Args) < 3 {
                fmt.Println("Usage: go run main.go <directory_to_scan> <output_baseline.json>")
                os.Exit(1)
        }

        targetDir := os.Args[1]
        outputFile := os.Args[2]

        // Instantiation of a slice to store constructs
        var baselineData []FileInfo

        fmt.Printf("[*} Scanning directory: %s\n", targetDir)

        // Walking the Directory Tree
        err := filepath.Walk(targetDir, func(path string, info os.FileInfo, err error) error {
                if err != nil {
                        fmt.Printf("[-] Skipping inaccesible path %s: %v\n", path, err)
                        return nil // Tells walk to ignore errors
                }

                // We only want files not directories
                if info.IsDir() {
                        return nil
                }

                hash, err := hashFile(path)
                if err != nil {
                        fmt.Printf("[-] Failed to hash %s: %v\n", path, err)
                        return nil
                }

                // Create an instance of our struct and append it to slice
                entry := FileInfo{Path: path, Hash: hash}
                baselineData = append(baselineData, entry)

                return nil

        })

        if err != nil {
                fmt.Printf("[!] Fatal error traversing directory: %v\n", err)
                os.Exit(1)

        }

        // Marshalling data into readable JSON formatting
        jsonData, err := json.MarshalIndent(baselineData, "", " ")
        if err != nil {
                fmt.Printf("[!] Failed to encode JSON: %v\n", err)
                os.Exit(1)

        }

        // Writing the data output to a file on disk
        err = os.WriteFile(outputFile, jsonData, 0644)
        if err != nil {
                fmt.Printf("[!] Failed to save baseline file: %v\n", err)
                os.Exit(1)

        }

        fmt.Printf("[+] Success! Baseline saved with %d tracked files to: %s\n", len(baselineData), outputFile)
}
