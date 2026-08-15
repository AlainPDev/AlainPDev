package main

import (
        "bufio"
        "fmt"
        "net"
        "os"
        "sync"
        "time"
)

func scanPort(ip string, port int, wg *sync.WaitGroup) {
        defer wg.Done()

        address := fmt.Sprintf("%s:%d", ip, port)
        dialer := net.Dialer{Timeout: 2 * time.second}

        conn, err := dialer.Dial("tcp", address)
        if err != nil {
                return
        }
        defer conn.Close()

        fmt.Printf("[+] Port %d is OPEN\n", port)

        conn.SetReadDeadline(time.Now().Add(2 * time.Second))
        reader := bufio.NewReader(conn)
        banner, err := reader.ReadString('\n')
        if err == nil && len(banner) > 0 {
                fmt.Printf("    --> Banner grabbed: %s", banner)
        }
}

func main() {
        if len(os.Args) < 2 {
                fmt.Println("Usage: go run main.go <Target-IP>")
                os.Exit(1)
        }

        targetIP := os.Args[1]
        var wg sync.WaitGroup

        fmt.Printf("Starting scan on target: %s\n\n", targetIP)

        for port := 1; port <= 1024; port++ {
                wg.Add(1)
                go scanPort(targetIP, port, &wg)
        }

        wg.Wait()
        fmt.Println("\nScan process complete.")
}
