package log

import (
	"reflect"
	"testing"

	"github.com/sirupsen/logrus"
)

func TestNewLoggerWithoutCustomFields(t *testing.T) {

	logger := NewLogger("qqq", "www").WithCustomFields(nil)

	// Validations
	if logger == nil {
		t.Error("Expecting logger object not be nil")
	}

	expectedFields := logrus.Fields{
		"app":     "qqq",
		"version": "www",
	}
	if !reflect.DeepEqual(logger.CurrentEntry(), expectedFields) {
		t.Error("Expecting some of the defined defaults field values to match!")
	}
}

func TestNewLoggerWithCustomFieldsOnly(t *testing.T) {
	logger := NewLogger("ee", "rr").WithCustomFields(map[string]interface{}{"component": "assets"})

	// Validations
	if logger == nil {
		t.Error("Expecting logger object not be nil")
	}

	expectedFields := logrus.Fields{
		"app":       "ee",
		"version":   "rr",
		"component": "assets",
	}
	if !reflect.DeepEqual(logger.CurrentEntry(), expectedFields) {
		t.Error("Expecting some of the defined defaults field values to match!")
	}
}

func TestNewLoggerVerifyVerifyMinimumLogLevel(t *testing.T) {
	logger := NewLogger("", "")

	returnVal := logger.V(3)
	if returnVal {
		t.Error("Not expecting to fail, Incorrect log level")
	}
}
