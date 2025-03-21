package libp2p

import (
	"context"

	net "github.com/johankristianss/evrium/pkg/p2p/network"
	log "github.com/sirupsen/logrus"
)

type Server struct {
	port      int
	messenger *Messenger
	ctx       context.Context
	cancel    context.CancelFunc
	msgChan   chan net.Message
	handler   func(msg net.Message)
}

func CreateServer(port int, handler func(msg net.Message)) (*Server, error) {
	messenger, err := CreateMessenger(port, "")
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())

	msgChan := make(chan net.Message, 1000)

	s := &Server{}
	s.port = port
	s.messenger = messenger
	s.ctx = ctx
	s.cancel = cancel
	s.msgChan = msgChan
	s.handler = handler

	return s, nil
}

func (s *Server) Shutdown() {
	s.cancel()
}

func (s *Server) ServerForever() {
	go func() {
		s.messenger.ListenForever(s.msgChan, s.ctx)
	}()

	go func() {
		select {
		case <-s.ctx.Done():
			log.Info("Server shutting down")
			return
		case msg := <-s.msgChan:
			log.Info("Received message: ", msg)
			s.handler(msg)
		}
	}()
}
