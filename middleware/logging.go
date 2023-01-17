package middleware

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/rrd1986/common-go-modules/log"
)

const RedactedLogMessage = "-- Removed from logging --"

func LoggingMiddleware(logger log.LoggerType) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return LoggingWrapper(next, logger, nil, nil)
	}
}

func LoggingWrapper(next http.Handler, logger log.LoggerType, filters *LogFilters, extraParams map[string]LoggingParam) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if filters == nil {
			filters = &LogFilters{Header: HeaderFilterDefault, Body: BodyFilterDefault}
		}

		requestId, _ := uuid.NewRandom()

		// Read body bytes stream and then reset value
		bodyBytes, _ := ioutil.ReadAll(r.Body)
		r.Body.Close()
		r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

		// log custom fields
		logger = logger.WithCustomFields(map[string]interface{}{
			"request-id":     requestId,
			"request-uri":    r.RequestURI,
			"request-method": r.Method,
			"request-header": filters.Header(r.RequestURI, r.Header),
			"request-body":   filters.Body(r.RequestURI, bodyBytes),
		})

		if extraParams != nil {
			for k := range extraParams {
				loggingParam := extraParams[k]
				param := loggingParam(r)
				if param != nil {
					logger = logger.WithCustomFields(map[string]interface{}{
						k: param,
					})
				}
			}
		}

		// create logging responseWriter to capture response body
		loggingRW := &loggingResponseWriter{
			ResponseWriter: w,
		}

		logger.Infof("Request to %s endpoint", r.RequestURI)

		start := makeTimestamp()

		// call next
		next.ServeHTTP(loggingRW, r)

		finish := makeTimestamp()

		logger = logger.WithCustomFields(map[string]interface{}{
			"response-code": loggingRW.status,
			"response-time": finish - start,
			"response-body": filters.Body(r.RequestURI, loggingRW.body),
		})

		logger.Infof("Response from %s endpoint", r.RequestURI)
	})
}

func makeTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

type loggingResponseWriter struct {
	status int
	body   []byte
	http.ResponseWriter
}

func (w *loggingResponseWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *loggingResponseWriter) Write(body []byte) (int, error) {
	w.body = body
	return w.ResponseWriter.Write(body)
}

type LogFilters struct {
	Header HeaderFilter
	Body   BodyFilter
}

type HeaderFilter func(uri string, header http.Header) string
type BodyFilter func(uri string, body []byte) string

func HeaderFilterDefault(uri string, header http.Header) string {
	// Filter swagger by default
	if strings.Contains(uri, "swagger") {
		return RedactedLogMessage
	} else {
		return fmt.Sprintf("%s", header)
	}
}

func BodyFilterDefault(uri string, body []byte) string {
	// Filter swagger by default
	if strings.Contains(uri, "swagger") {
		return RedactedLogMessage
	} else {
		return fmt.Sprintf("%s", body)
	}
}

type LoggingParam func(request *http.Request) *string
