# [Parse header](https://app.codecrafters.io/courses/http-server/stages/5)

本次lab需要解析请求的 `header` 中的 `User-Agent`



`app/http/header.go` 先加入该 header 的 key。

```go
package http

const (
	HeaderConnection    = "Connection"
	HeaderContentLength = "Content-Length"
	HeaderContentType   = "Content-Type"
	HeaderUserAgent     = "User-Agent" // add
)
```

`app/http/request.go` 加入 header.

```go
package http

const (
	MethodGet = "GET"
)

type Request struct {
	HttpVersion string
	Headers     map[string]string // add
	Method      string
	Path        string
}
```

`app/http/server.go` 更新读取请求报文的逻辑，在一个循环中读取 `headers`.

```go

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
```

`app/main.go` 只需要简单的从请求头冲获取header并放置到相应体中即可

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
	server.SetHandler("/user-agent", func(req *http.Request, res *http.Response) {
		res.SetContent("text/plain", req.Headers[http.HeaderUserAgent])
	})

	if err := server.Start(4221); err != nil {
		panic(err)
	}
}
```

测试结果：

```shell
[stage-5] Running tests for Stage #5: Parse headers
[stage-5] Running program
[stage-5] $ ./your_server.sh
[stage-5] You can use the following curl command to test this locally
[stage-5] $ curl -v -X GET http://localhost:4221/user-agent -H "User-Agent: orange/grape"
[stage-5] Sending request (status line): GET /user-agent HTTP/1.1
[stage-5] Sending request: (Messages with >>> prefix are part of this log)
[stage-5] >>> GET /user-agent HTTP/1.1
[stage-5] >>> Host: localhost:4221
[stage-5] >>> User-Agent: orange/grape
[stage-5] >>> Accept-Encoding: gzip
[stage-5] >>> 
[stage-5] >>> 
[your_program] 2024/04/28 13:35:28 HTTP server started on port 4221
[stage-5] Received response: (Messages with >>> prefix are part of this log)
[stage-5] >>> HTTP/1.1 200 OK
[stage-5] >>> Connection: close
[stage-5] >>> Content-Length: 12
[stage-5] >>> Content-Type: text/plain
[stage-5] >>> 
[stage-5] >>> orange/grape
[stage-5] Test passed.
[stage-5] Terminating program
[stage-5] Program terminated successfully
```

