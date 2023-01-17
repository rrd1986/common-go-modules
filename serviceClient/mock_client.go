package serviceClient

import (
	"context"
	"net/http"

	"github.com/stretchr/testify/mock"
)

/**********************************************************************************************************
*
*		Mock client
*
**********************************************************************************************************/

type MockClient struct {
	mock.Mock
}

func (m *MockClient) GetApiPath() string {
	args := m.Called()
	path, _ := args.Get(0).(string)
	return path
}

func (m *MockClient) Get(ctx context.Context, params *map[string]string, headers *map[string]string, url string) (HeaderedResponse, error) {
	args := m.Called(params, headers, url)
	resp, _ := args.Get(0).(HeaderedResponse)
	err, _ := args.Get(1).(error)
	return resp, err
}

func (m *MockClient) Post(ctx context.Context, body interface{}, headers *map[string]string, url string) (HeaderedResponse, error) {
	args := m.Called(body, headers, url)
	resp, _ := args.Get(0).(HeaderedResponse)
	err, _ := args.Get(1).(error)
	return resp, err
}

func (m *MockClient) Patch(ctx context.Context, body interface{}, headers *map[string]string, url string) (HeaderedResponse, error) {
	args := m.Called(body, headers, url)
	resp, _ := args.Get(0).(HeaderedResponse)
	err, _ := args.Get(1).(error)
	return resp, err
}

func (m *MockClient) Delete(ctx context.Context, params *map[string]string, headers *map[string]string, url string) (HeaderedResponse, error) {
	args := m.Called(params, headers, url)
	resp, _ := args.Get(0).(HeaderedResponse)
	err, _ := args.Get(1).(error)
	return resp, err
}

/**********************************************************************************************************
*
*		Mock client response
*
**********************************************************************************************************/

type MockClientResponse struct {
	Code  int
	Bytes []byte
	Head  http.Header
}

func (c MockClientResponse) StatusCode() (code int) {
	return c.Code
}

func (c MockClientResponse) Body() (body []byte) {
	return c.Bytes
}

func (c MockClientResponse) Header() (header http.Header) {
	return c.Head
}
