package server

import (
	"bytes"
	"errors"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
	"sync"
)

type HandleFunc func(conn net.Conn)

type Server struct {
	addr     string
	mu       sync.Mutex
	handlers map[string]HandleFunc
}

var rn = "\r\n"

func HeaderShortcut(body string) []byte {
	return []byte(
		"HTTP/1.1 200 OK" + rn +
			"Content-Length:" + strconv.Itoa(len(body)) + rn +
			"Content-Type:text/html" + rn +
			"Connection:close" + rn + rn +
			string(body),
	)
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
	listener, err := net.Listen("tcp", s.addr)

	if err != nil {
		log.Println("listener crash")
		return err
	}

	defer func() {
		if cerr := listener.Close(); cerr != nil {
			if err == nil {
				err = cerr
				return
			}
			log.Println(cerr)
		}
	}()

	log.Println("start")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		go s.handle(conn)
	}
}


func (s *Server) handle(conn net.Conn) (err error) {
	defer func() {
		if cerr := conn.Close(); cerr != nil {
			if err == nil {
				err = cerr
				return
			}
			log.Println(err)
		}
	}()

	buf := make([]byte, 4096)
	n, err := conn.Read(buf)

	log.Println(" >>>>> NEW CONNECTION <<<<< BYTES => ", n)
	
	// Recive header data from client
	if err != io.EOF {
		log.Printf("%s", buf[:n])
	}
	if err != nil {
		return err
	}


	// Set header data
	data := buf[:n]
	reqiestLineDelim := []byte(rn)
	reqiestLineEnd := bytes.Index(data, reqiestLineDelim)
	
	if reqiestLineEnd == -1 {
		return errors.New("Connection is broken")
	}

	requserLine := string(data[:reqiestLineEnd])
	parts := strings.Split(requserLine, " ")
	
	if (len(parts) != 3 ) {
		return errors.New("Connection is broken")
	}

	route := s.handlers[parts[1]]

	// return 
	if route != nil {
		for pathPresent, handler := range s.handlers {
			if pathPresent == parts[1] {
				handler(conn)
		}
	}}

	return nil
}
