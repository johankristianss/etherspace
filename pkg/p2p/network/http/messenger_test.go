package network

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	p2pnet "github.com/johankristianss/etherspace/pkg/p2p/network"
)

// TestSend checks if the HTTPMessenger's Send method correctly sends messages
func TestSend(t *testing.T) {
	// Create a mock HTTP server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST method, got %s", r.Method)
		}

		var msg p2pnet.Message
		if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
			t.Errorf("Failed to decode request body: %v", err)
		}

		// Validate message fields
		if msg.ID != "1" || msg.From.Name != "Sender" || msg.To.Name != "Receiver" {
			t.Errorf("Unexpected message received: %+v", msg)
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer mockServer.Close()

	// Create a Messenger instance
	messenger := NewHTTPMessenger("localhost:8080")

	// Create a test message
	testMsg := p2pnet.Message{
		ID:      "1",
		From:    p2pnet.Node{Name: "Sender", Addr: "localhost:8081"},
		To:      p2pnet.Node{Name: "Receiver", Addr: mockServer.Listener.Addr().String()},
		Type:    1,
		Payload: []byte("Test Payload"),
	}

	// Call the Send method
	err := messenger.Send(testMsg, context.Background())
	if err != nil {
		t.Errorf("Send failed: %v", err)
	}
}

// TestListenForever checks if the server correctly receives messages
func TestListenForever(t *testing.T) {
	// Create a Messenger instance that listens on :8081
	messenger := NewHTTPMessenger(":8081")

	// Message channel
	msgChan := make(chan p2pnet.Message, 1)

	// Context to manage the server
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start the listener in a separate goroutine
	go func() {
		_ = messenger.ListenForever(msgChan, ctx)
	}()
	time.Sleep(1 * time.Second) // Give the server time to start

	// Create a test message
	testMsg := p2pnet.Message{
		ID:      "2",
		From:    p2pnet.Node{Name: "Sender", Addr: "localhost:8082"},
		To:      p2pnet.Node{Name: "Receiver", Addr: "localhost:8081"},
		Type:    2,
		Payload: []byte("Hello, World"),
	}

	// Use the Send method to send the message
	sendMessenger := NewHTTPMessenger("localhost:8082") // Sender instance
	err := sendMessenger.Send(testMsg, context.Background())
	if err != nil {
		t.Fatalf("Failed to send message: %v", err)
	}

	// Verify that the message is received in the channel
	select {
	case receivedMsg := <-msgChan:
		if receivedMsg.ID != testMsg.ID || string(receivedMsg.Payload) != string(testMsg.Payload) {
			t.Errorf("Received incorrect message: %+v", receivedMsg)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("Message was not received in time")
	}
}
