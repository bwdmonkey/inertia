package webhook

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

// Constants for the generic webhook interface
var (
	PushEvent = "push"
	// PullEvent = "pull"
)

// Payload represents a generic webhook payload
type Payload interface {
	GetEventType() string
	GetRepoName() string
	GetRef() string
	GetGitURL() string
	GetSSHURL() string
}

// Parse takes in a webhook request and parses it into one of the supported types
func Parse(r *http.Request) (Payload, error) {
	if r.Header.Get("content-type") != "application/json" {
		return nil, errors.New("Content-Type must be JSON")
	}

	// Try Github
	githubEventHeader := r.Header.Get("x-github-event")
	if len(githubEventHeader) > 0 {
		fmt.Println("Github webhook detected")
		return parseGithubEvent(r, githubEventHeader)
	}

	// Try Gitlab
	gitlabEventHeader := r.Header.Get("x-gitlab-event")
	if len(gitlabEventHeader) > 0 {
		fmt.Println("Gitlab webhook detected")
		return nil, errors.New("Unsupported webhook received")
	}

	// Try Bitbucket
	userAgent := r.Header.Get("user-agent")
	if strings.Contains(userAgent, "Bitbucket") {
		fmt.Println("Bitbucket webhook detected")
		return nil, errors.New("Unsupported webhook received")
	}

	return nil, errors.New("Unsupported webhook received")
}
