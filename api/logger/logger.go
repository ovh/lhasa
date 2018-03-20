package logger

import (
	"os"

	"github.com/gin-gonic/gin/json"
	"github.com/sirupsen/logrus"
)

type syslogLevel uint

// Syslog levels
const (
	SyslogEmergency     syslogLevel = 0
	SyslogAlert         syslogLevel = 1
	SyslogCritical      syslogLevel = 2
	SyslogError         syslogLevel = 3
	SyslogWarning       syslogLevel = 4
	SyslogNotice        syslogLevel = 5
	SyslogInformational syslogLevel = 6
	SyslogDebugging     syslogLevel = 7
)

const (
	gelfVersion        = "1.1"
	defaultSyslogLevel = SyslogInformational
)

var (
	levelMap        map[logrus.Level]syslogLevel
	protectedFields map[string]bool
)

func init() {
	levelMap = map[logrus.Level]syslogLevel{
		logrus.DebugLevel: SyslogDebugging,
		logrus.InfoLevel:  SyslogInformational,
		logrus.WarnLevel:  SyslogWarning,
		logrus.ErrorLevel: SyslogError,
		logrus.FatalLevel: SyslogCritical,
		logrus.PanicLevel: SyslogEmergency,
	}
	protectedFields = map[string]bool{
		"version":       true,
		"host":          true,
		"short_message": true,
		"full_message":  true,
		"timestamp":     true,
		"level":         true,
	}
}

// NewLogger creates a configured logrus logger instance
func NewLogger(isVerbose, isDebug, isQuiet, isJSON bool) *logrus.Logger {
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
		log.Formatter = &gelfFormatter{}
	}
	return log
}

type gelfFormatter struct {
}

// Format implements logrus formatter
func (*gelfFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	host, _ := os.Hostname()
	nanosecond := float64(entry.Time.Nanosecond()) / 1e9
	gelfEntry := map[string]interface{}{
		"version":       gelfVersion,
		"short_message": entry.Message,
		"level":         toSyslogLevel(entry.Level),
		"timestamp":     float64(entry.Time.Unix()) + nanosecond,
		"host":          host,
	}
	for key, value := range entry.Data {
		if !protectedFields[key] {
			key = "_" + key
		}
		gelfEntry[key] = value
	}

	message, err := json.Marshal(gelfEntry)
	return append(message, '\n'), err
}

// Levels implements logrus hook
func toSyslogLevel(level logrus.Level) syslogLevel {
	syslog, ok := levelMap[level]
	if ok {
		return syslog
	}
	return defaultSyslogLevel
}
