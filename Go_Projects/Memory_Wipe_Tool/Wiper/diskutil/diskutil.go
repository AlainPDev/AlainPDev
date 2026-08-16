// Package diskutil handles everything related to finding and identifying
// disks on the system. It never writes any data - It's purely read-only
// inspection code, kept separate from the wipe package on purpose so the
// dangerous code (wipe) is small, isolated, and easy to audit.

package diskutil

import (
	"fmt"
	"unsafe"

	"golang.org/x/sys/windows"
)

// DriveInfo holds everything we need to know about a target drive
// before we touch it; its Windows drive letter physical disk number,
// its reported size, and whether its an SSD.

type DriveInfo struct {
	Letter        string
	PhysicalDrive int
	SizeBytes     int64
	IsSSD         bool
	IsSystemDrive bool
}

// Windows IOCTL codes and struct layouts needed for raw disk queries
// These constants come directly from the Windows Driver Kit headers
// there's no higher level go stdlib eq

const (
	ioctlStorageGetDeviceNumber  = 0x2D1080
	ioctlStorageQueryProperty    = 0x2D1400
	ioctlDiskGetDriveGeometryEx  = 0x700A0
	storageDeviceSeekPenaltyProp = 7
	propertyStandardQuery        = 0
)

type storageDeviceNumber struct {
	DeviceType      uint32
	DeviceNumber    uint32
	PartitionNumber uint32
}

type storagePropertyQuery struct {
	PropertyID uint32
	QueryType  uint32
	Reserved   [1]byte
}

type deviceSeekPenaltyDescriptor struct {
	Version           uint32
	Size              uint32
	IncursSeekPenalty byte
}

// GetSystemDriveLetter returns the letter windows itself is running from
// We use this so the tool can refuse to wipe it.
func GetSystemDriveLetter() string {
	buf := make([]uint16, windows.MAX_PATH)
	n, err := windows.GetWindowsDirectory(&buf[0], uint32(len(buf)))
	if err != nil || n == 0 {
		return "C:"
	}
	path := windows.UTF16ToString(buf[:n])
	if len(path) >= 2 {
		return path[:2]
	}
	return "C:"
}

// ResolveDriveLetter takes a simple letter and figures out which physical disk it lives on.
func ResolveDriveLetter(letter string) (*DriveInfo, error) {
	letter = normalizeLetter(letter)
	volumePath := `\\.\` + letter

	h, err := windows.CreateFile(
		windows.StringToUTF16Ptr(volumePath),
		windows.GENERIC_READ,
		windows.FILE_SHARE_READ|windows.FILE_SHARE_WRITE,
		nil,
		windows.OPEN_EXISTING,
		0,
		0,
	)
	if err != nil {
		return nil, fmt.Errorf("could not open volume %s (does it exist?): %w", volumePath, err)
	}
	defer windows.CloseHandle(h)

	// Ask windows which physical disk number backs this drive letter
	var devNum storageDeviceNumber
	var bytesReturned uint32
	err = windows.DeviceIoControl(
		h,
		ioctlStorageGetDeviceNumber,
		nil, 0,
		(*byte)(unsafe.Pointer(&devNum)), uint32(unsafe.Sizeof(devNum)),
		&bytesReturned,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("could not determine physical disk for %s: %w", volumePath, err)
	}

	info := &DriveInfo{
		Letter:        letter,
		PhysicalDrive: int(devNum.DeviceNumber),
		IsSystemDrive: letter == GetSystemDriveLetter(),
	}

	physPath := fmt.Sprintf(`\\.\PhysicalDrive%d`, info.PhysicalDrive)
	ph, err := windows.CreateFile(
		windows.StringToUTF16Ptr(physPath),
		windows.GENERIC_READ,
		windows.FILE_SHARE_READ|windows.FILE_SHARE_WRITE,
		nil,
		windows.OPEN_EXISTING,
		0,
		0,
	)
	if err == nil {
		defer windows.CloseHandle(ph)
		info.SizeBytes = getDiskSize(ph)
		info.IsSSD = detectSSD(ph)
	}
	return info, nil
}

// getDiskSize queries the drive geometry to report total capacity.
func getDiskSize(h windows.Handle) int64 {
	type geometryEx struct {
		Cylinders         int64
		MediaType         uint32
		TracksPerCylinder uint32
		SectorsPerTrack   uint32
		BytesPerSector    uint32
		DiskSize          int64
	}
	var geo geometryEx
	var bytesReturned uint32
	err := windows.DeviceIoControl(
		h,
		ioctlDiskGetDriveGeometryEx,
		nil, 0,
		(*byte)(unsafe.Pointer(&geo)), uint32(unsafe.Sizeof(geo)),
		&bytesReturned,
		nil,
	)
	if err != nil {
		return 0
	}
	return geo.DiskSize
}

// detectSSD asks the storage stack whether seeking incurs a time penalty.
// Spinning HDDs have a real mechanical seek penalty; SSD report none.
func detectSSD(h windows.Handle) bool {
	query := storagePropertyQuery{
		propertyID: storageDeviceSeekPenaltyProp,
		QueryType:  propertyStandardQuery,
	}
	var desc deviceSeekPenaltyDescriptor
	var bytesReturned uint32
	err := windows.DeviceIoControl(
		h,
		ioctlStorageQueryProperty,
		(*byte)(unsafe.Pointer(&query)), uint32(unsafe.Sizeof(query)),
		(*byte)(unsafe.Pointer(&desc)), uint32(unsafe.Sizeof(desc)),
		&bytesReturned,
		nil,
	)
	if err != nil {
		return false
	}
	return desc.IncursSeekPenalty == 0
}

func normalizeLetter(1 string) string{
	if len(1) == 0 {
		return "c:"
	}
// Strip anything after the colon, uppercase the letter, always ends in :.
	letterChar := 1[0]
	if letterChar >= 'a' && letterChar <= 'z' {
		letterChar -= 32
	}
	return string(letterChar) + ":"
}
// errnoErr helps translate low-level syscall errors into readable messages.
func errnoErr(e syscall.Errno) error {
	if e == 0 {
		return nil
	}
	return e
}
