package log

import (
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

// CustomLogger is a holder for logger
type CustomLogger struct {
	*logrus.Entry
}

type LoggerType interface {
	Tracef(format string, args ...interface{})
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Warningf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
	Panicf(format string, args ...interface{})
	Trace(args ...interface{})
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Warning(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})
	Panic(args ...interface{})
	Traceln(args ...interface{})
	Debugln(args ...interface{})
	Infoln(args ...interface{})
	Warnln(args ...interface{})
	Warningln(args ...interface{})
	Errorln(args ...interface{})
	Fatalln(args ...interface{})
	Panicln(args ...interface{})
	V(level int) bool
	WithCustomFields(fields map[string]interface{}) LoggerType
	CurrentEntry() logrus.Fields
}

func (dl *CustomLogger) CurrentEntry() logrus.Fields {
	return dl.Data
}

// V reports whether verbosity level l is at least the requested verbose level. Refer to grpclog.LoggerV2
func (dl *CustomLogger) V(level int) bool {
	l, ok := map[logrus.Level]int{
		logrus.DebugLevel: 0,
		logrus.InfoLevel:  1,
		logrus.WarnLevel:  2,
		logrus.ErrorLevel: 3,
		logrus.FatalLevel: 4,
	}[dl.Level]

	if !ok {
		return false
	}

	return l >= level
}

// WithCustomFields is a method which adds the ability to add custom fields as additional params to a log entry
func (dl *CustomLogger) WithCustomFields(fields map[string]interface{}) LoggerType {
	return &CustomLogger{dl.WithFields(fields)}
}

// NewLogger initializes a new logger instance
func NewLogger(app string, version string) LoggerType {
	tempLogger := logrus.Logger{
		Out: os.Stdout,
		Formatter: CustomFormatter{&logrus.JSONFormatter{
			TimestampFormat: time.RFC3339Nano,
		}},
		Hooks: make(logrus.LevelHooks),
		Level: logrus.DebugLevel,
	}

	return &CustomLogger{tempLogger.WithField("app", app).WithField("version", version)}
}
