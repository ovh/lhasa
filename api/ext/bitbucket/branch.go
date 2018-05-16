package bitbucket

import (
	"encoding/json"
)

// Branch client
type Branch struct {
	c *Client
}

// Create method
func (b *Branch) Create(bo *BranchOptions) (interface{}, error) {
	data := b.buildBranchBody(bo)
	urlStr := b.c.requestURL("/projects/%s/repos/%s/branches", bo.Owner, bo.RepoSlug)
	return b.c.execute("POST", urlStr, data, "application/json")
}

func (b *Branch) buildBranchBody(bo *BranchOptions) string {

	body := map[string]interface{}{}
	body["name"] = bo.Name
	body["startPoint"] = bo.StartPoint
	body["message"] = bo.Message

	data, err := json.MarshalIndent(body, "", "\t")
	if err != nil {
		return err.Error()
	}
	return string(data)
}
