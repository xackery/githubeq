package github

import (
	"context"
	"fmt"

	gh "github.com/google/go-github/github"
	"github.com/xackery/githubeq/database"

	"time"

	"github.com/xackery/githubeq/config"
	"golang.org/x/oauth2"
)

var authClient *gh.Client

type TokenSource struct {
	AccessToken string
}

func (t *TokenSource) Token() (*oauth2.Token, error) {
	token := &oauth2.Token{
		AccessToken: t.AccessToken,
	}
	return token, nil
}

func Client(cfg *config.Config) (client *gh.Client, err error) {

	//Client is already created, return it
	if authClient != nil {
		client = authClient
		return
	}

	//return client
	ts := &TokenSource{AccessToken: cfg.Github.PersonalAccessToken}
	tc := oauth2.NewClient(context.Background(), ts)
	client = gh.NewClient(tc)
	return
}

func SyncUpdatesOnIssues(cfg *config.Config, issues []database.Issue) (newIssues []database.Issue, err error) {

	client, err := Client(cfg)
	if err != nil {
		return
	}

	for _, issue := range issues {

		newIssue, resp, err := client.Issues.Get(context.Background(), cfg.Github.User, cfg.Github.Repository, issue.DB.GithubIssueID)
		if err != nil {
			return nil, fmt.Errorf("issues get: %s: %w", resp, err)
		}

		if newIssue.UpdatedAt.Before(issue.DB.LastReview) {
			continue
		}

		issue.Github = newIssue
		newIssues = append(newIssues, issue)
	}
	return
}

func CreateIssues(cfg *config.Config, issues []database.Issue) (newIssues []database.Issue, err error) {

	client, err := Client(cfg)
	if err != nil {
		return
	}

	for _, issue := range issues {
		//Create a new issue request
		newIssueRequest := gh.IssueRequest{}
		newIssueRequest.Labels = &[]string{}

		if issue.DB.CategoryName == "NPC" && cfg.Github.NPCLabel != "" {
			*newIssueRequest.Labels = append(*newIssueRequest.Labels, cfg.Github.NPCLabel)
		}

		//Truncate message to 25 max on title
		msg := issue.DB.BugReport
		if len(msg) > 25 {
			msg = msg[0:25] + "..."
		}

		//make title
		title := fmt.Sprintf("[#%d %s] %s", issue.DB.ID, issue.DB.ReporterName, msg)
		newIssueRequest.Title = &title

		//Create body
		body := fmt.Sprintf("**Message:** %s\n", issue.DB.BugReport)
		body = fmt.Sprintf("%s **User:** %s (cid: %d, accid: %d, client: %d) at %s\n", body, issue.DB.CharacterName, issue.DB.CharacterID, issue.DB.AccountID, issue.DB.ClientVersionID, issue.DB.ReportDatetime.Format(time.RFC822))
		body = fmt.Sprintf("%s **Location:** #zone %s %d %d %d\n", body, issue.DB.Zone, issue.DB.PosX, issue.DB.PosY, issue.DB.PosZ)
		if issue.DB.TargetID > 0 {
			body = fmt.Sprintf("%s **Target:** %s\n", body, issue.DB.TargetName)
		} else {
			body = fmt.Sprintf("%s **Target:** None\n", body)
		}

		isTrue := issue.DB.CanDuplicate == 1
		body = fmt.Sprintf("%s CanDuplicate? %t\n", body, isTrue)
		isTrue = issue.DB.CrashBug == 1
		body = fmt.Sprintf("%s CrashBug? %t\n\n", body, isTrue)

		body = fmt.Sprintf("%s System Info: %s\n", body, issue.DB.SystemInfo)

		newIssueRequest.Body = &body

		newIssue, resp, err := client.Issues.Create(context.Background(), cfg.Github.User, cfg.Github.Repository, &newIssueRequest)
		if err != nil {
			return nil, fmt.Errorf("issue create %s: %w", resp, err)
		}

		issue.Github = newIssue
		newIssues = append(newIssues, issue)
	}
	return
}
