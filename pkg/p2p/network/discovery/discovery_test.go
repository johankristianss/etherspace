package network

import (
	"context"
	"fmt"
	"testing"
	"time"

	"golang.org/x/exp/rand"
)

// GetRandomAvailableAddress generates a random local address for testing
func getRandomAvailableAddress() string {
	// Generate a random port in the range 10000-60000 to avoid conflicts
	port := rand.Intn(50000) + 10000
	// Use a fixed local IP (not bound to the OS)
	return fmt.Sprintf("127.0.0.1:%d", port)
}

// TestServerDiscovery checks if multiple servers can discover each other using multicast
func TestServerDiscovery(t *testing.T) {
	// Number of test servers
	numServers := 3
	serverAddresses := make([]string, numServers)
	pongChans := make([]chan string, numServers)
	ctxs := make([]context.Context, numServers)
	cancels := make([]context.CancelFunc, numServers)

	// Start multiple servers
	for i := 0; i < numServers; i++ {
		// Get a random available address for each server
		addr := getRandomAvailableAddress()

		serverAddresses[i] = addr
		pongChans[i] = make(chan string, numServers) // Buffered to prevent blocking
		ctxs[i], cancels[i] = context.WithCancel(context.Background())

		// Start listening for PINGs
		go StartMulticastListener(addr, ctxs[i])

		// Start sending PINGs
		go StartMulticastPinger(ctxs[i])

		// Start listening for PONG responses
		ListenForever(ctxs[i], pongChans[i])
	}

	// Wait for servers to discover each other
	time.Sleep(5 * time.Second)

	// Check if servers received at least one PONG message
	for i := 0; i < numServers; i++ {
		select {
		case discoveredServer := <-pongChans[i]:
			t.Logf("Server %s discovered %s", serverAddresses[i], discoveredServer)
		case <-time.After(3 * time.Second): // Timeout to prevent blocking
			t.Errorf("Server %s did not discover any peers", serverAddresses[i])
		}
	}

	// Clean up: Stop all servers
	for _, cancel := range cancels {
		cancel()
	}

	// Allow time for cleanup
	time.Sleep(2 * time.Second)
}
