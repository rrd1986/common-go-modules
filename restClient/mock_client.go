package restClient

import (
	"bytes"
	"errors"
	"net/http"
)

// MockHTTPClient structure implementing Client interface
type MockHTTPClient struct {
	ShowError      bool
	SampleResponse interface{}
}

// Post operation to execute HTTP POST operation with support for custom headers
// including setting BasicAuth header information
func (client MockHTTPClient) Post(serviceURL string, header http.Header, payload *bytes.Buffer, userName string, userPass string) (interface{}, error) {

	if client.ShowError {
		return nil, errors.New("external service error")
	}

	return client.SampleResponse, nil
}

// Get operation to execute HTTP GET operation with support for custom headers
// including setting BasicAuth header information
func (client MockHTTPClient) Get(serviceURL string, header http.Header, isTLS bool) (interface{}, error) {

	if client.ShowError {
		return nil, errors.New("external service error")
	}

	return client.SampleResponse, nil
}
