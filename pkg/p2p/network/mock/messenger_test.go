package mock

import (
	"context"
	"testing"
	"time"

	"github.com/go-playground/assert/v2"
	net "github.com/johankristianss/evrium/pkg/p2p/network"
)

func TestMessenger(t *testing.T) {
	n := CreateFakeNetwork()

	node1 := net.Node{Addr: "10.0.0.1:1111"}
	messenger1 := CreateMessenger(n, node1)

	node2 := net.Node{Addr: "10.0.0.2:1111"}
	messenger2 := CreateMessenger(n, node2)

	msgChan := make(chan net.Message)
	ctx := context.TODO()
	go func() {
		messenger1.ListenForever(msgChan, ctx)
	}()

	for {
		err := messenger2.Send(net.Message{From: node2, To: node1, Payload: []byte("Hello")}, context.TODO())
		if err != nil {
			time.Sleep(10 * time.Millisecond)
		} else {
			break
		}
	}

	msg := <-msgChan

	assert.Equal(t, string(msg.Payload), "Hello")
	assert.Equal(t, msg.From.Addr, "10.0.0.2:1111")
	assert.Equal(t, msg.To.Addr, "10.0.0.1:1111")

	err := messenger2.Send(net.Message{From: node2, To: node1, Payload: []byte("Hello 2")}, context.TODO())
	assert.Equal(t, err, nil)

	msg = <-msgChan
	assert.Equal(t, string(msg.Payload), "Hello 2")
}
