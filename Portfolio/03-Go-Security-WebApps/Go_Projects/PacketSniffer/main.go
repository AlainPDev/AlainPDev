package main

import (
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"syscall"
)

// EthernetHeader represents the parsed structures of a standard IEEE 802.3 Ethernet Frame
type EthernetHeader struct {
	DstMAC    net.HardwareAddr
	SrcMAC    net.HardwareAddr
	EtherType uint16
}

// parseEthernetFrame takes a raw slice of bytes from the wire and extracts header structures
func parseEthernetFrame(buffer []byte) (EthernetHeader, error) {
	var header EthernetHeader

	if len(buffer) < 14 {
		return header, fmt.Errorf("packet too short to contain a valid Ethernet header")
	}

	// 1. Slice out the MAC Addresses (6 bytes each)
	// Destination MAC occupies bytes 0 to 5. Source MAC occupies bytes 6 to 11.
	header.DstMAC = net.HardwareAddr(buffer[0:6])
	header.SrcMAC = net.HardwareAddr(buffer[6:12])

	// 2. Extract the EtherType (Bytes 12 and 13)
	// Network data arrives over the wire in Big-Endian (Network Byte Order).
	// We read those two specific bytes and convert them into a standard uint16.
	header.EtherType = binary.BigEndian.Uint16(buffer[12:14])

	return header, nil
}

func main() {
	fmt.Println("[*] Initializing Raw Layer-2 Network Sniffer...")

	// 3. Establish a low-level Raw Socket connection with the Linux Kernel
	// ETH_P_ALL tells the kernel to clone EVERY single incoming packet on the wire to our binary.
	// syscall.SOCK_RAW bypasses the TCP/IP stack completely to give us raw frames.
	fd, err := syscall.Socket(syscall.AF_PACKET, syscall.SOCK_RAW, int(htons(syscall.ETH_P_ALL)))
	if err != nil {
		fmt.Printf("[!] Failed to create raw socket: %v\n", err)
		fmt.Println("[!] Did you run this script using 'sudo' permissions?")
		os.Exit(1)
	}
	defer syscall.Close(fd)

	fmt.Println("[+] Raw Socket binding successful. Promiscuous parsing active. Press Ctrl+C to terminate...")
	fmt.Println("--------------------------------------------------------------------------------")

	// 4. Infinite capture loop using a pre-allocated byte buffer
	buffer := make([]byte, 65535) // 65,535 bytes is the maximum theoretical length of an IP packet
	for {
		// Read raw bytes straight from the hardware network interface
		numBytes, _, err := syscall.Recvfrom(fd, buffer, 0)
		if err != nil {
			fmt.Printf("[-] Packet dropping anomaly: %v\n", err)
			continue
		}

		// 5. Parse the captured slice data
		ethFrame, err := parseEthernetFrame(buffer[:numBytes])
		if err != nil {
			continue
		}

		// 6. Translate the parsed EtherType code to understand the packet payload
		var proto string
		switch ethFrame.EtherType {
		case 0x0800:
			proto = "IPv4"
		case 0x0806:
			proto = "ARP (Address Resolution Protocol)"
		case 0x86DD:
			proto = "IPv6"
		default:
			proto = fmt.Sprintf("Unknown EtherType (0x%04X)", ethFrame.EtherType)
		}

		// 7. Output live network telemetry to the console standard output
		fmt.Printf("[+] Captured Packet Frame: %d bytes\n", numBytes)
		fmt.Printf("    • Source MAC      : %s\n", ethFrame.SrcMAC.String())
		fmt.Printf("    • Destination MAC : %s\n", ethFrame.DstMAC.String())
		fmt.Printf("    • Network Payload : %s\n\n", proto)
	}
}

// htons converts a 16-bit integer from Host Byte Order to Network Byte Order.
// Big-Endian and Little-Endian architectures interpret numbers backwards.
// We use bitwise operations to ensure the Linux kernel reads our interface request flags correctly.
func htons(i uint16) uint16 {
	return (i<<8)&0xff00 | (i>>8)&0x00ff
}
