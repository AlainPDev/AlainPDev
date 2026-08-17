// Package wipelog writes a simple completion certificate after a wipe
// finishes. Sellers can keep this as proof for a buyer, and it's a nice
// product feature in its own right (resale marketplaces and IT asset
// disposition shops often want documented proof of erasure).
package wipelog

import (
	"fmt"
	"os"
	"time"
)

// Entry captures everything worth recording about one wipe operation.
type Entry struct {
	DriveLetter   string
	PhysicalDrive int
	SizeBytes     int64
	WasSSD        bool
	StartedAt     time.Time
	FinishedAt    time.Time
	Success       bool
	ErrorMessage  string
}

// Write appends a human-readable record to wipe-log.txt in the same
// directory as the executable, and returns the path it wrote to.
func Write(e Entry) (string, error) {
	logPath := "wipe-log.txt"

	f, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return "", fmt.Errorf("could not open log file: %w", err)
	}
	defer f.Close()

	status := "SUCCESS"
	if !e.Success {
		status = "FAILED: " + e.ErrorMessage
	}

	mediaType := "HDD"
	if e.WasSSD {
		mediaType = "SSD"
	}

	record := fmt.Sprintf(
		"----------------------------------------\n"+
			"Drive Wipe Certificate\n"+
			"----------------------------------------\n"+
			"Drive letter:     %s\n"+
			"Physical disk #:  %d\n"+
			"Media type:       %s\n"+
			"Size:             %.2f GB\n"+
			"Method:           Single-pass random overwrite (NIST SP 800-88 Clear)\n"+
			"Started:          %s\n"+
			"Finished:         %s\n"+
			"Duration:         %s\n"+
			"Result:           %s\n\n",
		e.DriveLetter,
		e.PhysicalDrive,
		mediaType,
		float64(e.SizeBytes)/(1024*1024*1024),
		e.StartedAt.Format(time.RFC1123),
		e.FinishedAt.Format(time.RFC1123),
		e.FinishedAt.Sub(e.StartedAt).Round(time.Second),
		status,
	)

	if _, err := f.WriteString(record); err != nil {
		return "", fmt.Errorf("could not write log entry: %w", err)
	}

	return logPath, nil
}
