package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rrd1986/common-go-modules/log"

	"github.com/gorilla/mux"
)

func createSuccessHandler(t *testing.T) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("response value"))

		headerVal := r.Header.Get("x-b3-traceid")
		contextVal := r.Context().Value("x-b3-traceid")
		if headerVal != contextVal {
			t.Error("Expecting i18n localizer to be defined!")
		}
	})
}

func Test_Trace_Propagation_Propagates_Header(t *testing.T) {
	appRouter := mux.NewRouter().StrictSlash(true)
	appRouter.Use(TracePropagationMiddleware(log.NewLogger("", "")))
	appRouter.Handle("/", createSuccessHandler(t)).Methods("GET")

	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Add("x-b3-traceid", "123")

	w := httptest.NewRecorder()
	appRouter.ServeHTTP(w, req)
}
