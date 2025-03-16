package network

import (
	"context"
	"fmt"
	"net"
	"strings"
	"time"
)

const (
	MulticastAddress = "224.0.0.250:9999" // Multicast IP and port for PING
	MulticastGroup   = "224.0.0.250"
	MulticastPort    = 9999
	BufferSize       = 1024
	PingMessage      = "PING"
	PongMessage      = "PONG"
)

// StartMulticastPinger sends periodic PING messages to discover servers
func StartMulticastPinger(ctx context.Context) {
	go func() {
		addr, err := net.ResolveUDPAddr("udp", MulticastAddress)
		if err != nil {
			fmt.Println("Error resolving multicast address:", err)
			return
		}

		conn, err := net.DialUDP("udp", nil, addr)
		if err != nil {
			fmt.Println("Error creating UDP connection:", err)
			return
		}
		defer conn.Close()

		ticker := time.NewTicker(2 * time.Second) // Send every 2 seconds
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				fmt.Println("Stopping multicast pinger...")
				return
			default:
				_, err := conn.Write([]byte(PingMessage))
				if err != nil {
					fmt.Println("Failed to send PING:", err)
				} else {
					fmt.Println("Sent PING")
				}
			}
			time.Sleep(2 * time.Second) // Adjust based on network needs
		}
	}()
}

// StartMulticastListener listens for PING messages and responds with PONG
func StartMulticastListener(serverAddr string, ctx context.Context) {
	go func() {
		addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", MulticastGroup, MulticastPort))
		if err != nil {
			fmt.Println("Error resolving multicast address:", err)
			return
		}

		conn, err := net.ListenMulticastUDP("udp", nil, addr)
		if err != nil {
			fmt.Println("Error joining multicast group:", err)
			return
		}
		defer conn.Close()

		buf := make([]byte, BufferSize)

		for {
			select {
			case <-ctx.Done():
				fmt.Println("Stopping multicast listener...")
				return
			default:
				conn.SetReadDeadline(time.Now().Add(1 * time.Second)) // Timeout for graceful shutdown
				n, senderAddr, err := conn.ReadFromUDP(buf)
				if err != nil {
					continue // Ignore timeout errors
				}

				message := string(buf[:n])
				fmt.Printf("Received: %s from %s\n", message, senderAddr)

				// If received PING, respond with PONG
				if strings.TrimSpace(message) == PingMessage {
					go sendPong(serverAddr, senderAddr)
				}
			}
		}
	}()
}

// sendPong responds to a PING message with a PONG
func sendPong(serverAddr string, recipient *net.UDPAddr) {
	conn, err := net.DialUDP("udp", nil, recipient)
	if err != nil {
		fmt.Println("Error sending PONG:", err)
		return
	}
	defer conn.Close()

	response := fmt.Sprintf("%s:%s", PongMessage, serverAddr)
	_, err = conn.Write([]byte(response))
	if err != nil {
		fmt.Println("Error writing PONG response:", err)
	} else {
		fmt.Printf("Sent PONG to %s\n", recipient)
	}
}

// ListenForever continuously listens for PONG responses in a separate goroutine
func ListenForever(ctx context.Context, pongChan chan<- string) {
	go func() {
		addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%d", MulticastPort))
		if err != nil {
			fmt.Println("Error resolving UDP address:", err)
			return
		}

		conn, err := net.ListenUDP("udp", addr)
		if err != nil {
			fmt.Println("Error starting UDP listener:", err)
			return
		}
		defer conn.Close()

		buf := make([]byte, BufferSize)

		for {
			select {
			case <-ctx.Done():
				fmt.Println("Stopping UDP listener...")
				return
			default:
				conn.SetReadDeadline(time.Now().Add(1 * time.Second)) // Timeout for graceful exit
				n, senderAddr, err := conn.ReadFromUDP(buf)
				if err != nil {
					continue // Ignore timeout errors
				}

				message := string(buf[:n])
				if strings.HasPrefix(message, PongMessage) {
					serverInfo := strings.TrimPrefix(message, PongMessage+":")
					fmt.Printf("Discovered server: %s (from %s)\n", serverInfo, senderAddr)
					pongChan <- serverInfo
				}
			}
		}
	}()
}
