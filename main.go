package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"time"

	gh "github.com/google/go-github/v60/github"
	"github.com/jamfesteq/githubeq/config"
	"github.com/jamfesteq/githubeq/tlog"
	"golang.org/x/oauth2"
)

var Version string
var cfg *config.Config
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

func main() {
	err := run()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func run() error {
	tlog.Init(nil, os.Stdout)
	if Version == "" {
		Version = "1.x.x EXPERIMENTAL"
	}
	tlog.Infof("Starting GithubEQ %s", Version)


	err := os.MkdirAll("issues", 0755)
	if err != nil {
		return fmt.Errorf("mkdir: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err = config.NewConfig(ctx)
	if err != nil {
		return fmt.Errorf("config: %w", err)
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)

	ticker := time.NewTicker(time.Duration(cfg.SyncFrequencyMinutes) * time.Minute)

	err = generate()
	if err != nil {
		return fmt.Errorf("generate: %w", err)
	}
	for {
		select {
		case <-ctx.Done():
			tlog.Infof("Exiting, interrupt signal sent")
			return nil
		case <-signalChan:
			tlog.Infof("Exiting, interrupt signal sent")
			return nil
		case <-ticker.C:
			err = generate()
			if err != nil {
				tlog.Errorf("Failed to generate: %s", err.Error())
				continue
			}
		}
	}
}

func generate() error {
	files, err := os.ReadDir("issues")
	if err != nil {
		return fmt.Errorf("read dir: %w", err)
	}
	totalParsed := 0
	for _, file := range files {
		data, err := os.ReadFile("issues/" + file.Name())
		if err != nil {
			return fmt.Errorf("read file: %w", err)
		}

		if filepath.Ext(file.Name()) != ".txt" {
			tlog.Debugf("Skipping file: %s", file.Name())
			continue
		}

		err = createIssue(string(data))
		if err != nil {
			return fmt.Errorf("create %s: %w", file.Name(), err)
		}

		err = os.Remove("issues/" + file.Name())
		if err != nil {
			return fmt.Errorf("remove file: %w", err)
		}

		totalParsed++
	}
	tlog.Infof("Total Parsed: %d", totalParsed)

	return nil
}

func createIssue(data string) error {

	newIssueRequest := gh.IssueRequest{}
	newIssueRequest.Labels = &[]string{}

	offset := strings.Index(data, "\n")
	if offset == -1 {
		return fmt.Errorf("author offset: -1")
	}
	author := data[:offset]
	data = data[offset+1:]

	offset = strings.Index(data, "\n")
	if offset == -1 {
		return fmt.Errorf("tags offset: -1")
	}
	tagData := data[:offset]
	data = data[offset+1:]

	tags := strings.Split(tagData, ",")

	offset = strings.Index(data, "\n")
	if offset == -1 {
		return fmt.Errorf("msg offset: -1")
	}
	title := data[:offset]
	body := title + "\n"
	data = data[offset+1:]
	if len(title) > 25 {
		title = title[:25] + "..."
	}
	title = fmt.Sprintf("[%s] %s", author, title)
	newIssueRequest.Title = &title
	if len(tags) > 0 {

	}
	//	*newIssueRequest.Labels = append(*newIssueRequest.Labels, tags...)

	offset = strings.Index(data, "-------\n")
	if offset == -1 {
		return fmt.Errorf("dataBody offset -1")
	}

	body += strings.ReplaceAll(data[0:offset], "\n", "<br>\n") + "\n"
	data = data[offset+8:]

	body += "<details>\n<summary>"
	offset = strings.Index(data, "\n")
	if offset == -1 {
		return fmt.Errorf("charInfo offset -1")
	}
	charInfo := data[:offset]
	body += charInfo
	body += "</summary>\n"
	data = data[offset+1:]

	offset = strings.Index(data, "-------\n")
	if offset == -1 {
		return fmt.Errorf("dataInventory offset -1")
	}
	body += strings.ReplaceAll(data[0:offset], "\n", "<br>\n") + "\n"
	body += "</details>\n"
	data = data[offset+8:]

	offset = strings.Index(data, "\n")
	if offset == -1 {
		return fmt.Errorf("inventory offset -1")
	}
	inventory := data[:offset]
	body += "<details>\n<summary>"
	body += inventory
	body += "</summary>\n"
	body += strings.ReplaceAll(data[offset+1:], "\n", "<br>\n") + "\n"
	body += "</details>\n"

	// fmt.Println("title:", title)
	// fmt.Println("tags:", tags)
	// fmt.Println("body:", body)
	newIssueRequest.Body = &body

	client := gh.NewClient(nil).WithAuthToken(cfg.Github.PersonalAccessToken)

	newIssue, resp, err := client.Issues.Create(context.Background(), cfg.Github.User, cfg.Github.Repository, &newIssueRequest)
	if err != nil {
		return fmt.Errorf("github create %+v: %w", resp, err)
	}

	tlog.Infof("Created issue: %d from %s %s", newIssue.ID, author, newIssue.GetHTMLURL())

	limit, resp, err := client.RateLimit.Get(context.Background())
	if err != nil {
		return fmt.Errorf("ratelimits: %v: %w", resp, err)
	}

	if limit.Core.Remaining < 5 {
		tlog.Infof("[Github] ratelimit close, ending early: %d/%d, resets in %0.2f minutes", limit.Core.Remaining, limit.Core.Limit, time.Until(limit.Core.Reset.Time).Minutes())
		return nil
	}

	return nil
}
