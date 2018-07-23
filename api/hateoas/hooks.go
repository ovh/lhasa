package hateoas

import (
	"github.com/gin-gonic/gin"
	"github.com/loopfz/gadgeto/tonic"
)

// RenderHookWrapper handles hateoas entity conversion
func RenderHookWrapper(hook tonic.RenderHook) tonic.RenderHook {
	return func(c *gin.Context, status int, payload interface{}) {
		baseURL := BaseURL(c)
		switch r := payload.(type) {
		case Resourceable:
			r.ToResource(baseURL)
		}
		if hook == nil {
			return
		}
		hook(c, status, payload)
	}
}
