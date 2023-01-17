package restClient

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/rrd1986/common-go-modules/trace"
)

// Client defines operations supported for invoking external services
type Client interface {
	Post(ctx context.Context, serviceURL string, Header http.Header, payload *bytes.Buffer, userName string, userPass string) (interface{}, error)
	Get(ctx context.Context, serviceURL string, header http.Header, isTLS bool) (interface{}, error)
}

// HTTPClient structure implementing Client interface
type HTTPClient struct {
}

// Post operation to execute HTTP POST operation with support for custom headers
// including setting BasicAuth header information
func (client HTTPClient) Post(ctx context.Context, serviceURL string, header http.Header, payload *bytes.Buffer, userName string, userPass string) (interface{}, error) {

	// create new request object
	request, reqErr := http.NewRequest("POST", serviceURL, payload)
	if reqErr != nil {
		return nil, reqErr
	}

	// set header parameters
	request.Header = header
	request = trace.RequestWithTraceHeaders(ctx, request)

	// set basic auth information, if data available
	request.SetBasicAuth(userName, userPass)

	// create http client
	httpClient := &http.Client{Timeout: time.Second * 5}

	// invoke service call
	serviceResponse, serviceErr := httpClient.Do(request)
	// throw exception if service fails
	if serviceErr != nil {
		return nil, serviceErr
	}

	//defer serviceResponse.Body.Close()

	if serviceResponse.StatusCode != http.StatusOK || serviceResponse.ContentLength == 0 {
		errorMessage := "Expected valid response from the invoked service, Please re-try again. status received: " +
			strconv.Itoa(serviceResponse.StatusCode)
		return nil, errors.New(errorMessage)
	}

	// gather results from service invoked and return
	body, readError := ioutil.ReadAll(serviceResponse.Body)
	if readError != nil {
		return nil, readError
	}

	var responseData interface{}
	jsonError := json.Unmarshal(body, &responseData)
	if jsonError != nil {
		return nil, jsonError
	}

	return responseData, nil
}

// Get operation to execute HTTP GET operation with support for custom headers and optional TLS flag
func (client HTTPClient) Get(ctx context.Context, serviceURL string, header http.Header, isTLS bool) (interface{}, error) {

	// create new request object
	request, reqErr := http.NewRequest("GET", serviceURL, nil)
	if reqErr != nil {
		return nil, reqErr
	}

	// set header parameters
	request.Header = header
	request = trace.RequestWithTraceHeaders(ctx, request)

	// create http client
	httpClient := &http.Client{Timeout: time.Second * 5}

	// TODO: Need to set the mTLS based on JWT TOKEN, Currently turned off
	if isTLS {
		transport := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		httpClient.Transport = transport
	}

	// invoke service call
	serviceResponse, serviceErr := httpClient.Do(request)

	// throw exception if service fails
	if serviceErr != nil {
		return nil, serviceErr
	}

	if serviceResponse.StatusCode != http.StatusOK || serviceResponse.ContentLength == 0 {
		errorMessage := "Expected valid response from the invoked service, Please re-try again. status received: " +
			strconv.Itoa(serviceResponse.StatusCode)
		return nil, errors.New(errorMessage)
	}

	// gather results from service invoked and return
	body, readError := ioutil.ReadAll(serviceResponse.Body)
	if readError != nil {
		return nil, readError
	}

	var responseData interface{}
	jsonError := json.Unmarshal(body, &responseData)
	if jsonError != nil {
		return nil, jsonError
	}

	return responseData, nil
}
