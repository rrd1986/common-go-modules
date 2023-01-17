package serviceClient

import (
	"net/http"
)

type BasicResponse interface {
	StatusCode() int
	Body() []byte
}

type HeaderedResponse interface {
	BasicResponse
	Header() http.Header
}
