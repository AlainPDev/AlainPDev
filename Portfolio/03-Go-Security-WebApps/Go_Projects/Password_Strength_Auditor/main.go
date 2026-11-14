package main

import (
	"fmt"
	"os"
	"unicode"
)

// 1. Define a struct to represent our custome Password Auditor configurations
type PasswordAuditor struct {
	MinLength      int
	RequireUpper   bool
	RequireLower   bool
	RequireNumber  bool
	RequireSpecial bool
}

// 2. Define a struct to hold the detailed findings of a password audit
type AuditResult struct {
	IsValid     bool
	Score       int      // Score from 0 (terrible) to 5 (excellent)
	FailReasons []string // List of policies the password violated
}

// 3. Attach a Method to the PasswordAuditor struct to analyze a password string
func (pa PasswordAuditor) Audit(password string) AuditResult {
	var result AuditResult

	var hasUpper, hasLower, hasNumber, hasSpecial bool
	var length = len(password)

	for _, ch := range password {
		switch {
		case unicode.IsUpper(ch):
			hasUpper = true
		case unicode.IsLower(ch):
			hasLower = true
		case unicode.IsDigit(ch):
			hasNumber = true
		case unicode.IsPunct(ch) || unicode.IsSymbol(ch):
			hasSpecial = true
		}
	}

	// Validate against policies configured in the struct
	if length < pa.MinLength {
		result.FailReasons = append(result.FailReasons, fmt.Sprintf("Password is too short (Minimum %d chars required)", pa.MinLength))
	}
	if pa.RequireUpper && !hasUpper {
		result.FailReasons = append(result.FailReasons, "Missing at least one uppercase letter (A-Z)")
	}
	if pa.RequireLower && !hasLower {
		result.FailReasons = append(result.FailReasons, "Missing at least one lowercase letter (a-z)")
	}
	if pa.RequireNumber && !hasNumber {
		result.FailReasons = append(result.FailReasons, "Missing at least one number (0-9)")
	}
	if pa.RequireSpecial && !hasSpecial {
		result.FailReasons = append(result.FailReasons, "Missing at least one special character (!, @, #, $, etc.")
	}

	// Calculate Strength score
	score := 0
	if length >= pa.MinLength {
		score++
	}
	if hasUpper {
		score++
	}
	if hasLower {
		score++
	}
	if hasNumber {
		score++
	}
	if hasSpecial {
		score++
	}
	result.Score = score

	if len(result.FailReasons) == 0 {
		result.IsValid = true
	}

	return result
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <password_to_test>")
		os.Exit(1)
	}

	testPassword := os.Args[1]

	policy := PasswordAuditor{
		MinLength:      12,
		RequireUpper:   true,
		RequireLower:   true,
		RequireNumber:  true,
		RequireSpecial: true,
	}

	report := policy.Audit(testPassword)

	// 6. Display formatted audit metrics
	fmt.Println("================Password Audit Report ============")
	fmt.Printf("Password Evaluated : %s\n", testPassword)
	fmt.Printf("Strength Score		: %d/5\n", report.Score)

	if report.IsValid {
		fmt.Println("[+] Compliance Status: COMPLIANT (Strong Password)")
	} else {
		fmt.Println("[-] Compliance Status: NON-COMPLIANT (Weak Password)")
		fmt.Println("\nViolations Detected:")
		for _, reason := range report.FailReasons {
			fmt.Printf("  • %s\n", reason)
		}
	}
	fmt.Println("================================================================")
}
