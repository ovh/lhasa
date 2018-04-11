package logger

import (
	"context"
	"crypto/tls"
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/fabienm/go-logrus-formatters"
	loghook "github.com/ovh/cds/sdk/log/hook"
	"github.com/sirupsen/logrus"
)

var (
	hook *loghook.Hook
)

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
		hostname, _ := os.Hostname()
		log.Formatter = formatters.NewGelf(hostname)
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
