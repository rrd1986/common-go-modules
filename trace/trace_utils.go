package trace

import (
	"context"
	"net/http"
	"os"
	"strings"
)

const TracingContextHeaders = "TRACE_HEADERS"

var TraceHeadersToPropagate []string

func init() {
	tracingHeaders := os.Getenv(TracingContextHeaders)
	if tracingHeaders == "" {
		tracingHeaders = "x-request-id,x-b3-traceid,x-b3-spanid,x-b3-parentspanid,x-b3-sampled,x-b3-flags,x-ot-span-context"
	}
	TraceHeadersToPropagate = strings.Split(tracingHeaders, ",")
}

func ContextWithTraceHeaders(ctx context.Context, r *http.Request) context.Context {
	for _, header := range TraceHeadersToPropagate {
		value := r.Header.Get(header)
		if value != "" {
			ctx = context.WithValue(ctx, header, value)
		}
	}
	return ctx
}

func RequestWithTraceHeaders(ctx context.Context, r *http.Request) *http.Request {
	for _, header := range TraceHeadersToPropagate {
		if value, ok := ctx.Value(header).(string); ok {
			r.Header.Add(header, value)
		}
	}
	return r
}

func GetHeaders(ctx context.Context) map[string]string {
	headers := map[string]string{}
	for _, header := range TraceHeadersToPropagate {
		if value, ok := ctx.Value(header).(string); ok {
			headers[header] = value
		}
	}
	return headers
}

func SetTraceHeaders(ctx context.Context, h http.Header) http.Header {
	if h == nil {
		h = make(http.Header)
	}
	for _, header := range TraceHeadersToPropagate {
		if value, ok := ctx.Value(header).(string); ok {
			h.Add(header, value)
		}
	}
	return h
}
