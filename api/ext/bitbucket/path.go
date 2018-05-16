package bitbucket

import (
	"bytes"
	"mime/multipart"
)

// Path client
type Path struct {
	c *Client
}

// Create method
func (p *Path) Create(po *PathOptions) (interface{}, error) {
	data, content := p.buildPathBody(po)
	urlStr := p.c.requestURL("/projects/%s/repos/%s/browse/%s", po.Owner, po.RepoSlug, po.Name)
	return p.c.execute("PUT", urlStr, data, content)
}

func (p *Path) buildPathBody(po *PathOptions) (string, string) {

	// Prepare a form that you will submit to that URL.
	var buffer bytes.Buffer
	w := multipart.NewWriter(&buffer)
	// Add the other fields
	fw, err := w.CreateFormField("content")
	if err != nil {
		return "", "err"
	}
	if _, err := fw.Write([]byte(po.Content)); err != nil {
		return "", "err"
	}
	// Add the other fields
	if fw, err = w.CreateFormField("message"); err != nil {
		return "", "err"
	}
	if _, err = fw.Write([]byte(po.Message)); err != nil {
		return "", "err"
	}
	// Add the other fields
	if fw, err = w.CreateFormField("branch"); err != nil {
		return "", "err"
	}
	if _, err = fw.Write([]byte(po.Branch)); err != nil {
		return "", "err"
	}
	// Don't forget to close the multipart writer.
	// If you don't close it, your request will be missing the terminating boundary.
	w.Close()

	form := buffer.String()

	return form, w.FormDataContentType()
}
