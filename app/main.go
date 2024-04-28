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
