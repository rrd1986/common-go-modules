package log

import (
	"github.com/sirupsen/logrus"
)

// CustomFormatter is a custom formatter to add additional params in the log entries
type CustomFormatter struct {
	logrus.Formatter
}

// Format log entry to include date time in UTC
func (u CustomFormatter) Format(e *logrus.Entry) ([]byte, error) {
	e.Time = e.Time.UTC()
	return u.Formatter.Format(e)
}
