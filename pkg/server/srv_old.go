package server

// import (
// 	"bufio"
// 	"bytes"
// 	"errors"
// 	"io"
// 	"log"
// 	"net"
// 	"net/url"
// 	"strconv"
// 	"strings"
// 	"sync"
// )

// var rn = "\r\n"
// var errConnectionBroken = errors.New("connection is broken")
// var errHttpSupport = errors.New("your browser can't support HTTP version 1.1")

// type HandleFunc func(req *Request)

// // type RequestFunc func(req Request)
// type Request struct {
// 	Conn        net.Conn
// 	QueryParams url.Values

// 	Body       []byte
// 	Headers    map[string]string
// 	PathParams map[string]string
// }

// type Server struct {
// 	addr     string
// 	mu       sync.Mutex
// 	handlers map[string]HandleFunc
// }

// func NewServer(addr string) *Server {
// 	return &Server{addr: addr, handlers: make(map[string]HandleFunc)}
// }

// func HeaderShortcut(body string) []byte {
// 	return []byte(
// 		"HTTP/1.1 200 OK" + rn +
// 			"Content-Length:" + strconv.Itoa(len(body)) + rn +
// 			"Content-Type:text/html" + rn +
// 			"Connection:close" + rn + rn +
// 			string(body),
// 	)
// }

// func (s *Server) Register(path string, handler HandleFunc) {
// 	s.mu.Lock()
// 	defer s.mu.Unlock()
// 	s.handlers[path] = handler
// }

// func (s *Server) Start() error {
// 	listener, err := net.Listen("tcp", s.addr)
// 	if err != nil {
// 		log.Println("=========| Start(): listener crash |=========")
// 		return err
// 	}

// 	defer func() {
// 		if cerr := listener.Close(); cerr != nil {
// 			if err == nil {
// 				err = cerr
// 				return
// 			}
// 			log.Println(cerr)
// 		}
// 	}()

// 	log.Println("=========| Start(): Start server |=========")

// 	for {
// 		conn, err := listener.Accept()
// 		if err != nil {
// 			log.Println("=========| Start(): ", err, " |=========")
// 			continue
// 		}

// 		go s.handle(conn)
// 	}
// }

// func (s *Server) handle(conn net.Conn) (err error) {
// 	defer func() {
// 		if cerr := conn.Close(); cerr != nil {
// 			if err == nil {
// 				err = cerr
// 				return
// 			}
// 			log.Println("=========| handle(): ", err, "|=========")
// 		}
// 	}()

// 	// GET DATA
// 	log.Println("=========| handle(): New Connection |=========")
// 	var req Request
// 	buf := make([]byte, 4096)
// 	n, err := conn.Read(buf)

// 	if err != io.EOF {
// 		log.Printf("%s", buf[:n])
// 	}
// 	if err != nil {
// 		return err
// 	}
// 	// --


// 	// PARSE PATH, BROWSER VERSION
// 	// * re   -- request line
// 	// * rle  -- request line end
// 	// * rld  -- request line delim
// 	data := buf[:n]

// 	rld := []byte(rn)
// 	rle := bytes.Index(data, rld)

// 	if rle == -1 {
// 		return errConnectionBroken
// 	}

// 	rl := string(data[:rle])
// 	parts := strings.Split(rl, " ")

// 	if len(parts) != 3 {
// 		return errConnectionBroken
// 	}

// 	path, version := parts[1], parts[2]

// 	if version != "HTTP/1.1" {
// 		return errHttpSupport
// 	}
// 	// --


// 	// SAVE HEADER
// 	// 
// 	headerData := string(data[rle+2:bytes.Index(data, []byte(rn+rn))])
// 	scanner := bufio.NewScanner(strings.NewReader(headerData))

// 	headersMap := make(map[string]string)
// 	for scanner.Scan() {
// 		t := scanner.Text()
// 		headersMap[t[:strings.Index(t, ":")]] = t[strings.Index(t, " ")+1:]
// 	}

// 	// log.Println(headersMap)
// 	req.Headers = headersMap
// 	// --


// 	// SAVE BODY
// 	bodyData := data[bytes.Index(data, []byte(rn+rn))+4:]
// 	req.Body = bodyData
// 	log.Println(req.Body, bodyData)

// 	// PARSE URL REQUEST
// 	// Decode symbols not range on ascii table
// 	decoded, _ := url.PathUnescape(path)
// 	if err != nil {
// 		log.Print(err)
// 		return
// 	}

// 	url, _ := url.ParseRequestURI(decoded)
// 	if err != nil {
// 		log.Print(err)
// 		return
// 	}

// 	req.Conn = conn
// 	req.QueryParams = url.Query()

// 	handler := func(req *Request) { req.Conn.Close() }

// 	s.mu.Lock()
// 	pathParam, hFunc := s.FindRoute(url.Path)

// 	if hFunc != nil {
// 		req.PathParams = pathParam
// 		handler = hFunc
// 	}

// 	s.mu.Unlock()
// 	handler(&req)

// 	return nil
// }

// //
// func (s *Server) FindRoute(path string) (map[string]string, HandleFunc) {
// 	tempIndx := 0
// 	mapOFparams := make(map[string]string)
// 	registRoutes := make([]string, len(s.handlers))

// 	for i := range s.handlers {
// 		registRoutes[tempIndx] = i
// 		tempIndx++
// 	}

// 	for i := 0; i < len(registRoutes); i++ {
// 		flag := false
// 		tempRoute := registRoutes[i]
// 		cliParts := strings.Split(path, "/")
// 		serverParts := strings.Split(tempRoute, "/")

// 		log.Println(serverParts) // []string
// 		log.Println(cliParts)    // []string

// 		for i, j := range serverParts {
// 			if j != "" {
// 				fIndx := j[0:1]
// 				lIndx := j[len(j)-1:]

// 				if fIndx == "{" && lIndx == "}" {
// 					mapOFparams[j[1:len(j)-1]] = cliParts[i]
// 					flag = true

// 				} else if cliParts[i] != j {
// 					// log.Println(j, v, "|||", f,l, mapOFparams, " <<<<<<<========================================== JV")

// 					splitStr := strings.Split(j, "{")

// 					if len(splitStr) > 0 {
// 						key := splitStr[1][:len(splitStr[1])-1]
// 						mapOFparams[key] = cliParts[i][len(splitStr[0]):]
// 						flag = true
// 					} else {
// 						flag = false
// 						break
// 					}
// 				}
// 				flag = true
// 			}
// 		}

// 		if flag {
// 			if function, status := s.handlers[tempRoute]; status {
// 				log.Println(function, status)
// 				return mapOFparams, function
// 			}
// 			break
// 		}
// 	}

// 	return nil, nil
// }
