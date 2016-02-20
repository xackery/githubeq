package github

import (
	"fmt"
	gh "github.com/google/go-github/github"
	"github.com/xackery/eqemuconfig"
	"github.com/xackery/githubeq/database"

	"golang.org/x/oauth2"
	"time"
)

var config *eqemuconfig.Config
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
	ts := &TokenSource{AccessToken: config.Github.PersonalAccessToken}
	tc := oauth2.NewClient(oauth2.NoContext, ts)
	client = gh.NewClient(tc)
	/*user, _, err := client.Users.Get("")
	if err != nil {
		fmt.Printf("client.Users.Get() faled with '%s'\n", err)
		return
	}*/
	//fmt.Println(user)
	return
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

	for _, issue := range issues {
		//Create a new issue request
		newIssueRequest := gh.IssueRequest{}
		newIssueRequest.Labels = &[]string{}

		//Add labels
		if issue.DB.Tar_is_client == 1 && config.Github.CharacterLabel != "" {
			*newIssueRequest.Labels = append(*newIssueRequest.Labels, config.Github.CharacterLabel)
		}
		if issue.DB.Tar_is_npc == 1 && config.Github.NPCLabel != "" {
			*newIssueRequest.Labels = append(*newIssueRequest.Labels, config.Github.NPCLabel)
		}
		if issue.DB.Item_id > 0 {
			*newIssueRequest.Labels = append(*newIssueRequest.Labels, config.Github.ItemLabel)
		}

		//Truncate message to 25 max on title
		msg := issue.DB.Message
		if len(msg) > 25 {
			msg = msg[0:25] + "..."
		}

		//make title
		title := fmt.Sprintf("[#%d %s] %s", issue.DB.Id, issue.DB.My_name, msg)
		newIssueRequest.Title = &title

		db := issue.DB
		//Create body
		body := fmt.Sprintf("%s\n", db.Message)
		body = fmt.Sprintf("%s **User:** %s (cid: %d, accid: %d, client: %s) at %s\n", body, db.My_name, db.My_character_id, db.My_account_id, db.Client, db.Create_date.Format(time.RFC822))
		body = fmt.Sprintf("%s **Location:** %f, %f, %f (zone: %d)\n", body, db.My_x, db.My_y, db.My_z, db.My_zone_id)
		if db.Tar_is_client > 0 {
			body = fmt.Sprintf("%s **Target:** Client %s (cid: %d, accid: %d)\n", body, db.Tar_name, db.Tar_character_id, db.Tar_account_id)
		} else if db.Tar_is_npc > 0 {
			body = fmt.Sprintf("%s **Target:** NPC %s (%d) spawngroup %d\n", body, db.Tar_name, db.Tar_npc_type_id, db.Tar_npc_spawngroup_id)
		} else {
			body = fmt.Sprintf("%s **Target:** None\n", body)
		}

		if db.Item_id > 0 {
			body = fmt.Sprintf("%s **Item**: %s (%d)", body, db.Item_name, db.Item_id)
		}
		newIssueRequest.Body = &body

		newIssue, resp, newErr := client.Issues.Create(config.Github.RepoUser, config.Github.RepoName, &newIssueRequest)
		if newErr != nil {
			fmt.Println(resp)
			err = fmt.Errorf("Failed to request issues: %s", newErr.Error())
			return
		}

		issue.Github = newIssue
		newIssues = append(newIssues, issue)
		return
	}
	//req.Milestone
	return
}
