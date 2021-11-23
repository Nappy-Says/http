package server

import (
	"io"
	"log"
	"net"
	"strconv"
	"sync"
)

type HandleFunc func(conn net.Conn)

type Server struct {
	addr 	string
	mu 			sync.Mutex
	handlers 	map[string]HandleFunc
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

	defer func()  {
		if cerr := listener.Close(); cerr != nil {
			if err == nil {
				err = cerr
				return
			}
			log.Println(cerr)
		}
	} ()

	log.Println("start")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		err = handle(conn)
		if err != nil {
			log.Println(err)
			continue
		}
	}
}



func handle(conn net.Conn) (err error) {
	defer func() {
		if cerr := conn.Close(); cerr != nil {
			if err == nil {
				err = cerr
				return
			}
			log.Println(err)
		}
	} ()


	buf := make([]byte, 4096)

	for {
		n, err := conn.Read(buf)

		log.Println(" >>>>> NEW CONNECTION <<<<< BYTES => ", n)
		
		if err != io.EOF {
			log.Printf("%s", buf[:n])
			return nil
		}
		if err != nil {
			return err
		}
	}
}
