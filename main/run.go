package main

import (
	"context"
	"dnsServer/api"
	"dnsServer/server"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	// Start the DNS server and get the stop channel
	Run()
}

func Run() {
	var wg sync.WaitGroup
	wg.Add(1)
	dnsServer, err := server.NewDNSServer(":53")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	dnsServer.Start()

	handle := api.StartApiServer(":8080")
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-stopChan // Wait for SIGINT/SIGTERM
		println("Shutting down server...")

		// Create a context with a timeout
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := handle.Shutdown(ctx); err != nil {
			println("Shutdown error:", err)
		} else {
			println("Server gracefully stopped")
		}
		dnsServer.Stop()

		// Wait a moment for the server to shut down gracefully
		time.Sleep(time.Second)

		fmt.Println("Server stopped")
		wg.Done()
	}()
	// Send stop signal to the DNS server

	wg.Wait() // Wait for server to shut down

}
