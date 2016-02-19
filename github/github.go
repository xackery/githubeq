package github

import (
	"fmt"
	gh "github.com/google/go-github/github"
	"github.com/xackery/eqemuconfig"
	"github.com/xackery/githubeq/database"

	//"golang.org/x/oauth2"
)

var config *eqemuconfig.Config
var authClient *gh.Client

func GetClient() (client *gh.Client, err error) {
	//Client is already created, return it
	if authClient != nil {
		client = authClient
		return
	}

	//Config isn't set, load it
	if config == nil {
		config, err = eqemuconfig.GetConfig()
		if err != nil {
			fmt.Errorf("Error getting eqemuconfig: %s", err.Error())
			return
		}
	}

	client = gh.NewClient(nil)

	//ts := oauth2.TokenSource{Token: &oauth2.Token{AccessToken: config.Shortname}}
	/*oauth2.ReuseTokenSource(t, src)

	ts := &oauth2.Token{AccessToken: config.Shortname}
	tc := oauth2.NewClient(oauth2.NoContext, ts)


	client = gh.NewClient(tc)
	authClient = client*/
	return
	//ts := oauth2.StaticTokenSource()
	//client := gh.NewClient(nil)

}

func CreateIssues(issues []database.Issue) (newIssues []database.Issue, err error) {
	if config == nil {
		config, err = eqemuconfig.GetConfig()
		if err != nil {
			fmt.Errorf("Error getting eqemuconfig: %s", err.Error())
			return
		}
	}

	client, err := GetClient()
	if err != nil {
		return
	}

	//milestones, resp, err := client.Issues.ListMilestones("Xackery", "xackery/rebuildeq", nil)

	//for _, milestone := range mielstones {
	//if milestones[0].Title == config.Github
	//}

	for _, issue := range issues {
		newIssueRequest := &gh.IssueRequest{}
		//msg, err := issue.DB.Message.Value()
		//title := fmt.Sprintf("[#%d %s] %s", issue.DB.Id, issue.DB.My_name, msg[0:25])
		//newIssue.Title = &title
		newIssue, _, newErr := client.Issues.Create("Xackery", "rebuildeq", newIssueRequest)
		if newErr != nil {
			err = newErr
			return
		}
		//issue.
		issue.Github = newIssue
		newIssues = append(newIssues, issue)
	}
	//req.Milestone
	return
}
