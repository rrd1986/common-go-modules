package utils

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/rrd1986/common-go-modules/log"
	"github.com/rrd1986/common-go-modules/middleware"
)

// ContentType represents HTTP Conent-Type header
const ContentType = "Content-Type"

// AcceptType represents HTTP Accept header param
const AcceptType = "accept"

// JSONContentType is to represent JSON header param value type
const JSONContentType = "application/json"

// XMLContentType is to represent XML header param value type
const XMLContentType = "application/xml"

// PlainTextContentType is to represent plain text param value type
const PlainTextContentType = "text/plain"

// PDFContentType is to represent pdf header param value type
const PDFContentType = "application/pdf"

// hstsHeader represents the HSTS Header param
const HSTSHeader = "Strict-Transport-Security"

var logger log.LoggerType

func SetLogger(l log.LoggerType) {
	logger = l
}

// WriteResponse is an utility to format the service response based on the Accept HTTP header information
func WriteResponse(w http.ResponseWriter, status int, responseData interface{}, contentType string) {

	if contentType != "" {
		w.Header().Set(AcceptType, contentType)
		w.Header().Set(ContentType, contentType)
	}

	var dataToWrite []byte
	var writeError error
	switch contentType {
	case JSONContentType:
		dataToWrite, writeError = JSONMarshall(responseData)
	case XMLContentType:
		writeError = errors.New("XML is not a supported format at the moment, please change the header value to application/json")
	case PlainTextContentType, PDFContentType:
		str, ok := responseData.(string)
		if !ok {
			writeError = errors.New("error converting data to plain text")
		} else {
			dataToWrite = []byte(str)
		}
	default:
		dataToWrite, writeError = JSONMarshall(responseData)
	}

	if writeError != nil {
		w.WriteHeader(http.StatusInternalServerError)
		updateResponseWriter(w, []byte(writeError.Error()))
	}

	// write response in JSON format
	w.WriteHeader(status)
	updateResponseWriter(w, dataToWrite)
}

func updateResponseWriter(w http.ResponseWriter, data []byte) {
	// write response data
	contentLength, writeErr := w.Write(data)
	// log exception
	if writeErr != nil {
		logger.Error("Failed to write service response data, reason: ", writeErr)
	}
	// set content length
	w.Header().Set("Content-Length", strconv.Itoa(contentLength))
}

func SetSecurityHeaders(h http.Handler, swaggerPath string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if (r.RequestURI == swaggerPath) || (strings.HasSuffix(r.RequestURI, ".html")) || (strings.HasSuffix(r.RequestURI, ".htm")) {
			w.Header().Set(HSTSHeader, "max-age=31536000")
			w.Header().Set(middleware.CSPHeader, "frame-ancestors 'self'")
		}
		h.ServeHTTP(w, r)
	})
}
