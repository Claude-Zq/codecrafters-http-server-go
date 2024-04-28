# [Bind to a port](https://app.codecrafters.io/courses/http-server/stages/1)

取消注释即可

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

	_, err = l.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}
}

```

测试结果：

```shell
qing@QingdeMacBook-Pro note % codecrafters test
Initiating test run...

⚡ This is a turbo test run. https://codecrafters.io/turbo

Running tests. Logs should appear shortly...

Debug = true

[stage-1] Running tests for Stage #1: Bind to a port
[stage-1] Running program
[stage-1] $ ./your_server.sh
[stage-1] Connecting to localhost:4221 using TCP
[your_program] Logs from your program will appear here!
[stage-1] Success! Closing connection
[stage-1] Test passed.
[stage-1] Terminating program
[stage-1] Program terminated successfully

All tests passed. Congrats!
```

