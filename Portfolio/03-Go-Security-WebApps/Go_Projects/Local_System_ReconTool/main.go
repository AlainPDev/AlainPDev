package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// ReconData structure maps out the system analytics we intend to pull
type ReconData struct {
	OS           string `json:"operating_system"`
	Architecture string `json:"architecture"`
	Currentuser  string `json:"current_user"`
	HomeDir      string `json:"home_directory"`
	KernelInfo   string `json:"kernel_info_raw"`
}

// executeCommand runs a low-level terminal binary and returns its trimmed output string
func executeCommand(binary string, arguments ...string) (string, error) {
	// Initialize the command executor execution block
	cmd := exec.Command(binary, arguments...)

	// Create memory buffers to capture what the program prints to stdout and stderr
	var outBuffer bytes.Buffer
	var errBuffer bytes.Buffer
	cmd.Stdout = &outBuffer
	cmd.Stderr = &errBuffer

	// Run the system command
	err := cmd.Run()
	if err != nil {
		// Return the error message printed by the OS terminal if it fails
		return "", fmt.Errorf("%v: %s", err, strings.TrimSpace(errBuffer.String()))
	}

	return strings.TrimSpace(outBuffer.String()), nil
}

func main() {
	fmt.Println("[*] Initiating local system intelligence audit...")

	// 1. Gather environmental details native to the Go runtime engine
	data := ReconData{
		OS:           runtime.GOOS,
		Architecture: runtime.GOARCH,
	}

	// 2. Extract values using the environment variable lookup package
	if runtime.GOOS == "windows" {
		data.Currentuser = os.Getenv("USERNAME")
		data.HomeDir = os.Getenv("USERPROFILE")
	} else {
		data.Currentuser = os.Getenv("USER")
		data.HomeDir = os.Getenv("HOME")
	}

	// 3. Execute platform-specific kernel investigations
	var kernelOutput string
	var err error

	if runtime.GOOS == "windows" {
		// Windows system information query
		kernelOutput, err = executeCommand("cmd", "/c", "ver")
	} else {
		// Linux/macOS system information query (runs: uname -a)
		kernelOutput, err = executeCommand("uname", "-a")
	}

	if err != nil {
		data.KernelInfo = fmt.Sprintf("Error retrieving info: %v", err)
	} else {
		data.KernelInfo = kernelOutput
	}

	// 4. Transform our structured dataset into an enterprise JSON output format
	jsonReport, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Printf("[!] JSON Serialization failure: %v\n", err)
		os.Exit(1)
	}

	// 5. Output findings directly to the console standard output
	fmt.Println("\n=================== RECON SUMMARY REPORT ===================")
	fmt.Println(string(jsonReport))
	fmt.Println("==============================================================")
}
