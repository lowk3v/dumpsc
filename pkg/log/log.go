package log

import (
	"github.com/sirupsen/logrus"
)

type Logger struct {
	*logrus.Logger
}

func New(level string) *Logger {
	lvl, err := logrus.ParseLevel(level)
	if err != nil {
		lvl = logrus.InfoLevel
	}

	l := &Logger{
		Logger: logrus.New(),
	}
	l.SetReportCaller(false)
	l.SetLevel(lvl)
	l.Formatter = &logrus.TextFormatter{
		DisableColors:   false,
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	}
	return l
}
