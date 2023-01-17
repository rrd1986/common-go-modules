package middleware

import (
	"net/http"

	"github.com/rrd1986/common-go-modules/log"
	"github.com/rrd1986/common-go-modules/trace"
)

func TracePropagationMiddleware(logger log.LoggerType) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return TracePropagationWrapper(next, logger)
	}
}

func TracePropagationWrapper(next http.Handler, logger log.LoggerType) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		newContext := trace.ContextWithTraceHeaders(r.Context(), r)

		r = r.WithContext(newContext)

		// call next
		next.ServeHTTP(w, r)
	})
}
