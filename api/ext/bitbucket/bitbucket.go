package bitbucket

import (
	"github.com/ovh/lhasa/api/config"
)

func getAPIBaseURL() string {
	baseURL := ""
	bitbucket, ok := config.ExtractValue("bitbucket").(map[string]interface{})
	if ok {
		baseURL, _ = bitbucket["url"].(string)
	}
	return baseURL
}

// RepositoriesOptions options
type RepositoriesOptions struct {
	Owner string `json:"owner"`
	Team  string `json:"team"`
	Role  string `json:"role"`
}

// RepositoryOptions options
type RepositoryOptions struct {
	Owner       string `json:"owner"`
	RepoSlug    string `json:"repo_slug"`
	Scm         string `json:"scm"`
	Description string `json:"description"`
	Language    string `json:"language"`
	Project     string `json:"project"`
}

// PullRequestsOptions options
type PullRequestsOptions struct {
	Owner             string   `json:"owner"`
	RepoSlug          string   `json:"repo_slug"`
	Title             string   `json:"title"`
	Description       string   `json:"description"`
	CloseSourceBranch bool     `json:"close_source_branch"`
	SourceBranch      string   `json:"source_branch"`
	SourceRepository  string   `json:"source_repository"`
	DestinationBranch string   `json:"destination_branch"`
	DestinationCommit string   `json:"destination_repository"`
	Message           string   `json:"message"`
	Reviewers         []string `json:"reviewers"`
}

// BranchOptions option
type BranchOptions struct {
	Owner      string `json:"owner"`
	RepoSlug   string `json:"repo_slug"`
	Name       string `json:"name"`
	StartPoint string `json:"startPoint"`
	Message    string `json:"message"`
}

// PathOptions options
type PathOptions struct {
	Owner    string `json:"owner"`
	RepoSlug string `json:"repo_slug"`
	Name     string `json:"name"`
	Branch   string `json:"branch"`
	Content  string `json:"content"`
	Message  string `json:"message"`
}
