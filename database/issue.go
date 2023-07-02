package database

import (
	gh "github.com/google/go-github/github"
)

type Issue struct {
	DB     *BugReports
	Github *gh.Issue
}
