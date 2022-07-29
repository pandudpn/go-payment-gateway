package pg

import (
	"fmt"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

const layoutDate = "2006-01-02T15:04:05-07:00"

type logger struct {
	log bool
	*logrus.Logger
}

// Log global variable to use as a logging interface
var Log Logging

// Logging define all method to be implemented as logging
// only support print, warn and error
type Logging interface {
	// Print write debug message without format into stdout
	Print(args ...interface{})

	// Printf write debug message with format into stdout
	Printf(format string, args ...interface{})

	// Println write debug message with new line into stdout
	Println(args ...interface{})

	// Warn write warn message without format into stdout
	Warn(args ...interface{})

	// Warnf write warn message with format into stdout
	Warnf(format string, args ...interface{})

	// Error write error message without format into stdout
	Error(args ...interface{})

	// Errorf print error message with format
	Errorf(format string, args ...interface{})
}

func (l *logger) Print(args ...interface{}) {
	if !l.log {
		return
	}

	l.Logger.Print(args...)
}

func (l *logger) Printf(format string, args ...interface{}) {
	if !l.log {
		return
	}

	l.Logger.Printf(format, args...)
}

func (l *logger) Println(args ...interface{}) {
	if !l.log {
		return
	}

	l.Logger.Println(args...)
}

func (l *logger) Warn(args ...interface{}) {
	if !l.log {
		return
	}

	l.Logger.Warn(args...)
}

func (l *logger) Warnf(format string, args ...interface{}) {
	if !l.log {
		return
	}

	l.Logger.Warnf(format, args...)
}

func (l *logger) Error(args ...interface{}) {
	if !l.log {
		return
	}

	l.Logger.Error(args...)
}

func (l *logger) Errorf(format string, args ...interface{}) {
	if !l.log {
		return
	}

	l.Logger.Errorf(format, args...)
}

type logFormatter struct {
	logrus.TextFormatter
}

// setLogger create a new logrus.Logger
func setLogger() *logrus.Logger {
	l := &logrus.Logger{
		Out:          os.Stderr,
		ReportCaller: true,
		Level:        logrus.DebugLevel,
		Formatter: &logFormatter{
			logrus.TextFormatter{
				ForceColors:            true,
				FullTimestamp:          true,
				TimestampFormat:        layoutDate,
				DisableLevelTruncation: false,
			},
		},
	}
	// set logrus into Log variable
	Log = &logger{true, l}

	return l
}

// NewLogger create instance of Logging
func NewLogger() *logrus.Logger {
	return setLogger()
}

// DisableLogging set configuration for disable logging
func DisableLogging() {
	Log = &logger{log: false}
}

// Format create formatted logging
func (l *logFormatter) Format(e *logrus.Entry) ([]byte, error) {
	var (
		levelColor int
	)

	switch e.Level {
	case logrus.WarnLevel:
		levelColor = 33 // yellow
	case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
		levelColor = 31 // red
	default:
		levelColor = 34 // blue
	}

	return []byte(
		fmt.Sprintf(
			"time=%s level=\u001B[%dm%s\u001B[0m %s\n",
			e.Time.Format(l.TimestampFormat),
			levelColor,
			strings.ToUpper(e.Level.String()),
			e.Message,
		),
	), nil
}
