package service

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/xackery/eqemuconfig"
	"github.com/xackery/githubeq/database"
	"github.com/xackery/githubeq/github"
	"log"
	"time"
)

var config *eqemuconfig.Config

func Start() (err error) {
	config, err = eqemuconfig.GetConfig()
	if err != nil {
		err = fmt.Errorf("eqemuconfig error:", err.Error())
		return
	}

	for {
		time.Sleep(1 * time.Second)
		var issues []database.Issue

		//=== Get Old Issues, and sync it with Github ===
		issues, err = getIssuesFromDB(false)
		if err != nil {
			fmt.Errorf("Error getting old issues from DB: %s", err.Error())
		}
		log.Printf("[DB] %d old issues", len(issues))

		//Get any updates from github on previously created issues
		issues, err = github.GetUpdatesOnIssues(issues)
		if err != nil {
			fmt.Errorf("Error getting updates on issues: %s", err.Error())
			return
		}

		//Update DB on any changes
		err = updateDBWithGithubChanges(issues)
		if err != nil {
			fmt.Errorf("Issues upding DB with github changes: %s", err.Error())
			return
		}

		//=== Get New Issues, and sync it with Github ===
		issues, err = getIssuesFromDB(true)
		if err != nil {
			fmt.Errorf("Error getting new issues from DB: %s", err.Error())
			return
		}

		if len(issues) < 1 {
			return
		}
		log.Printf("[DB] %d new issues", len(issues))

		issues, err = github.CreateIssues(issues)
		if err != nil {
			fmt.Errorf("Had issues adding new issues to Github: %s", err.Error())
			return
		}

		log.Printf("[Github] %d added issues", len(issues))

		err = updateDBWithGithubChanges(issues)
		if err != nil {
			fmt.Errorf("Issues upding DB with github changes: %s", err.Error())
			return
		}
		log.Println("Done!")
		return
	}

	return
}

//Get any issues in DB without a github issue #
func getIssuesFromDB(isNewIssue bool) (issues []database.Issue, err error) {
	db, err := sqlx.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=true", config.Database.Username, config.Database.Password, config.Database.Host, config.Database.Port, config.Database.Db))
	if err != nil {
		return
	}

	query := `SELECT * from issues WHERE github_issue_id != 0 AND is_fixed = 0`
	if isNewIssue {
		query = `SELECT * from issues WHERE github_issue_id = 0 AND is_fixed = 0`
	}

	rows, err := db.Queryx(query)
	if err != nil {
		fmt.Errorf("Error getting non-issued issues: %s", err.Error())
		return
	}
	defer db.Close()

	for rows.Next() {
		issue := database.Issue{}
		err = rows.StructScan(&issue.DB)
		if err != nil {
			fmt.Errorf("Error scanning issue to struct: %s", err.Error())
			return
		}
		issues = append(issues, issue)
	}
	return
}

func updateDBWithGithubChanges(addedIssues []database.Issue) (err error) {
	db, err := sqlx.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=true", config.Database.Username, config.Database.Password, config.Database.Host, config.Database.Port, config.Database.Db))
	if err != nil {
		return
	}
	defer db.Close()

	for _, issue := range addedIssues {
		isFixed := 0
		if issue.Github.ClosedAt != nil {
			isFixed = 1
		}

		isInProgress := 0
		if issue.Github.Assignee != nil {
			isInProgress = 1
		}

		_, err = db.Exec("UPDATE issues SET github_issue_id = ?, is_fixed = ?, is_in_progress = ? WHERE id = ?", *issue.Github.Number, isFixed, isInProgress, issue.DB.Id)
		if err != nil {
			err = fmt.Errorf("Error updating github issue for id %d: %s", issue.DB.Id, err.Error())
			return
		}
	}
	return
}
