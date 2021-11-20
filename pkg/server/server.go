package server

import (
	"net"
	"sync"
)

type HandleFunc func(conn net.Conn)

type Server struct {
	addr 	string
	mu 			sync.Mutex
	handlers 	map[string]HandleFunc
}


func NewServer(addr string) *Server {
	return &Server{addr: addr, handlers: make(map[string]HandleFunc)}
}

func (s *Server) Register(path string, handler HandleFunc) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.handlers[path] = handler
}

func (s *Server) Start() error {
	return nil
}