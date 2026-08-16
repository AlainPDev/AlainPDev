module laptop-wipe

go 1.21

// golang.org/x/sys gives access to raw Windows syscalls
// CreateFile with disk-level flags, etc. that are not in standard library
require golang.org/x/sys v0.20.0
