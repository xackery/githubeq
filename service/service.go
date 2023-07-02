package service

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/xackery/githubeq/config"
	"github.com/xackery/githubeq/database"
	"github.com/xackery/githubeq/github"
	"github.com/xackery/githubeq/tlog"
)

var Version string

func Start() (err error) {
	tlog.Init(nil, os.Stdout)
	if Version == "" {
		Version = "1.x.x EXPERIMENTAL"
	}
	tlog.Infof("Starting GithubEQ %s", Version)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := config.NewConfig(ctx)
	if err != nil {
		return fmt.Errorf("config: %w", err)
	}

	err = prepDatabase(cfg)
	if err != nil {
		return fmt.Errorf("prepDatabase: %w", err)
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)

	ticker := time.NewTicker(time.Duration(cfg.CheckFrequencyMinutes) * time.Minute)

	err = sync(cfg)
	if err != nil {
		tlog.Errorf("Error syncing: %s", err.Error())
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
			err = sync(cfg)
			if err != nil {
				tlog.Errorf("Error syncing: %s", err.Error())
			}
		}
	}
}

func sync(cfg *config.Config) error {
	var issues []database.Issue
	var err error

	tlog.Infof("Starting sync at %s", time.Now().Format(time.RFC3339))
	issues, err = syncIssuesFromDB(cfg, false)
	if err != nil {
		tlog.Errorf("Failed to get old issues from DB: %s", err.Error())
	}
	tlog.Infof("[DB] %d old issues", len(issues))

	//Get any updates from github on previously created issues
	issues, err = github.SyncUpdatesOnIssues(cfg, issues)
	if err != nil {
		return fmt.Errorf("syncUpdatesOnIssues: %w", err)
	}

	//Update DB on any changes
	err = syncDBWithGithubChanges(cfg, issues)
	if err != nil {
		return fmt.Errorf("syncDBWithGithubChanges: %w", err)
	}

	//=== Get New Issues, and sync it with Github ===
	issues, err = syncIssuesFromDB(cfg, true)
	if err != nil {
		return fmt.Errorf("syncIssuesFromDB: %w", err)
	}

	tlog.Infof("[DB] %d new issues", len(issues))

	issues, err = github.CreateIssues(cfg, issues)
	if err != nil {
		return fmt.Errorf("createIssues on Github: %w", err)
	}

	tlog.Infof("[Github] %d added issues", len(issues))

	err = syncDBWithGithubChanges(cfg, issues)
	if err != nil {
		return fmt.Errorf("syncDBWithGithubChanges: %w", err)
	}

	return nil
}

// Get any issues in DB without a github issue #
func syncIssuesFromDB(cfg *config.Config, isNewIssue bool) ([]database.Issue, error) {
	var issues []database.Issue

	db, err := sqlx.Connect("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&multiStatements=true&interpolateParams=true&collation=utf8mb4_unicode_ci&charset=utf8mb4,utf8", cfg.Database.Username, cfg.Database.Password, cfg.Database.Host, cfg.Database.Port, cfg.Database.Database))
	if err != nil {
		return nil, fmt.Errorf("connect db: %w", err)
	}
	defer db.Close()

	query := `SELECT * FROM bug_reports WHERE github_issue_id != 0 AND bug_status = 0`
	if isNewIssue {
		query = `SELECT * FROM bug_reports WHERE github_issue_id = 0 AND bug_status = 0`
	}

	rows, err := db.Queryx(query)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}

	for rows.Next() {
		issue := database.Issue{
			DB: &database.BugReports{},
		}
		err = rows.StructScan(issue.DB)
		if err != nil {
			return nil, fmt.Errorf("struct scan: %w", err)
		}
		issues = append(issues, issue)
	}
	return issues, nil
}

func syncDBWithGithubChanges(cfg *config.Config, addedIssues []database.Issue) error {
	db, err := sqlx.Connect("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&multiStatements=true&interpolateParams=true&collation=utf8mb4_unicode_ci&charset=utf8mb4,utf8", cfg.Database.Username, cfg.Database.Password, cfg.Database.Host, cfg.Database.Port, cfg.Database.Database))
	if err != nil {
		return fmt.Errorf("connect db: %w", err)
	}
	defer db.Close()

	for _, issue := range addedIssues {
		isFixed := 0
		if issue.Github.ClosedAt != nil {
			isFixed = 1
		}

		_, err = db.Exec("UPDATE bug_reports SET github_issue_id = ?, bug_status = ? WHERE id = ?", *issue.Github.Number, isFixed, issue.DB.ID)
		if err != nil {
			return fmt.Errorf("update %d: %w", issue.DB.ID, err)
		}
	}
	return nil
}

func prepDatabase(cfg *config.Config) error {
	db, err := sqlx.Connect("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&multiStatements=true&interpolateParams=true&collation=utf8mb4_unicode_ci&charset=utf8mb4,utf8", cfg.Database.Username, cfg.Database.Password, cfg.Database.Host, cfg.Database.Port, cfg.Database.Database))
	if err != nil {
		return fmt.Errorf("connect db: %w", err)
	}
	defer db.Close()

	_, err = db.Exec(`ALTER TABLE bug_reports ADD github_issue_id int(11) unsigned NOT NULL DEFAULT 0`)
	if err != nil {
		if strings.Contains(err.Error(), "Duplicate column name 'github_issue_id'") {
			return nil
		}
		return fmt.Errorf("alter table: %w", err)
	}
	tlog.Infof("Successfully added github_issue_id column to bug_reports table")
	return nil
}
