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

// Version is the version of the application
var Version string
var cfg *config.Config

// TokenSource is a custom oauth2.TokenSource
type TokenSource struct {
	AccessToken string
}

// Token returns a new oauth2.Token
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

	offset := strings.Index(data, "-------\n")
	if offset == -1 {
		return fmt.Errorf("message offset: -1")
	}
	msg := data[0:offset]
	data = data[offset:]

	offset = strings.Index(msg, "\n")
	if offset == -1 {
		return fmt.Errorf("title offset -1")
	}
	title := msg[:offset]
	if len(title) > 25 {
		title = title[:25] + "..."
	}

	author := ""
	category := ""
	body := msg

	for {
		offset = strings.Index(data, "-------\n")
		if offset == -1 {
			break
		}
		data = data[offset+8:]
		offset = strings.Index(data, "\n")
		if offset == -1 {
			break
		}
		label := data[:offset]
		data = data[offset+1:]

		offset = strings.Index(data, "-------\n")
		if offset == -1 {
			offset = len(data)
		}

		content := data[0:offset]

		if label == "preview" {
			body += content
			data = data[offset:]
			continue
		}

		body += "<details>\n"
		body += fmt.Sprintf("<summary>%s</summary>\n", label)

		if label == "bug info" {
			lines := strings.Split(content, "\n")
			for _, line := range lines {
				if strings.HasPrefix(line, "category_name: ") {
					category = strings.TrimSpace(line[14:])
				}
			}
		}
		if label == "character info" {
			lines := strings.Split(content, "\n")
			for _, line := range lines {
				if strings.HasPrefix(line, "Account: ") {
					author = strings.TrimSpace(line[9:])
				}
			}
		}

		body += strings.ReplaceAll(content, "\n", "<br>\n")
		body += "</details>\n"

		data = data[offset:]
	}

	title = fmt.Sprintf("[%s] %s", author, title)
	if len(title) > 60 {
		title = title[:60] + "..."
	}

	labels := []string{}
	if cfg.Github.BugLabel != "" {
		labels = append(labels, cfg.Github.BugLabel)
	}

	label := labelByCategory(category)
	if label != "" && label != cfg.Github.BugLabel {
		labels = append(labels, label)
	}

	newIssueRequest.Title = &title
	newIssueRequest.Body = &body
	newIssueRequest.Labels = &labels

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

func labelByCategory(category string) string {
	labels := map[string]string{
		"Other":          cfg.Github.OtherLabel,
		"Video":          cfg.Github.VideoLabel,
		"Audio":          cfg.Github.AudioLabel,
		"Pathing":        cfg.Github.PathingLabel,
		"Quest":          cfg.Github.QuestLabel,
		"Tradeskills":    cfg.Github.TradeskillsLabel,
		"Spell Stacking": cfg.Github.SpellStackingLabel,
		"Doors/Portal":   cfg.Github.DoorsPortalLabel,
		"Items":          cfg.Github.ItemsLabel,
		"NPC":            cfg.Github.NPCLabel,
		"Dialogs":        cfg.Github.DialogsLabel,
		"LoN - TCG":      cfg.Github.LoNLabel,
		"Mercenaries":    cfg.Github.MercenariesLabel,
	}

	label, ok := labels[category]
	if !ok {
		label = cfg.Github.FallbackLabel
	}
	if label == "" {
		return cfg.Github.FallbackLabel
	}

	return label
}
