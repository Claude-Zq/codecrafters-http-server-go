## [Post a file](https://app.codecrafters.io/courses/http-server/stages/8)

支持上传文件

`app/http/request.go` 引入`post`方法，`Request` 加入body

```go
package http

const (
	MethodGet  = "GET"
	MethodPost = "POST" // add
)

type Request struct {
	HttpVersion string
	Headers     map[string]string
	Method      string
	Path        string
	Body        string // add
}
```

`app/http/server.go` 读取请求报文中的body.

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

	if req.Headers[HeaderContentLength] == "" {
		return &req, nil
	}
	contentLength, err := strconv.Atoi(req.Headers[HeaderContentLength])
	if err != nil {
		return nil, err
	}
	body := strings.Builder{}
	for i := 0; i < contentLength; i++ {
		b, err := reader.ReadByte()
		if err != nil {
			return nil, err
		}
		body.WriteByte(b)
	}
	req.Body = body.String()
	return &req, err

}
```

`app/http/response.go` 

```go
const (
	StatusOK                = "200 OK"
	StatusCreated           = "201 Created" // add
	StatusBadRequest        = "400 Bad Request"
	StatusNotFound          = "404 Not Found"
	StatusInternalServerErr = "500 Internal Server Error" // add
)
```

更新 `main.go`

```go
package main

import (
	"flag"
	"os"
	"path"
	"strings"

	"github.com/Claude-Zq/codecrafters-http-server-go/app/http"
)

func main() {
	var directory string
	flag.StringVar(&directory, "directory", "files", "path to the files directory")
	flag.Parse()
	if err := os.MkdirAll(directory, 0777); err != nil {
		panic(err)
	}

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

	server.SetHandler("/files/**", func(req *http.Request, res *http.Response) {
		filePath := path.Join(directory, req.Path[len("/files/"):])
		switch req.Method {
		case http.MethodGet:
			if content, err := os.ReadFile(filePath); err != nil {
				res.Status = http.StatusNotFound
			} else {
				res.SetContent("application/octet-stream", string(content))
			}
		case http.MethodPost:
			if err := os.WriteFile(filePath, []byte(req.Body), 0644); err != nil {
				res.Status = http.StatusInternalServerErr
			} else {
				res.Status = http.StatusCreated
			}
		}
	})

	if err := server.Start(4221); err != nil {
		panic(err)
	}
}
```

测试：

```shell
qing@QingdeMacBook-Pro http % codecrafters test         
Initiating test run...

⚡ This is a turbo test run. https://codecrafters.io/turbo

Running tests. Logs should appear shortly...

Debug = true

[stage-8] Running tests for Stage #8: Post a file
[stage-8] Running program
[stage-8] $ ./your_server.sh --directory /tmp/data/codecrafters.io/http-server-tester/
[stage-8] You can use the following curl command to test this locally
[stage-8] $ curl -v -X POST http://localhost:4221/files/raspberry_pineapple_orange_banana -d 'strawberry orange pineapple apple grape apple strawberry orange'
[stage-8] Sending request (status line): POST /files/raspberry_pineapple_orange_banana HTTP/1.1
[stage-8] Sending request: (Messages with >>> prefix are part of this log)
[stage-8] >>> POST /files/raspberry_pineapple_orange_banana HTTP/1.1
[stage-8] >>> Host: localhost:4221
[stage-8] >>> User-Agent: Go-http-client/1.1
[stage-8] >>> Content-Length: 63
[stage-8] >>> Accept-Encoding: gzip
[stage-8] >>> 
[stage-8] >>> strawberry orange pineapple apple grape apple strawberry orange
[your_program] 2024/04/28 14:39:31 HTTP server started on port 4221
[stage-8] Received response: (Messages with >>> prefix are part of this log)
[stage-8] >>> HTTP/1.1 201 Created
[stage-8] >>> Connection: close
[stage-8] >>> Content-Length: 0
[stage-8] >>> 
[stage-8] >>> 
[stage-8] Validating file `raspberry_pineapple_orange_banana` exists on disk
[stage-8] Validating file `raspberry_pineapple_orange_banana` content
[stage-8] Test passed.
[stage-8] Terminating program
[stage-8] Program terminated successfully
```

