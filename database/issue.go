package database

import (
	gh "github.com/google/go-github/github"
	"github.com/xackery/goeq/issue"
)

type Issue struct {
	DB     issue.Issue
	Github *gh.Issue
}
