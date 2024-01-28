package main

import (
	"dnsServer/server"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// Start the DNS server and get the stop channel
	dnsServer, err := server.NewDNSServer(":8080")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	dnsServer.Start()

	// Set up channel to listen for OS interrupt signal (e.g., CTRL+C)
	osSignal := make(chan os.Signal, 1)
	signal.Notify(osSignal, syscall.SIGINT, syscall.SIGTERM)

	// Block until an OS signal is received
	sig := <-osSignal
	fmt.Printf("Received signal: %v, shutting down...\n", sig)

	// Send stop signal to the DNS server
	dnsServer.Stop()

	// Wait a moment for the server to shut down gracefully
	time.Sleep(time.Second)

	fmt.Println("Server stopped")
}
