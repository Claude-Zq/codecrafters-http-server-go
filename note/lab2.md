## [Respond with 200](https://app.codecrafters.io/courses/http-server/stages/2)

有手就行，直接从连接里面读取请求，并发送相应即可。

```go
package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	conn, err := l.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}
	defer conn.Close()

	msg := make([]byte, 1024)
	_, err = conn.Read(msg)
	_, err = conn.Write([]byte("HTTP/1.1 200 OK \r\n\r\n"))
	if err != nil {
		fmt.Println("Error reading message: ", err.Error())
		os.Exit(1)
	}
}
```

测试

```shell
qing@QingdeMacBook-Pro note % codecrafters test
Initiating test run...

⚡ This is a turbo test run. https://codecrafters.io/turbo

Running tests. Logs should appear shortly...

Debug = true

[stage-2] Running tests for Stage #2: Respond with 200
[stage-2] Running program
[stage-2] $ ./your_server.sh
[stage-2] You can use the following curl command to test this locally
[stage-2] $ curl -v -X GET http://localhost:4221/
[stage-2] Sending request (status line): GET / HTTP/1.1
[stage-2] Sending request: (Messages with >>> prefix are part of this log)
[stage-2] >>> GET / HTTP/1.1
[stage-2] >>> Host: localhost:4221
[stage-2] >>> User-Agent: Go-http-client/1.1
[stage-2] >>> Accept-Encoding: gzip
[stage-2] >>> 
[stage-2] >>> 
[your_program] Logs from your program will appear here!
[stage-2] Received response: (Messages with >>> prefix are part of this log)
[stage-2] >>> HTTP/1.1 200 OK 
[stage-2] >>> Connection: close
[stage-2] >>> 
[stage-2] >>> 
[stage-2] Test passed.
[stage-2] Terminating program
[stage-2] Program terminated successfully
```

