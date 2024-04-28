package http

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"path"
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
		go s.handleConnection(conn)
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
		handler := s.getHandler(req.Path)
		if handler != nil {
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
	req := Request{Headers: make(map[string]string)}
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

	for {
		headerLine, err := readLine(reader)
		if err != nil {
			return nil, err
		} else if len(headerLine) == 0 {
			break
		}
		headerParts := strings.SplitN(string(headerLine), ": ", 2)
		if len(headerParts) != 2 {
			continue
		}
		req.Headers[headerParts[0]] = headerParts[1]
	}

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
	msg.WriteString(res.content)

	conn.Write([]byte(msg.String()))
}

func (s Server) getHandler(requestPath string) Handler {
	// TODO: Implement path pattern-matching
	handler, ok := s.handlers[requestPath]
	if ok {
		return handler
	}
	for requestPath != "/" {
		requestPath = path.Dir(requestPath)
		handler, ok := s.handlers[requestPath+"/**"]
		if ok {
			return handler
		}
	}
	return nil
}
