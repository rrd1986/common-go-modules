package trace

import (
	"context"
	"net/http"
	"testing"
)

func TestJSONMarshallSuccess(t *testing.T) {
	var h http.Header
	h = SetTraceHeaders(context.WithValue(context.TODO(), "x-request-id", "22"), h)
	if h == nil {
		t.Errorf("Headers not initialised")
	}
}
