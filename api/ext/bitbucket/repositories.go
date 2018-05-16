package bitbucket

// Repositories client
type Repositories struct {
	c            *Client
	PullRequests *PullRequests
	Branch       *Branch
	Path         *Path
}
