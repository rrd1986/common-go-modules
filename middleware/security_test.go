package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func newHttpHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	})
}
func Test_Security_Wrapper(t *testing.T) {
	respRecorder := httptest.NewRecorder()
	testRequest := &http.Request{Method: "GET", RequestURI: "test.txt"}
	h := newHttpHandler()
	hf := SecurityWrapper(h)
	hf.ServeHTTP(respRecorder, testRequest)
	if respRecorder.Header().Get(PragmaHeader) == "" {
		t.Errorf("Pragma header not set")
	}
	if respRecorder.Header().Get(CacheControlHeader) == "" {
		t.Errorf("Cache Control header not set")
	}
	if respRecorder.Header().Get(XSSProtectionHeader) == "" {
		t.Errorf("XSS Protection header not set")
	}
	if respRecorder.Header().Get(OptionsContentTypeHeader) == "" {
		t.Errorf("X-Content-Type-Options header not set")
	}
	if respRecorder.Header().Get(CSPHeader) == "" {
		t.Errorf("Content-Security-Policy header not set")
	}
	if respRecorder.Header().Get(HSTSHeader) == "" {
		t.Errorf("Strict-Transport-Security header not set")
	}

}
