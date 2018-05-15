package proxy

import (
	"net"

	log "github.com/sirupsen/logrus"
)

//ServerProxifier inteface for proxy runner
type ServerProxifier interface {
	ListenAndServe(network, address string) error
	Serve(li net.Listener) error
}

//Server container for data
type Server struct {
	Logger *log.Logger
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
	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			s.Logger.Error(err)
			continue
		}
		c := NewConnection(conn, s.Logger)
		go c.Serve()
	}
}
