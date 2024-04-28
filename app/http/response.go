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
