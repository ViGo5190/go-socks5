package proxy

import (
	"net"

	log "github.com/sirupsen/logrus"
)

type ServerProxifier interface {
	ListenAndServe(network, address string) error
	Serve(li net.Listener) error
}

type Server struct {
	Logger *log.Logger
}

func New(logger *log.Logger) *Server {
	return &Server{
		Logger: logger,
	}
}

func (s *Server) ListenAndServe(network, address string) error {
	listener, err := net.Listen(network, address)
	defer listener.Close()
	if err != nil {
		return err
	}
	return s.Serve(listener)
}

func (s *Server) Serve(listener net.Listener) error {
	for {
		conn, err := listener.Accept()
		if err != nil {
			s.Logger.Error(err)
			//return err
		}
		c := NewConnection(conn, s.Logger)
		go c.Serve()
		//go s.ServeConn(conn)
	}
}
