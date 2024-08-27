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
