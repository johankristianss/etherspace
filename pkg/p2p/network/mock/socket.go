package mock

import (
	"context"

	net "github.com/johankristianss/evrium/pkg/p2p/network"
)

type Socket interface {
	Send(msg net.Message) error
	Receive(context.Context) (net.Message, error)
}
