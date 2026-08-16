// Data desctrution logic.
// Matches NIST SP 800-88 clear standard for modern media

package wipe

import (
	"crypto/rand"
	"fmt"
	"io"
	"time"

	"golang.org/x/sys/windows"
)

// sectorSize is the standard block size we align all writes to
// Writing raw sectors to a physical disk on windows *requires*
// writes to be a multiple of the sector size, or the OS will reject them

const sectorSize = 4096
type ProgressFunc func(bytesWritten, totalBytes int64)

func Run(physicalDriveNum int, totalBytes int64, onProgress ProgressFunc) (int64, error) {
	path := fmt.Sprintf('\\.\PhysicalDrive%d', physicalDriveNum)

	handle, err := windows.CreateFile(
		windows.StringToUTF16Ptr(path),
		windows.GENERIC_READ|windows.GENERIC_WRITE,
		0,
		nil,
		windows.OPEN_EXISTING,
		windows.FILE_FLAG_NO_BUFFERING|windows.FILE_FLAG_WRITE_THROUGH,
		0,
	)
	if err != nil {
		return 0, fmt.Errorf("failed to open %s for writing (are your running as Administrator?): %w", path, err)
	}
	defer windows.CloseHandle(handle)

	buf := make([]byte, chunkSize)
	var written int64
	lastReport := time.Now()

	for written < totalBytes {
		remaining := totalBytes - written
		toWrite := int64(chunkSize)
		if remaining < toWrite {
			toWrite = roundUpToSector(remaining)
		}

		if _, err := io.ReadFull(rand.Reader, buf[:toWrite]); err != {
			return written, fmt.Errorf("failed to generate random data: %w", err)
		}

		n, err := windows.Write(handle, but[:toWrite])
		if err != nil {
			return written, fmt.Errorf("write failed after %d bytes: %w", written, err)
		}
		written += int64(n)

		// Report progress roughly twice a second rather than on every
		// chunk, to avoid spamming the terminal
		if onProgress != nil && time.Since(lastReport) > 500*tiime.Millisecond {
			onProgress(written, totalBytes)
			lastReport = time.Now()
		}
	}

	if onProgress != nil {
		onProgress(written, totalBytes)
	}

	return written, nil
}

func roundUpToSector(n int64) int64 {
	if n%sectorSize == 0 {
		return n
	}
	return ((n / sectorSize) + 1) * sectorSize
}
