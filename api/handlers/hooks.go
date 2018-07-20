package handlers

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/juju/errors"
	"github.com/loopfz/gadgeto/tonic"
)

var (
	// binary Simple binary binding
	binary binaryBinding
	// whitelist Simple content white list
	whitelist = []string{"text/plain", "application/json"}
)

// MediaResource defines a media resource behaviour
type MediaResource interface {
	GetContentType() string
	GetBytes() []byte
	SetBytes([]byte)
}

type binaryBinding struct{}

// Name returns "binary" as a name for this binding
func (binaryBinding) Name() string {
	return "binary"
}

// Bind perform this binary binding
func (binaryBinding) Bind(req *http.Request, obj interface{}) error {
	// Scan interface implementation
	media, ok := obj.(MediaResource)
	if !ok {
		return nil
	}
	// Check white list
	for _, value := range req.Header["Content-Type"] {
		if !stringInSlice(value, whitelist) {
			return fmt.Errorf("unsupported content-type %s", value)
		}
	}

	buf, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return err
	}
	media.SetBytes(buf)
	return nil
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// BindHook hook for lhasa, override default one
func BindHook(c *gin.Context, i interface{}) error {
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, tonic.MaxBodyBytes)
	if c.Request.ContentLength == 0 || c.Request.Method == http.MethodGet {
		return nil
	}
	defaultBindError := errors.New("error parsing request body")
	if err := c.ShouldBindWith(i, binary); err != nil && err != io.EOF {
		return errors.Wrap(err, defaultBindError)
	}
	if err := c.ShouldBindWith(i, binding.JSON); err != nil && err != io.EOF {
		return errors.Wrap(err, defaultBindError)
	}
	return nil
}

// RenderHook hook for lhasa, override default one
func RenderHook(c *gin.Context, status int, payload interface{}) {
	// Scan interface implementation
	media, ok := payload.(MediaResource)
	if ok {
		c.Data(http.StatusOK, media.GetContentType(), media.GetBytes())
		return
	}

	tonic.DefaultRenderHook(c, status, payload)
}
