package log

import (
	"github.com/sirupsen/logrus"
)

// DellEmcFormatter is a custom formatter to add additional params in the log entries
type DellEmcFormatter struct {
	logrus.Formatter
}

// Format log entry to include date time in UTC
func (u DellEmcFormatter) Format(e *logrus.Entry) ([]byte, error) {
	e.Time = e.Time.UTC()
	return u.Formatter.Format(e)
}
