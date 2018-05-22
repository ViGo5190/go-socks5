package proxy

import (
	"context"
	"net"
	"sync"

	"github.com/rs/zerolog/log"
)

//ServerProxifier inteface for proxy runner
type ServerProxifier interface {
	ListenAndServe(network, address string) error
	Serve(li net.Listener) error
}

//Server container for data
type Server struct {
	ctx    context.Context
	done   context.CancelFunc
	closed bool
	wg     sync.WaitGroup
}

//ListenAndServe create listener and serve
func (s *Server) ListenAndServe(network, address string) error {
	listener, err := net.Listen(network, address)
	if err != nil {
		return err
	}
	s.Serve(listener)
	return nil
}

//Serve serve listener -> connections
func (s *Server) Serve(listener net.Listener) {
	s.ctx, s.done = context.WithCancel(context.Background())
	defer listener.Close()

	newConns := make(chan net.Conn)

	go func(l net.Listener, ss *Server) {
		for {
			c, err := l.Accept()
			if err != nil {
				if !ss.closed {
					log.Error().Msgf("Error of accepting connections: %v", err)
				}
				return
			}
			newConns <- c
		}
	}(listener, s)

	for {
		select {
		case <-s.ctx.Done():
			s.wg.Wait()
			return
		case conn := <-newConns:
			c := NewConnection(conn)
			s.wg.Add(1)
			go c.Serve(&s.wg)
		}
	}
}

//Stop server
func (s *Server) Stop() {
	s.closed = true
	s.done()
}
