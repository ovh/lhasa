package binding

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/loopfz/gadgeto/tonic"
	"github.com/ovh/lhasa/api/v1"
)

var (
	// BINARY Simple binary binding
	BINARY = binaryBinding{}
	// WHITELIST Simple content white list
	WHITELIST = []string{"text/plain"}
)

type binaryBinding struct{}

func (binaryBinding) Name() string {
	return "binary"
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func (binaryBinding) Bind(req *http.Request, obj interface{}) error {
	// Scan interface implementation
	media, assert := (obj).(v1.MediaResource)
	if assert {
		// Check white list
		for header, values := range req.Header {
			if strings.ToLower(header) == "content-type" {
				for _, value := range values {
					if !stringInSlice(value, WHITELIST) {
						return errors.New("content-type " + value + " is white listed")
					}
				}
			}
		}

		buf, err := ioutil.ReadAll(req.Body)
		if err != nil {
			return err
		}
		media.SetBytes(buf)
	}

	return nil
}

// BindHook hook for lhasa, override default one
func BindHook(c *gin.Context, i interface{}) error {
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, tonic.MaxBodyBytes)
	if c.Request.ContentLength == 0 || c.Request.Method == http.MethodGet {
		return nil
	}
	// If tag binary is on any field binding will override json binding
	if err := c.ShouldBindWith(i, BINARY); err != nil && err != io.EOF {
		return fmt.Errorf("error parsing request body: %s", err.Error())
	}
	if err := c.ShouldBindWith(i, binding.JSON); err != nil && err != io.EOF {
		return fmt.Errorf("error parsing request body: %s", err.Error())
	}
	return nil
}

// RenderHook hook for lhasa, override default one
func RenderHook(c *gin.Context, status int, payload interface{}) {
	// Scan interface implementation
	media, assert := (payload).(v1.MediaResource)
	if assert {
		c.Data(http.StatusOK, media.GetContentType(), media.GetBytes())
		return
	}

	tonic.DefaultRenderHook(c, status, payload)
}
