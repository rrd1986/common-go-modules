package middleware

import (
	"bufio"
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/rrd1986/common-go-modules/log"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func Test_Logging_Middleware(t *testing.T) {
	buff := make([]byte, 0)
	buffer := bytes.NewBuffer(buff)
	templog := logrus.Logger{
		Out:       buffer,
		Formatter: log.DellEmcFormatter{&logrus.JSONFormatter{}},
		Hooks:     make(logrus.LevelHooks),
		Level:     logrus.DebugLevel,
	}
	logger := log.DellEmcLogger{Entry: templog.WithField("app", "app")}
	req, _ := http.NewRequest("GET", "/test", bytes.NewBuffer([]byte("{request:body}")))

	appRouter := setupLoggingMiddleware(t, &logger, successHandler)
	w := httptest.NewRecorder()
	appRouter.ServeHTTP(w, req)

	_, tokens, _ := bufio.ScanLines(buffer.Bytes(), false)
	request := string(tokens)

	assert.Contains(t, request, "request-id")
	assert.Contains(t, request, "request-uri")
	assert.Contains(t, request, "request-method")
	assert.Contains(t, request, "request-header")
	assert.Contains(t, request, "request-body")
	assert.Contains(t, request, "GET")

	_, tokens, _ = bufio.ScanLines(buffer.Bytes()[len(tokens)+1:], false)
	response := string(tokens)

	assert.Contains(t, response, "response-code")
	assert.Contains(t, response, "response-time")
	assert.Contains(t, response, "response-body")
	assert.Contains(t, response, "response-code\":200")

}

func setupLoggingMiddleware(t *testing.T, logger log.LoggerType, handler func(http.ResponseWriter, *http.Request)) *mux.Router {

	appRouter := mux.NewRouter().StrictSlash(true)
	appRouter.Use(LoggingMiddleware(logger))
	appRouter.HandleFunc("/test", handler).Methods("GET")

	return appRouter
}

func successHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Write([]byte("response value"))
}
