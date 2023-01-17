package serviceClient

import (
	"context"
	"crypto/tls"
	"net/url"
	"os"

	"github.com/rrd1986/common-go-modules/utils"

	ngciErrors "github.com/rrd1986/common-go-modules/errors"
	"github.com/rrd1986/common-go-modules/trace"
	resty "gopkg.in/resty.v1"
)

/*
   This client should be used for communication between micro-services
   It allows communication from outside the cluster to within for development work
   It also provides an easily mocked response framework for testing
*/

// Env var
const restyDebugEnvVar = "REQUEST_DEBUG"

// Error Codes
const ErrorCodeResourceNotFound = 101
const ErrorCodeServiceError = 102
const ErrorCodeTypeConversionError = 103
const ErrorCodeInvalidId = 104

type ClientInput struct {
	DefaultProtocol string
	ProtocolEnvVar  string
	DefaultHostName string
	HostNameEnvVar  string
	DefaultPort     string
	PortEnvVar      string
	BaseUrl         string
}

type Type interface {
	GetApiPath() string
	Get(ctx context.Context, params *map[string]string, headers *map[string]string, url string) (HeaderedResponse, error)
	Post(ctx context.Context, body interface{}, headers *map[string]string, url string) (HeaderedResponse, error)
	Patch(ctx context.Context, body interface{}, headers *map[string]string, url string) (HeaderedResponse, error)
	Delete(ctx context.Context, params *map[string]string, headers *map[string]string, url string) (HeaderedResponse, error)
}

type client struct {
	client  *resty.Client
	secure  bool
	ApiPath string
}

func New(ci *ClientInput) Type {
	c := CreateNewRestyClient()
	client := client{client: c}
	if ci != nil {
		client.initClient(ci)
	}
	return &client
}

func CreateNewRestyClient() *resty.Client {
	return resty.New().SetDebug(isDebugMode())
}

func (c *client) GetApiPath() string {
	return c.ApiPath
}

func (c *client) Get(ctx context.Context, params *map[string]string, headers *map[string]string, url string) (HeaderedResponse, error) {
	r := c.getBasicRequest(ctx)
	if headers != nil {
		r.SetHeaders(*headers)
	}
	if params != nil {
		r.SetQueryParams(*params)
	}
	return r.Get(url)
}

func (c *client) Post(ctx context.Context, body interface{}, headers *map[string]string, url string) (HeaderedResponse, error) {
	r := c.getBasicRequest(ctx)
	if headers != nil {
		r.SetHeaders(*headers)
	}
	if body != nil {
		r.SetBody(body)
	}
	return r.Post(url)
}

func (c *client) Patch(ctx context.Context, body interface{}, headers *map[string]string, url string) (HeaderedResponse, error) {
	r := c.getBasicRequest(ctx)
	if headers != nil {
		r.SetHeaders(*headers)
	}
	if body != nil {
		r.SetBody(body)
	}
	return r.Patch(url)
}

func (c *client) Delete(ctx context.Context, params *map[string]string, headers *map[string]string, url string) (HeaderedResponse, error) {
	r := c.getBasicRequest(ctx)
	if headers != nil {
		r.SetHeaders(*headers)
	}
	if params != nil {
		r.SetQueryParams(*params)
	}
	return r.Delete(url)
}

func (c *client) initClient(ci *ClientInput) {
	scheme := ci.DefaultProtocol
	if os.Getenv(ci.ProtocolEnvVar) != "" {
		scheme = os.Getenv(ci.ProtocolEnvVar)
	}

	c.secure = scheme == "https"

	hostName := ci.DefaultHostName
	if os.Getenv(ci.HostNameEnvVar) != "" {
		hostName = os.Getenv(ci.HostNameEnvVar)
	}

	port := ci.DefaultPort
	if os.Getenv(ci.PortEnvVar) != "" {
		port = os.Getenv(ci.PortEnvVar)
	}

	u := url.URL{
		Scheme: scheme,
		Host:   hostName + port,
	}
	c.ApiPath = u.String() + ci.BaseUrl
}

func (c *client) getBasicRequest(ctx context.Context) (r *resty.Request) {
	var req *resty.Request
	if c.secure {
		config := &tls.Config{InsecureSkipVerify: true} // Skip cert validation if hitting the asset service from outside the cluster
		c.client.SetTLSClientConfig(config)             // Only used for development, will not be possible in production TODO: Remove before RTS
		req = c.client.R()
	} else {
		req = c.client.R()
	}
	req.SetHeaders(trace.GetHeaders(ctx))
	return req
}

func GetResourceNotFoundError(resourceName string) error {
	return ngciErrors.NewErrorStr(resourceName+" not found", ErrorCodeResourceNotFound)
}

func GetServiceError(resourceName string) error {
	return ngciErrors.NewErrorStr("error retrieving "+resourceName, ErrorCodeServiceError)
}

func GetTypeConversionError(resourceName string) error {
	return ngciErrors.NewErrorStr("cannot convert "+resourceName, ErrorCodeTypeConversionError)
}

func AppendToUri(uri string, suffixes ...string) string {
	const sep = "/"
	for _, suffix := range suffixes {
		uri = uri + sep + suffix
	}
	return uri
}

func GetAcceptJsonHeader() *map[string]string {
	headers := map[string]string{utils.AcceptType: utils.JSONContentType}
	return &headers
}

func isDebugMode() bool {
	debugMode := os.Getenv(restyDebugEnvVar)
	if debugMode == "TRUE" {
		return true
	}
	return false
}
