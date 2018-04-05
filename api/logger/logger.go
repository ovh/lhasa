package logger

import (
	"context"
	"crypto/tls"
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/gin-gonic/gin/json"
	loghook "github.com/ovh/cds/sdk/log/hook"
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
	hook            *loghook.Hook
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

// Hook fix hook
func Hook(log *logrus.Logger, conf *Conf) *logrus.Logger {
	Initialize(log, conf)
	return log
}

// Conf contains log configuration
type Conf struct {
	Level             string
	GraylogHost       string
	GraylogPort       string
	GraylogProtocol   string
	GraylogExtraKey   string
	GraylogExtraValue string
	Ctx               context.Context
}

// Initialize graylog hook if any
func Initialize(mylog *logrus.Logger, conf *Conf) {
	if conf.GraylogHost != "" && conf.GraylogPort != "" {

		graylogcfg := &loghook.Config{
			Addr:      fmt.Sprintf("%s:%s", conf.GraylogHost, conf.GraylogPort),
			Protocol:  conf.GraylogProtocol,
			TLSConfig: &tls.Config{ServerName: conf.GraylogHost},
		}

		extra := map[string]interface{}{}
		if conf.GraylogExtraKey != "" && conf.GraylogExtraValue != "" {
			keys := strings.Split(conf.GraylogExtraKey, ",")
			values := strings.Split(conf.GraylogExtraValue, ",")
			if len(keys) != len(values) {
				mylog.Errorf("Error while initialize log: extraKey (len:%d) does not have same corresponding number of values on extraValue (len:%d)", len(keys), len(values))
			} else {
				for i := range keys {
					extra[keys[i]] = values[i]
				}
			}
		}

		extra["OS"] = runtime.GOOS
		extra["Arch"] = runtime.GOARCH

		// no need to check error here
		hostname, _ := os.Hostname()
		extra["Hostname"] = hostname

		var errhook error
		hook, errhook = loghook.NewHook(graylogcfg, extra)

		if errhook != nil {
			mylog.Errorf("Error while initialize graylog hook: %v", errhook)
		} else {
			mylog.AddHook(hook)
		}
	}

	if conf.Ctx == nil {
		conf.Ctx = context.Background()
	}
	go func() {
		<-conf.Ctx.Done()
		mylog.Infof("Draining logs")
		if hook != nil {
			hook.Flush()
		}
	}()
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
