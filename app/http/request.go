package http

const (
	MethodGet  = "GET"
	MethodPost = "POST"
)

type Request struct {
	HttpVersion string
	Headers     map[string]string
	Method      string
	Path        string
	Body        string
}
