// Command laptop-wipe is a lightweight, fast alternative to DBAN for
// people reselling an old Windows laptop. It performs a single-pass
// random overwrite of a chosen drive (NIST SP 800-88 "Clear" standard),
// which is dramatically faster than DBAN's legacy multi-pass schemes
// while still rendering the previous data practically unrecoverable
// for a home/resale threat model.
//
// IMPORTANT: This tool must be run as Administrator, and is only
// safe/meaningful to run from outside the drive being wiped (e.g. a
// bootable Windows PE / USB recovery environment), since it will refuse
// to touch the currently-running system drive as a safety measure.
package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"laptop-wipe/diskutil"
	"laptop-wipe/wipe"
	"laptop-wipe/wipelog"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("==========================================")
	fmt.Println(" Laptop Wipe — Fast Secure Erase Tool")
	fmt.Println("==========================================")
	fmt.Println("This will PERMANENTLY destroy all data on the drive you choose.")
	fmt.Println("It cannot be undone. Make sure you have backed up anything you")
	fmt.Println("want to keep before continuing.")
	fmt.Println()

	// --- Step 1: ask which drive letter to wipe -----------------------
	fmt.Print("Enter the drive letter to wipe (e.g. D): ")
	letterInput, _ := reader.ReadString('\n')
	letterInput = strings.TrimSpace(letterInput)

	info, err := diskutil.ResolveDriveLetter(letterInput)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	// --- Step 2: hard safety block on the system drive -----------------
	// This is the single most important guardrail in the whole program.
	// No amount of user confirmation should be able to override this —
	// wiping the drive Windows is currently running from would crash
	// mid-wipe and isn't what any legitimate user actually wants anyway.
	if info.IsSystemDrive {
		fmt.Println()
		fmt.Printf("Refusing to wipe %s — this is your active Windows system drive.\n", info.Letter)
		fmt.Println("Boot from a USB recovery/installer drive instead, then run this")
		fmt.Println("tool from there to wipe the internal drive.")
		os.Exit(1)
	}

	mediaType := "HDD"
	if info.IsSSD {
		mediaType = "SSD"
	}
	sizeGB := float64(info.SizeBytes) / (1024 * 1024 * 1024)

	// --- Step 3: simple, clear confirmation ----------------------------
	// Kept intentionally simple per product requirements: normal resale
	// users know their drive letter, not model numbers or serials, so we
	// confirm using the same plain terms they already understand.
	fmt.Println()
	fmt.Printf("Target drive:  %s\n", info.Letter)
	fmt.Printf("Type:          %s\n", mediaType)
	fmt.Printf("Size:          %.2f GB\n", sizeGB)
	fmt.Println()
	fmt.Printf("Type WIPE to permanently erase drive %s, or anything else to cancel: ", info.Letter)
	confirm, _ := reader.ReadString('\n')
	confirm = strings.TrimSpace(confirm)

	if confirm != "WIPE" {
		fmt.Println("Cancelled. No data was touched.")
		os.Exit(0)
	}

	// --- Step 4: run the wipe with a live progress bar ------------------
	fmt.Println()
	fmt.Println("Wiping... do not close this window or turn off the computer.")
	startedAt := time.Now()

	written, err := wipe.Run(info.PhysicalDrive, info.SizeBytes, func(done, total int64) {
		pct := float64(done) / float64(total) * 100
		fmt.Printf("\r  Progress: %6.2f%%  (%.2f GB / %.2f GB)   ",
			pct,
			float64(done)/(1024*1024*1024),
			float64(total)/(1024*1024*1024),
		)
	})
	finishedAt := time.Now()
	fmt.Println() // move past the progress line

	// --- Step 5: write the completion log/certificate -------------------
	logEntry := wipelog.Entry{
		DriveLetter:   info.Letter,
		PhysicalDrive: info.PhysicalDrive,
		SizeBytes:     written,
		WasSSD:        info.IsSSD,
		StartedAt:     startedAt,
		FinishedAt:    finishedAt,
		Success:       err == nil,
	}
	if err != nil {
		logEntry.ErrorMessage = err.Error()
	}

	logPath, logErr := wipelog.Write(logEntry)

	if err != nil {
		fmt.Println("Wipe failed:", err)
		if logErr == nil {
			fmt.Println("Details recorded in", logPath)
		}
		os.Exit(1)
	}

	fmt.Println()
	fmt.Printf("Done. Drive %s has been wiped (%.2f GB in %s).\n",
		info.Letter, sizeGB, finishedAt.Sub(startedAt).Round(time.Second))
	if logErr == nil {
		fmt.Println("A wipe certificate was saved to", logPath)
	}
}
