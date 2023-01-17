package log

import (
	"os"
	"testing"

	"github.com/sirupsen/logrus"
)

func TestFormat(t *testing.T) {

	formatter := CustomFormatter{&logrus.JSONFormatter{}}

	logger := logrus.Logger{
		Out:       os.Stdout,
		Formatter: CustomFormatter{&logrus.JSONFormatter{}},
		Hooks:     make(logrus.LevelHooks),
		Level:     logrus.DebugLevel,
	}
	response, error := formatter.Format(logger.WithField("", ""))

	if error != nil {
		t.Error("Not expecting exception here!")
	}
	if response == nil {
		t.Error("Expecting formatted response here!")
	}
}
