package main

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"
)

func main() {
	printHeader("1. Checking Server Status")
	if !isServerRunning("localhost:8080") {
		fmt.Println("❌ Error: Server is NOT running on localhost:8080")
		fmt.Println("Try running 'go run cmd/server/main.go'...")
		cmd := exec.Command("go", "run", "cmd/server/main.go")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		go func() {
			if err := cmd.Run(); err != nil {
				fmt.Printf("❌ Failed to run server: %v\n", err)
				os.Exit(1)
			}
		}()
	}
	fmt.Println("✅ Server is reachable.")

	printHeader("2. Opening Browser Client")
	htmlPath := filepath.Join("cmd", "test", "index.html")
	absPath, err := filepath.Abs(htmlPath)
	if err != nil {
		fmt.Printf("❌ Could not resolve path: %v\n", err)
		os.Exit(1)
	}

	if err := openBrowser(absPath); err != nil {
		fmt.Printf("⚠️  Could not open browser automatically: %v\n", err)
		fmt.Printf("   Please manually open: %s\n", absPath)
	} else {
		fmt.Println("✅ Browser opened.")
	}

	// Give the browser a moment to connect via WebSocket
	fmt.Println("   Waiting 2 seconds for WebSocket connection...")
	time.Sleep(2 * time.Second)

	printHeader("3. Running Publisher")
	cmd := exec.Command("go", "run", "cmd/publisher/main.go")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Printf("❌ Failed to run publisher: %v\n", err)
		os.Exit(1)
	}

	printHeader("TEST COMPLETE")
	fmt.Println("Check your browser window. You should see the notification!")
}

// Helper to check if a TCP port is open
func isServerRunning(address string) bool {
	timeout := 1 * time.Second
	conn, err := net.DialTimeout("tcp", address, timeout)
	if err != nil {
		return false
	}
	if conn != nil {
		conn.Close()
		return true
	}
	return false
}

// Helper to open file in default browser based on OS
func openBrowser(path string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", "", path)
	case "darwin":
		cmd = exec.Command("open", path)
	default:
		cmd = exec.Command("xdg-open", path)
	}

	return cmd.Start()
}

func printHeader(msg string) {
	fmt.Printf("\n--- %s ---\n", msg)
}
