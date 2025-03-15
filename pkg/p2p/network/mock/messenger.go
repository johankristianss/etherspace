package mock

import (
	"context"

	net "github.com/johankristianss/etherspace/pkg/p2p/network"
)

type MockMessenger struct {
	network Network
	node    net.Node
}

func CreateMessenger(network Network, node net.Node) *MockMessenger {
	return &MockMessenger{network: network, node: node}
}

func (m *MockMessenger) Send(msg net.Message, ctx context.Context) error {
	socket, err := m.network.Dial(msg.To.String())
	if err != nil {
		return err
	}
	msg.From = m.node
	return socket.Send(msg)
}

func (m *MockMessenger) ListenForever(msgChan chan net.Message, ctx context.Context) error {
	socket, err := m.network.Listen(m.node.String())
	if err != nil {
		return err
	}

	for {
		msg, _ := socket.Receive(ctx)
		select {
		case <-ctx.Done():
			return nil
		default:
			msgChan <- msg
		}
	}
}
