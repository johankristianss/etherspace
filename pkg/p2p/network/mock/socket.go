package mock

import (
	"context"

	net "github.com/johankristianss/etherspace/pkg/p2p/network"
)

type Socket interface {
	Send(msg net.Message) error
	Receive(context.Context) (net.Message, error)
}
