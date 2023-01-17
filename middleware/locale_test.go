package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rrd1986/common-go-modules/log"

	"github.com/gorilla/mux"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

func verifyMiddlewareFnCall(t *testing.T, expectedLocale string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// verify locale
		actualLocale := r.Context().Value(ContextLangKey).(string)
		if actualLocale != expectedLocale {
			t.Errorf("Expecting locale to be %s but got %s", expectedLocale, actualLocale)
		}

		// ensure localizer instance is available
		localizer := r.Context().Value(ContextLocalizerKey).(*i18n.Localizer)
		if localizer == nil {
			t.Error("Expecting i18n localizer to be defined!")
		}
	})
}

func setupMiddleware(t *testing.T, locale string) *mux.Router {

	i18nBundle := i18n.NewBundle(language.English)

	languageMatcher := language.NewMatcher([]language.Tag{
		language.English, language.German,
	})

	appRouter := mux.NewRouter().StrictSlash(true)
	appRouter.Use(LocaleSelectionMiddleware(i18nBundle, languageMatcher, log.NewLogger("", "")))
	appRouter.Handle("/", verifyMiddlewareFnCall(t, locale)).Methods("GET")

	return appRouter
}

func TestLangSelectionMiddlewareWithAcceptLanguageHeaderDefinedExpectsDELocale(t *testing.T) {

	locale := "de"
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Add("Accept-Language", locale)

	appRouter := setupMiddleware(t, locale)
	w := httptest.NewRecorder()
	appRouter.ServeHTTP(w, req)
}

func TestLangSelectionMiddlewareWithAcceptLanguageHeaderMissingExpectsENLocale(t *testing.T) {

	locale := "en"

	req, _ := http.NewRequest("GET", "/", nil)

	appRouter := setupMiddleware(t, locale)
	w := httptest.NewRecorder()
	appRouter.ServeHTTP(w, req)
}

func TestLangSelectionMiddlewareWithAcceptLanguageHeaderUnsupportedExpectsENLocale(t *testing.T) {

	locale := "es"

	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Add("Accept-Language", locale)

	appRouter := setupMiddleware(t, "en")
	w := httptest.NewRecorder()
	appRouter.ServeHTTP(w, req)
}
