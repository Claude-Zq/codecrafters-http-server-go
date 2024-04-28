package http

const (
	MethodGet = "GET"
)

type Request struct {
	HttpVersion string
	Method      string
	Path        string
}
