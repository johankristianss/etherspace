package mock

import (
	"context"

	net "github.com/johankristianss/evrium/pkg/p2p/network"
)

type FakeSocket struct {
	conn chan net.Message
}

func (socket *FakeSocket) Send(msg net.Message) error {
	socket.conn <- msg
	return nil
}

func (socket *FakeSocket) Receive(ctx context.Context) (net.Message, error) {
	select {
	case msg := <-socket.conn:
		return msg, nil
	case <-ctx.Done():
		return net.Message{}, ctx.Err()
	}
}
