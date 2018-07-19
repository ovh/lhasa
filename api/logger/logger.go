package logger

import (
	"os"

	"github.com/fabienm/go-logrus-formatters"
	"github.com/sirupsen/logrus"
)

// NewLogger creates a configured logrus logger instance
func NewLogger(isVerbose, isDebug, isQuiet, isJSON bool) logrus.FieldLogger {
	log := logrus.New()
	log.Level = logrus.WarnLevel
	if isVerbose {
		log.Level = logrus.InfoLevel
	}
	if isQuiet {
		log.Level = logrus.FatalLevel
	}
	if isDebug {
		log.Level = logrus.DebugLevel
	}
	if isJSON {
		hostname, _ := os.Hostname()
		log.Formatter = formatters.NewGelf(hostname)
	}

	return log
}
