package http

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
)

// Handler is a function that handles a request.
type Handler func(req *Request, res *Response)

// Server is a HTTP server. It listens for incoming connections and routes them to the appropriate handler.
type Server struct {
	handlers map[string]Handler
}

// NewServer creates a new server
func NewServer() *Server {
	return &Server{
		handlers: make(map[string]Handler),
	}
}

// SetHandler sets a handler for a given path pattern.
func (s *Server) SetHandler(pattern string, handler Handler) {
	s.handlers[pattern] = handler
}

// Start starts the server on the given port.
func (s Server) Start(port int) error {
	l, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", port))
	if err != nil {
		return err
	}
	log.Printf("HTTP server started on port %d\n", port)
	for {
		conn, err := l.Accept()
		if err != nil {
			return err
		}
		s.handleConnection(conn)
	}
}

// handleConnection handles a single connection to the server.
func (s Server) handleConnection(conn net.Conn) {
	defer conn.Close()
	res := NewResponse(StatusOK)
	res.Headers[HeaderConnection] = "close"
	res.Headers[HeaderContentLength] = "0"
	req, err := readRequest(conn)
	if err == nil {
		handler, ok := s.handlers[req.Path]
		if ok {
			handler(req, res)
		} else {
			res.Status = StatusNotFound
		}
	} else {
		res.Status = StatusBadRequest
	}
	sendResponse(conn, res)
}

// readLine reads a line from the reader.
func readLine(reader *bufio.Reader) ([]byte, error) {
	var line []byte
	for {
		next, prefix, err := reader.ReadLine()
		if err != nil {
			return nil, err
		}
		line = append(line, next...)
		if !prefix {
			break
		}
	}
	return line, nil
}

// Request represents an HTTP request.
func readRequest(conn net.Conn) (*Request, error) {
	var req Request
	reader := bufio.NewReader(conn)
	requestLine, err := readLine(reader)
	if err != nil {
		return nil, err
	}
	requestParts := strings.Fields(string(requestLine))
	if len(requestParts) != 3 {
		return nil, fmt.Errorf("Malformed request: %s", string(requestLine))
	}
	req.Method = requestParts[0]
	req.Path = requestParts[1]
	req.HttpVersion = requestParts[2]
	return &req, nil
}

// sendResponse sends a response to the client.
func sendResponse(conn net.Conn, res *Response) {
	msg := strings.Builder{}
	msg.WriteString(fmt.Sprintf("HTTP/1.1 %s\r\n", res.Status))
	for key, value := range res.Headers {
		msg.WriteString(fmt.Sprintf("%s: %s\r\n", key, value))
	}
	msg.WriteString("\r\n")
	conn.Write([]byte(msg.String()))
}
