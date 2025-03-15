package mock

import (
	"errors"
	"sync"

	net "github.com/johankristianss/etherspace/pkg/p2p/network"
)

type FakeNetwork struct {
	Hosts map[string]*FakeSocket
	mutex sync.Mutex
}

func CreateFakeNetwork() *FakeNetwork {
	return &FakeNetwork{Hosts: make(map[string]*FakeSocket)}
}

func (n *FakeNetwork) Listen(addr string) (Socket, error) {
	n.mutex.Lock()
	defer n.mutex.Unlock()
	socket := &FakeSocket{conn: make(chan net.Message, 1000)}
	n.Hosts[addr] = socket
	return socket, nil
}

func (n *FakeNetwork) Dial(addr string) (Socket, error) {
	n.mutex.Lock()
	defer n.mutex.Unlock()
	if _, ok := n.Hosts[addr]; !ok {
		return nil, errors.New("No such host " + addr)
	}
	return n.Hosts[addr], nil
}
