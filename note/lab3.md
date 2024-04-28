# [Respond with 200](https://app.codecrafters.io/courses/http-server/stages/3)

项目结构大改动，实现一个 http server 的主题框架

```shell
app
├── http
│   ├── header.go
│   ├── request.go
│   ├── response.go
│   └── server.go
└── main.go
```

## 定义请求题格式

```go
package http

const (
	MethodGet = "GET"
)

type Request struct {
	HttpVersion string
	Method      string
	Path        string
}
```

## 定义响应

```go
package http

const (
	StatusOK         = "200 OK"
	StatusBadRequest = "400 Bad Request"
	StatusNotFound   = "404 Not Found"
)

type Response struct {
	Headers map[string]string
	Status  string
}

func NewResponse(status string) *Response {
	return &Response{
		Headers: make(map[string]string),
		Status:  status,
	}
}
```

## 定义常用header

```go
package http

const (
	HeaderConnection    = "Connection"
	HeaderContentLength = "Content-Length"
	HeaderContentType   = "Content-Type"
)
```

## 实现server

```go
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
```

## 使用server

```go
package main

import "github.com/codecrafters-io/http-server-starter-go/app/http"

func main() {
	server := http.NewServer()
	server.SetHandler("/", func(req *http.Request, res *http.Response) {
		res.Status = http.StatusOK
	})
	if err := server.Start(4221); err != nil {
		panic(err)
	}
}
```

## 测试

```shell
[stage-3] Running tests for Stage #3: Respond with 404
[stage-3] Running program
[stage-3] $ ./your_server.sh
[stage-3] You can use the following curl command to test this locally
[stage-3] $ curl -v -X GET http://localhost:4221/grape
[stage-3] Sending request (status line): GET /grape HTTP/1.1
[stage-3] Sending request: (Messages with >>> prefix are part of this log)
[stage-3] >>> GET /grape HTTP/1.1
[stage-3] >>> Host: localhost:4221
[stage-3] >>> User-Agent: Go-http-client/1.1
[stage-3] >>> Accept-Encoding: gzip
[stage-3] >>> 
[stage-3] >>> 
[your_program] 2024/04/28 13:02:52 HTTP server started on port 4221
[stage-3] Received response: (Messages with >>> prefix are part of this log)
[stage-3] >>> HTTP/1.1 404 Not Found
[stage-3] >>> Connection: close
[stage-3] >>> Content-Length: 0
[stage-3] >>> 
[stage-3] >>> 
[stage-3] Test passed.
[stage-3] Terminating program
[stage-3] Program terminated successfully
```

