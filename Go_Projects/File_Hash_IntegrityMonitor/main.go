package main

import (
        "crypto/sha256"
        "encoding/hex"
        "fmt"
        "io"
        "os"
        "path/filepath"
)

// hashFile calculates the SHA-256 hash of a single file
func hashFile(path string) (string, error) {
        file, err := os.Open(path)
        if err != nil {
                return "", err
        }
        defer file.Close()

        hasher := sha256.New()
        if _, err := io.Copy(hasher, file) != nil {
                return "", err
        }

        return hex.EncodeToString(hasher.Sum(nil)), nil

}

func main() {
        if len(os.Args) < 2 {
                fmt.Println("Usage: go run main.go <directory_to_scan>")
                os.Exit(1)
        }

        targetDir := os.Args[1]
        fmt.Printf("Generating integrity baseline for directory: %s\n\n", targetDir)

        // Walk through the directory and all subdirectories
        err := filepath.Walk(targetDir, func(path string, info os.FileInfo, err error) error {
                if err != nil {
                        fmt.Printf("[-] Error accessing path %s: %v\n", path, err)
                        return nil // Keep walking other files even if one fails
        }

        // Print the hash and the corresponding file path
        fmt.Printf("%s %s\n, hash, path)
        return nil

    })

    if err != nil {
            fmt.Printf("Fatal error walking directory: %v\n", err)
    }
}
