# [Concurrent connections](https://app.codecrafters.io/courses/http-server/stages/6)

并发支持

```go
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
```

可以看到我们的http` server`当前在同一时间只能处理一个请求。

得益于 `go` 语言对并发强有力的支持，我们只需要加一个关键字即可，像极了作弊。

```
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
```

测试：

```shell
qing@QingdeMacBook-Pro codecrafters-http-server-go % codecrafters test         
Initiating test run...

⚡ This is a turbo test run. https://codecrafters.io/turbo

Running tests. Logs should appear shortly...

Debug = true

[stage-6] Running tests for Stage #6: Concurrent connections
[stage-6] Running program
[stage-6] $ ./your_server.sh
[stage-6] Creating 3 parallel connections
[stage-6] Creating connection 1
[your_program] 2024/04/28 13:44:33 HTTP server started on port 4221
[stage-6] Creating connection 2
[stage-6] Creating connection 3
[stage-6] Sending request on 3 (status line): GET / HTTP/1.1
[stage-6] Sending request on 3: (Messages with >>> prefix are part of this log)
[stage-6] >>> GET / HTTP/1.1
[stage-6] >>> 
[stage-6] >>> 
[stage-6] >>> 
[stage-6] Received response: (Messages with >>> prefix are part of this log)
[stage-6] >>> HTTP/1.1 200 OK
[stage-6] >>> Connection: close
[stage-6] >>> Content-Length: 0
[stage-6] >>> 
[stage-6] >>> 
[stage-6] Sending request on 2 (status line): GET / HTTP/1.1
[stage-6] Sending request on 2: (Messages with >>> prefix are part of this log)
[stage-6] >>> GET / HTTP/1.1
[stage-6] >>> 
[stage-6] >>> 
[stage-6] >>> 
[stage-6] Received response: (Messages with >>> prefix are part of this log)
[stage-6] >>> HTTP/1.1 200 OK
[stage-6] >>> Connection: close
[stage-6] >>> Content-Length: 0
[stage-6] >>> 
[stage-6] >>> 
[stage-6] Sending request on 1 (status line): GET / HTTP/1.1
[stage-6] Sending request on 1: (Messages with >>> prefix are part of this log)
[stage-6] >>> GET / HTTP/1.1
[stage-6] >>> 
[stage-6] >>> 
[stage-6] >>> 
[stage-6] Received response: (Messages with >>> prefix are part of this log)
[stage-6] >>> HTTP/1.1 200 OK
[stage-6] >>> Connection: close
[stage-6] >>> Content-Length: 0
[stage-6] >>> 
[stage-6] >>> 
[stage-6] Closing connection 3
[stage-6] Closing connection 2
[stage-6] Closing connection 1
[stage-6] Test passed.
[stage-6] Terminating program
[stage-6] Program terminated successfully

```

