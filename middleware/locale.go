package middleware

import (
	"context"
	"net/http"

	"github.com/rrd1986/common-go-modules/log"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

// context key to be used to send locale information
type contextKey string

// ContextLangKey is of type contextKey to save locale key information
const ContextLangKey contextKey = "langKey"

// ContextLocalizerKey is of type contextKey to save localizer instance
const ContextLocalizerKey contextKey = "localizerKey"

// LocaleSelectionMiddleware is a gorilla mux middleware to set language selection for the request
// based on information provided in the request params else the default language will be utilized
func LocaleSelectionMiddleware(i18nBundle *i18n.Bundle, languageMatcher language.Matcher, logger log.LoggerType) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return LocaleSelectionWrapper(next, i18nBundle, languageMatcher, logger)
	}
}

func LocaleSelectionWrapper(next http.Handler, i18nBundle *i18n.Bundle, languageMatcher language.Matcher, logger log.LoggerType) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// retrieve ACCEPT-LANGUAGE header information
		acceptLangHeaderValue := r.Header.Get("Accept-Language")

		languageTag, _ := language.MatchStrings(languageMatcher, acceptLangHeaderValue)
		logger.Info("locale: ", languageTag.String())

		// create new context by adding locale information identified based on provided information
		newCtx := context.WithValue(r.Context(), ContextLangKey, languageTag.String())

		// add localizer to the context
		i18nLocalizer := i18n.NewLocalizer(i18nBundle, acceptLangHeaderValue)
		newCtx = context.WithValue(newCtx, ContextLocalizerKey, i18nLocalizer)

		// call next
		next.ServeHTTP(w, r.WithContext(newCtx))
	})
}
