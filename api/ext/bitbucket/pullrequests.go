package bitbucket

import (
	"encoding/json"
)

// PullRequests client
type PullRequests struct {
	c *Client
}

// Create method
func (p *PullRequests) Create(po *PullRequestsOptions) (interface{}, error) {
	data := p.buildPullRequestBody(po)
	urlStr := p.c.requestURL("/projects/%s/repos/%s/pull-requests", po.Owner, po.RepoSlug)
	return p.c.execute("POST", urlStr, data, "application/json")
}

func (p *PullRequests) buildPullRequestBody(po *PullRequestsOptions) string {

	body := map[string]interface{}{}
	body["title"] = ""
	body["description"] = ""
	body["state"] = "OPEN"
	body["open"] = true
	body["closed"] = false
	body["locked"] = false
	body["reviewers"] = make([]interface{}, len(po.Reviewers))

	if n := len(po.Reviewers); n > 0 {
		for i, user := range po.Reviewers {
			body["reviewers"].([]interface{})[i] = map[string]interface{}{"user": nil}
			body["reviewers"].([]interface{})[i].(map[string]interface{})["user"] = map[string]interface{}{"name": user}
		}
	}

	body["fromRef"] = map[string]interface{}{}
	body["toRef"] = map[string]interface{}{}

	body["fromRef"].(map[string]interface{})["id"] = "refs/heads/" + po.SourceBranch
	body["fromRef"].(map[string]interface{})["repository"] = map[string]interface{}{"slug": po.RepoSlug}
	body["fromRef"].(map[string]interface{})["repository"].(map[string]interface{})["project"] = map[string]interface{}{"key": po.Owner}
	body["toRef"].(map[string]interface{})["id"] = "refs/heads/master"
	body["toRef"].(map[string]interface{})["repository"] = map[string]interface{}{"slug": po.RepoSlug}
	body["toRef"].(map[string]interface{})["repository"].(map[string]interface{})["project"] = map[string]interface{}{"key": po.Owner}

	if po.Title != "" {
		body["title"] = po.Title
	}

	if po.Description != "" {
		body["description"] = po.Description
	}

	if po.Message != "" {
		body["message"] = po.Message
	}

	data, err := json.MarshalIndent(body, "", "\t")
	if err != nil {
		return err.Error()
	}
	return string(data)
}
