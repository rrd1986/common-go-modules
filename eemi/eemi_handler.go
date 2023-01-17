package eemi

import (
	"context"
	"net/http"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/rrd1986/common-go-modules/log"
	"github.com/rrd1986/common-go-modules/middleware"
	"github.com/rrd1986/common-go-modules/utils"
)

type Handler struct {
	handle   func(ee http.ResponseWriter, rr *http.Request) error
	eemiData map[string]Config
}

var logger log.LoggerType

func SetLogger(l log.LoggerType) {
	logger = l
}

func NewHandler(handle func(w http.ResponseWriter, r *http.Request) error, eemiData map[string]Config) Handler {
	return Handler{handle: handle, eemiData: eemiData}
}

// Wrap http handler that allows handler functions to return an error
// This centralises the logic for creating error responses
// Loads eemi info based on the supplied map. Localises message and responseaction values
func (fn Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := fn.handle(w, r); err != nil {
		eemiError, ok := err.(Error)
		if !ok {
			eemiError = New(err, "UNHANDLED") // If no eemi code was detected, return with a predefined `unhandled` generic error message
		}
		logger.Error(eemiError)

		// get data for eemi error based on messageId in eemi error
		eemiData := fn.eemiData[eemiError.MessageID]

		// translate textual parts of response
		message := translateMessageByID(r.Context(), eemiData.Message, eemiError.TemplateData)
		responseAction := translateMessageByID(r.Context(), eemiData.ResponseAction, eemiError.TemplateData)

		eemiResponse := NewEemiResponse(eemiData.MessageID, message, responseAction, eemiData.Category, eemiData.SeverityLevel)
		utils.WriteResponse(w, eemiData.Status, eemiResponse, utils.JSONContentType)
	}
}

// translateMessageByID utility function to translate message based on current locale information
func translateMessageByID(ctx context.Context, messageID string, messageArgs map[string]string) string {
	// retrieve localizer instance from context
	localizer := ctx.Value(middleware.ContextLocalizerKey).(*i18n.Localizer)

	// construct message
	message := localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID: messageID,
		},
		TemplateData: messageArgs,
	})

	return message
}
