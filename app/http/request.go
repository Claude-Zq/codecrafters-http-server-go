package http

const (
	MethodGet = "GET"
)

type Request struct {
	HttpVersion string
	Headers     map[string]string
	Method      string
	Path        string
}
