package ghttp

import (
	"io"
	"net/http"
)

type Response interface {
	GetStatusCode() int
	GetBody() io.ReadCloser
}

type response struct {
	stdResponse *http.Response
}

func newResponse(stdResonse *http.Response) *response {
	return &response{
		stdResponse: stdResonse,
	}
}

func (r *response) GetStatusCode() int {
	return r.stdResponse.StatusCode
}

func (r *response) GetBody() io.ReadCloser {
	return r.stdResponse.Body
}
