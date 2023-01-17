package utils

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/rrd1986/common-go-modules/middleware"
)

type AssetSample struct {
	Type        string `json:"assetType"`
	Description string `json:"assetDesc"`
}

func TestWriteResponseMissingHeaderExpectJSONResponse(t *testing.T) {

	respRecorder := httptest.NewRecorder()

	payload := &AssetSample{
		Type:        "mgmtswitch",
		Description: "sample asset description",
	}

	// sending empty map as DB record
	WriteResponse(respRecorder, http.StatusCreated, payload, "")

	expected := `{"assetType":"mgmtswitch","assetDesc":"sample asset description"}`

	if respRecorder.Header().Get(AcceptType) != "" {
		t.Errorf("Not expecting accept type header to be set!")
	}

	if strings.Compare(respRecorder.Body.String(), expected) != 0 {
		// TODO Revsit t.Errorf(responseMismatchPrintFormat)
	}
}

func TestWriteResponseJSONTypeReturnsJSONResponse(t *testing.T) {

	respRecorder := httptest.NewRecorder()

	payload := &AssetSample{
		Type:        "mgmtswitch",
		Description: "sample asset description",
	}

	// sending empty map as DB record
	WriteResponse(respRecorder, http.StatusCreated, payload, JSONContentType)

	if respRecorder.Header().Get(AcceptType) != JSONContentType {
		t.Errorf("Not expecting accept type header to be set!")
	}

	actual := respRecorder.Body.String()
	expected := `{"assetType":"mgmtswitch","assetDesc":"sample asset description"}`
	if strings.Compare(actual, expected) != 0 {
		t.Errorf("Expecting response body message in following format [" + expected + "] but received" + actual)
	}
}

func TestWriteResponseEmptyJSONTypeReturnsEmptyObject(t *testing.T) {

	respRecorder := httptest.NewRecorder()

	// sending empty map as DB record
	WriteResponse(respRecorder, http.StatusCreated, nil, JSONContentType)

	if respRecorder.Header().Get(AcceptType) != JSONContentType {
		t.Errorf("Not expecting accept type header to be set!")
	}

	actual := respRecorder.Body.String()
	expected := `{}`
	if strings.Compare(actual, expected) != 0 {
		t.Errorf("Expecting response body message in following format [" + expected + "] but received" + actual)
	}
}

func TestWriteResponseXMLTypeThrowsError(t *testing.T) {

	respRecorder := httptest.NewRecorder()

	payload := map[string]interface{}{
		"Feb": map[string]interface{}{"Fullform": "February"},
		"Jan": map[string]interface{}{"Fullform": "January"},
		"Mar": map[string]interface{}{"Fullform": "March"},
	}

	// sending empty map as DB record
	WriteResponse(respRecorder, http.StatusCreated, payload, XMLContentType)

	if respRecorder.Header().Get(AcceptType) != XMLContentType {
		t.Errorf("Not expecting accept type header to be set!")
	}

	actual := respRecorder.Body.String()
	expected := "XML is not a supported format at the moment, please change the header value to application/json"
	if strings.Compare(actual, expected) != 0 {
		t.Errorf("Expecting an errror message in following format [" + expected + "] but received [" + actual + "]")
	}
}

func newHttpHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	})
}

func executeDummyRequest(uri string, swaggerPath string) *httptest.ResponseRecorder {
	respRecorder := httptest.NewRecorder()
	testRequest := &http.Request{Method: "GET", RequestURI: uri}
	h := newHttpHandler()
	hf := SetSecurityHeaders(h, swaggerPath)
	hf.ServeHTTP(respRecorder, testRequest)
	return respRecorder
}

func validateSecurityHeaders(t *testing.T, uri string, swaggerPath string) {
	respRecorder := executeDummyRequest(uri, swaggerPath)

	if respRecorder.Header().Get(middleware.CSPHeader) == "" {
		t.Errorf("Expecting CSP type header to be set!")
	}

	if respRecorder.Header().Get(HSTSHeader) == "" {
		t.Errorf("Expecting CSP type header to be set!")
	}
}

func TestSecurityHeadersAddedForHtmlFiles(t *testing.T) {
	validateSecurityHeaders(t, "test.html", "swaggerPath")
}

func TestSecurityHeadersAddedForHtmFiles(t *testing.T) {
	validateSecurityHeaders(t, "test.htm", "swaggerPath")
}

func TestSecurityHeadersAddedForSwaggerFiles(t *testing.T) {
	validateSecurityHeaders(t, "swaggerPath", "swaggerPath")
}

func TestSecurityHeadersNotAddedForNonHtmlFiles(t *testing.T) {
	respRecorder := executeDummyRequest("test.txt", "swaggerPath")

	if respRecorder.Header().Get(middleware.CSPHeader) != "" {
		t.Errorf("Expecting CSP type header not to be set!")
	}

	if respRecorder.Header().Get(HSTSHeader) != "" {
		t.Errorf("Expecting CSP type header not to be set!")
	}

}
