## [Get a file](https://app.codecrafters.io/courses/http-server/stages/7)

支持文件获取，非常简单，从命令行获取文件服务的根目录目录，然后写一个handler即可。

```go
package main

import (
	"flag"
	"os"
	"path"
	"strings"

	"github.com/Claude-Zq/http-server-starter-go/app/http"
)

func main() {
	var directory string
	flag.StringVar(&directory, "directory", "files", "path to the files directory")
	flag.Parse()

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
		if content, err := os.ReadFile(filePath); err != nil {
			res.Status = http.StatusNotFound
		} else {
			res.SetContent("application/octet-stream", string(content))
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

[stage-7] Running tests for Stage #7: Get a file
[stage-7] Running program
[stage-7] $ ./your_server.sh --directory /tmp/data/codecrafters.io/http-server-tester/
[stage-7] Testing existing file
[stage-7] Creating file orange_raspberry_orange_banana in /tmp/data/codecrafters.io/http-server-tester/
[stage-7] File Content: "mango grape banana banana pineapple pineapple orange pineapple"
[stage-7] You can use the following curl command to test this locally
[stage-7] $ curl -v -X GET http://localhost:4221/files/orange_raspberry_orange_banana
[stage-7] Sending request (status line): GET /files/orange_raspberry_orange_banana HTTP/1.1
[stage-7] Sending request: (Messages with >>> prefix are part of this log)
[stage-7] >>> GET /files/orange_raspberry_orange_banana HTTP/1.1
[stage-7] >>> Host: localhost:4221
[stage-7] >>> User-Agent: Go-http-client/1.1
[stage-7] >>> Accept-Encoding: gzip
[stage-7] >>> 
[stage-7] >>> 
[your_program] 2024/04/28 14:24:00 HTTP server started on port 4221
[stage-7] Received response: (Messages with >>> prefix are part of this log)
[stage-7] >>> HTTP/1.1 200 OK
[stage-7] >>> Connection: close
[stage-7] >>> Content-Length: 62
[stage-7] >>> Content-Type: application/octet-stream
[stage-7] >>> 
[stage-7] >>> mango grape banana banana pineapple pineapple orange pineapple
[stage-7] Testing non existent file returns 404
[stage-7] You can use the following curl command to test this locally
[stage-7] $ curl -v -X GET http://localhost:4221/files/non-existentgrape_orange_grape_grape
[stage-7] Sending request (status line): GET /files/non-existentgrape_orange_grape_grape HTTP/1.1
[stage-7] Sending request: (Messages with >>> prefix are part of this log)
[stage-7] >>> GET /files/non-existentgrape_orange_grape_grape HTTP/1.1
[stage-7] >>> Host: localhost:4221
[stage-7] >>> User-Agent: Go-http-client/1.1
[stage-7] >>> Accept-Encoding: gzip
[stage-7] >>> 
[stage-7] >>> 
[stage-7] Received response: (Messages with >>> prefix are part of this log)
[stage-7] >>> HTTP/1.1 404 Not Found
[stage-7] >>> Connection: close
[stage-7] >>> Content-Length: 0
[stage-7] >>> 
[stage-7] >>> 
[stage-7] Test passed.
[stage-7] Terminating program
[your_program] /files/orange_raspberry_orange_banana/files/non-existentgrape_orange_grape_grape
[stage-7] Program terminated successfully
```

