package db

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

func getLogger(log logrus.FieldLogger) gorm.Logger {
	return gorm.Logger{LogWriter: dbLogWriter{log}}
}

type dbLogWriter struct {
	log logrus.FieldLogger
}

// Println implements gorm's LogWriter interface
func (l dbLogWriter) Println(v ...interface{}) {
	if l.log == nil {
		return
	}
	if len(v) != 5 {
		return
	}

	var message interface{} = v
	file := strings.Split(strings.TrimSuffix(strings.TrimPrefix(v[0].(string), "\u001b[35m("), ")\u001b[0m"), ":")
	durationStr := strings.TrimSuffix(strings.TrimPrefix(v[2].(string), " \u001b[36;1m["), "]\u001b[0m ")
	rows, _ := strconv.Atoi(strings.TrimSuffix(strings.TrimPrefix(v[4].(string), " \n\u001b[36;31m["), " rows affected or returned ]\u001b[0m "))
	duration, _ := time.ParseDuration(durationStr)
	line, _ := strconv.Atoi(file[1])
	fields := logrus.Fields{
		"file":         file[0],
		"line":         line,
		"duration":     duration.Seconds(),
		"full_message": v[3],
		"rows":         rows,
	}
	switch v[3].(type) {
	case string:
		message = fmt.Sprintf("sql query: %s...", v[3].(string)[:int(math.Min(50, float64(len(v[3].(string)))))])
		l.log.WithFields(fields).Debug(message)
	case error:
		l.log.WithFields(fields).WithError(v[3].(error)).Error(message)
	default:
		l.log.WithFields(fields).Debug(message)
	}
}
