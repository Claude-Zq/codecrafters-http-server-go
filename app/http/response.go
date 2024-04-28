package http

const (
	StatusOK         = "200 OK"
	StatusBadRequest = "400 Bad Request"
	StatusNotFound   = "404 Not Found"
)

type Response struct {
	Headers map[string]string
	Status  string
}

func NewResponse(status string) *Response {
	return &Response{
		Headers: make(map[string]string),
		Status:  status,
	}
}
