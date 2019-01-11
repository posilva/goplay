package server

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"time"
)

// Server holds the data to run the server
type Server struct {
	port     int
	host     string
	connType string
	listener net.Listener
	errChan  chan error
	running  bool
}

// New creates a new server
func New() *Server {
	s := &Server{
		port:     9000,
		host:     "localhost",
		connType: "tcp",
	}
	return s
}

// Listen creates a socket listener
func (s *Server) Listen() error {
	var err error
	s.listener, err = net.Listen(s.connType, s.host+":"+strconv.Itoa(s.port))
	if err != nil {
		return fmt.Errorf("failed creating the listener %v", err.Error())
	}
	fmt.Printf("Server is listening on port %v \n", s.port)
	return nil
}

// Start accepting connections
func (s *Server) Start() int {
	defer s.listener.Close()
	s.errChan = make(chan error)
	s.running = true
	for s.running {
		conn, err := s.listener.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			return 1
		}
		go s.handleErrors()
		go s.handleRequest(conn)
	}
	return 0
}

// Stop the server
func (s *Server) Stop() {
	s.running = false
}
func (s *Server) handleRequest(conn net.Conn) {
	defer conn.Close()
	timeoutDuration := 5 * time.Second
	buf := make([]byte, 1024)
	log.Printf("handling connection from %v \n", conn.RemoteAddr())
	for {
		err := conn.SetReadDeadline(time.Now().Add(timeoutDuration))
		if err != nil {
			s.errChan <- fmt.Errorf("failed to set connection read deadline %v", err)
			return
		}
		_, err = conn.Read(buf)
		if nerr, ok := err.(net.Error); ok && nerr.Timeout() {
			s.errChan <- fmt.Errorf("read deadline reached")
		} else {
			if err != nil {
				s.errChan <- fmt.Errorf("failed to read from socket %v", err.Error())
				return
			} else {
				fmt.Printf("%s", string(buf))
			}
		}
	}
}
func (s *Server) handleErrors() {
	for {
		select {
		case err := <-s.errChan:
			fmt.Printf("client error received: %v\n", err)
		}
	}
}
