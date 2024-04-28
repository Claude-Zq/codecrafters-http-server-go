## [Respond with content](https://app.codecrafters.io/courses/http-server/stages/4)

本次lab需要根据请求的path返回响应的请求，有两个需要解决的问题：

如何实现路径匹配？即`/echo/abc` `/echo/bcd` 如何映射到同一个handler上？

如何获取路径参数：如何获取 `/echo/abc` `/echo/bcd` 中的 `abc` `bcd`。

正规一点的做法是需要 `Trie`树的，可以参考：[Gee](https://geektutu.com/post/gee-day3.html)，但这里为了简化，直接使用编程技巧解决：



`app/http/response.go` 加入设置响应报文的方法

```go
package http

import "strconv"

const (
	StatusOK         = "200 OK"
	StatusBadRequest = "400 Bad Request"
	StatusNotFound   = "404 Not Found"
)

type Response struct {
	Headers map[string]string
	Status  string

	content string
}

func NewResponse(status string) *Response {
	return &Response{
		Headers: make(map[string]string),
		Status:  status,
	}
}

func (res *Response) SetContent(contentType, contentBody string) {
	res.content = contentBody
	res.Headers[HeaderContentLength] = strconv.Itoa(len(contentBody))

	res.Headers[HeaderContentType] = contentType
}

```

`app/http/server.go`修改匹配handler的逻辑

```go
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
```

`main.go `加入新的 `handler`

```go
package main

import (
	"strings"

	"github.com/codecrafters-io/http-server-starter-go/app/http"
)

func main() {
	server := http.NewServer()
	server.SetHandler("/", func(req *http.Request, res *http.Response) {})
	server.SetHandler("/echo/**", func(req *http.Request, res *http.Response) {
		index := strings.Index(req.Path, "/echo/")
		if index > -1 {
			res.SetContent("text/plain", req.Path[index+6:])
		}
	})

	if err := server.Start(4221); err != nil {
		panic(err)
	}
}

```

测试

```shell
qing@QingdeMacBook-Pro codecrafters-http-server-go % codecrafters test         
Initiating test run...

⚡ This is a turbo test run. https://codecrafters.io/turbo

Running tests. Logs should appear shortly...

Debug = true

[stage-4] Running tests for Stage #4: Respond with content
[stage-4] Running program
[stage-4] $ ./your_server.sh
[stage-4] You can use the following curl command to test this locally
[stage-4] $ curl -v -X GET http://localhost:4221/echo/pineapple
[stage-4] Sending request (status line): GET /echo/pineapple HTTP/1.1
[stage-4] Sending request: (Messages with >>> prefix are part of this log)
[stage-4] >>> GET /echo/pineapple HTTP/1.1
[stage-4] >>> Host: localhost:4221
[stage-4] >>> User-Agent: Go-http-client/1.1
[stage-4] >>> Accept-Encoding: gzip
[stage-4] >>> 
[stage-4] >>> 
[your_program] 2024/04/28 13:15:02 HTTP server started on port 4221
[stage-4] Received response: (Messages with >>> prefix are part of this log)
[stage-4] >>> HTTP/1.1 200 OK
[stage-4] >>> Connection: close
[stage-4] >>> Content-Length: 9
[stage-4] >>> Content-Type: text/plain
[stage-4] >>> 
[stage-4] >>> pineapple
[stage-4] Test passed.
[stage-4] Terminating program
[stage-4] Program terminated successfully
```

