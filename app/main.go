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
