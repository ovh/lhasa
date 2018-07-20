package handlers

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ovh/lhasa/api/hateoas"
	"github.com/ovh/lhasa/api/security"
	"github.com/sirupsen/logrus"
)

// Gin context keys
const (
	KeyRoles    = "roles"
	KeyLogEntry = "logger"
)

var (
	dunno     = []byte("???")
	centerDot = []byte("·")
	dot       = []byte(".")
	slash     = []byte("/")
)

// RecoveryWithLogger returns a middleware for a given logger that recovers from any panics and writes a 500 if there was one.
func RecoveryWithLogger(log logrus.FieldLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				if log != nil {
					stack := stack(3)
					httprequest, _ := httputil.DumpRequest(c.Request, false)
					log.WithField("error", err).
						WithField("full_message", stack).
						WithField("request", string(httprequest)).
						Error("panic recovered")
				}
				c.AbortWithStatus(500)
			}
		}()
		c.Next()
	}
}

// stack returns a nicely formatted stack frame, skipping skip frames.
func stack(skip int) []byte {
	buf := new(bytes.Buffer) // the returned data
	// As we loop, we open files and read them. These variables record the currently
	// loaded file.
	var lines [][]byte
	var lastFile string
	for i := skip; ; i++ { // Skip the expected number of frames
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		// Print this much at least.  If we can't find the source, it won't show.
		fmt.Fprintf(buf, "%s:%d (0x%x)\n", file, line, pc)
		if file != lastFile {
			data, err := ioutil.ReadFile(file)
			if err != nil {
				continue
			}
			lines = bytes.Split(data, []byte{'\n'})
			lastFile = file
		}
		fmt.Fprintf(buf, "\t%s: %s\n", function(pc), source(lines, line))
	}
	return buf.Bytes()
}

// source returns a space-trimmed slice of the n'th line.
func source(lines [][]byte, n int) []byte {
	n-- // in stack trace, lines are 1-indexed but our array is 0-indexed
	if n < 0 || n >= len(lines) {
		return dunno
	}
	return bytes.TrimSpace(lines[n])
}

// function returns, if possible, the name of the function containing the PC.
func function(pc uintptr) []byte {
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return dunno
	}
	name := []byte(fn.Name())
	// The name includes the path name to the package, which is unnecessary
	// since the file name is already included.  Plus, it has center dots.
	// That is, we see
	//	runtime/debug.*T·ptrmethod
	// and want
	//	*T.ptrmethod
	// Also the package path might contains dot (e.g. code.google.com/...),
	// so first eliminate the path prefix
	if lastslash := bytes.LastIndex(name, slash); lastslash >= 0 {
		name = name[lastslash+1:]
	}
	if period := bytes.Index(name, dot); period >= 0 {
		name = name[period+1:]
	}
	name = bytes.Replace(name, centerDot, dot, -1)
	return name
}

// LoggingMiddleware logs before and after incoming gin requests
func LoggingMiddleware(logHeaders []string, log logrus.FieldLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		if log == nil {
			c.Next()
			return
		}
		entry := log.WithFields(logrus.Fields{
			"method": c.Request.Method,
			"path":   c.Request.URL.Path,
		})
		for _, h := range logHeaders {
			header := http.CanonicalHeaderKey(h)
			if header != "" {
				entry = entry.WithField(header, c.GetHeader(header))
			}
		}
		c.Set(KeyLogEntry, entry)

		entry.Debug("incoming request")
		startTime := time.Now()

		c.Next()

		entry.WithFields(logrus.Fields{
			"status":   c.Writer.Status(),
			"duration": time.Since(startTime).Seconds(),
		}).Debug("done")

		for _, err := range c.Errors.Errors() {
			if err != hateoas.ErrorCreated.Error() {
				entry.Error(err)
			}
		}
	}
}

// HasOne returns a gin handler that checks the request against the given role
func HasOne(roles ...security.Role) gin.HandlerFunc {
	return func(c *gin.Context) {
		rolesRaw, ok := c.Get(KeyRoles)
		rolePolicy, castok := rolesRaw.(security.RolePolicy)
		log := GetLogger(c)

		if !ok || !castok {
			if log != nil {
				log.Warnf("unusable security policy, please check your configuration")
				abortUnauthorized(c)
				return
			}
		}

		if rolePolicy.HasOne(roles...) {
			c.Next()
			return
		}
		if log != nil {
			log.WithField("expectedRoles", roles).WithField("roles", rolePolicy.ToSlice()).Info("unauthorized access")
		}
		abortUnauthorized(c)
	}
}

func abortUnauthorized(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized access"})
}

// AuthMiddleware returns a middleware that populate the Gin Context with security data
func AuthMiddleware(policy security.CompiledPolicy) gin.HandlerFunc {
	return func(c *gin.Context) {
		roles := security.BuildRolePolicy(policy, c.Request)
		c.Set(KeyRoles, roles)
		log := GetLogger(c)
		if log != nil {
			log.WithField("roles", roles.ToSlice()).Debugf("user has roles")
		}
		c.Next()
	}
}

// GetLogger returns the request logger from the gin context, nil if it does not exist
func GetLogger(c *gin.Context) *logrus.Entry {
	logRaw, found := c.Get(KeyLogEntry)
	log, castok := logRaw.(*logrus.Entry)
	if !found || !castok {
		return nil
	}
	return log
}
