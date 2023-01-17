package eemi

import (
	"fmt"
)

// NGCIError defines exception details confirming to EEMI (Enhanced Error Message Initiative) standard.
type Error struct {
	cause        error
	MessageID    string
	TemplateData map[string]string
}

// implement default error interface
func (e Error) Error() string {
	return fmt.Sprintf("id: %s, cause: %s", e.MessageID, e.cause)
}

// FormatError method to format the service response for any exception occurred
func (e *Error) FormatError() map[string]interface{} {
	return map[string]interface{}{
		"messageID": e.MessageID,
	}
}

// CreateEEMIError create a new instance of Error
func New(err error, messageID string) Error {
	return Error{cause: err, MessageID: messageID, TemplateData: nil}
}

// To be used for EEMI messages where the response needs uses templates in the message
func NewWithTemplateData(err error, messageID string, templateData map[string]string) Error {
	return Error{cause: err, MessageID: messageID, TemplateData: templateData}
}
