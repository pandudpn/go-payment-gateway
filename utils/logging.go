package utils

import (
	"fmt"
	"os"
	"strings"
	
	"github.com/sirupsen/logrus"
)

const layoutDate = "2006-01-02T15:04:05 -07:00"

type logFormatter struct {
	logrus.TextFormatter
}

func Logger() *logrus.Logger {
	return &logrus.Logger{
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
}

func (l *logFormatter) Format(e *logrus.Entry) ([]byte, error) {
	var (
		levelColor int
	)
	
	switch e.Level {
	case logrus.DebugLevel, logrus.TraceLevel:
		levelColor = 35 // purple
	case logrus.WarnLevel:
		levelColor = 33 // yellow
	case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
		levelColor = 31 // red
	default:
		levelColor = 36 // blue
	}
	
	return []byte(
		fmt.Sprintf(
			"[%s] - \u001B[%dm%s \u001B[0m - %s\n",
			e.Time.Format(l.TimestampFormat),
			levelColor,
			strings.ToUpper(e.Level.String()),
			e.Message,
		),
	), nil
}
