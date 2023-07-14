package database

import "time"

// BugReports represents the bug reports table
type BugReports struct {
	ID                 int       `db:"id"`                   // int(11) unsigned NOT NULL AUTO_INCREMENT,
	Zone               string    `db:"zone"`                 // varchar(32) NOT NULL DEFAULT 'Unknown',
	ClientVersionID    int       `db:"client_version_id"`    // int(11) unsigned NOT NULL DEFAULT 0,
	ClientVersionName  string    `db:"client_version_name"`  // varchar(24) NOT NULL DEFAULT 'Unknown',
	AccountID          int       `db:"account_id"`           // int(11) unsigned NOT NULL DEFAULT 0,
	CharacterID        int       `db:"character_id"`         // int(11) unsigned NOT NULL DEFAULT 0,
	CharacterName      string    `db:"character_name"`       // varchar(64) NOT NULL DEFAULT 'Unknown',
	ReporterSpoof      int       `db:"reporter_spoof"`       // tinyint(1) NOT NULL DEFAULT 1,
	CategoryID         int       `db:"category_id"`          // int(11) unsigned NOT NULL DEFAULT 0,
	CategoryName       string    `db:"category_name"`        // varchar(64) NOT NULL DEFAULT 'Other',
	ReporterName       string    `db:"reporter_name"`        // varchar(64) NOT NULL DEFAULT 'Unknown',
	UIPath             string    `db:"ui_path"`              // varchar(128) NOT NULL DEFAULT 'Unknown',
	PosX               int       `db:"pos_x"`                // float NOT NULL DEFAULT 0,
	PosY               int       `db:"pos_y"`                // float NOT NULL DEFAULT 0,
	PosZ               int       `db:"pos_z"`                // float NOT NULL DEFAULT 0,
	Heading            int       `db:"heading"`              // int(11) unsigned NOT NULL DEFAULT 0,
	TimePlayed         int       `db:"time_played"`          // int(11) unsigned NOT NULL DEFAULT 0,
	TargetID           int       `db:"target_id"`            // int(11) unsigned NOT NULL DEFAULT 0,
	TargetName         string    `db:"target_name"`          // varchar(64) NOT NULL DEFAULT 'Unknown',
	OptionalInfoMask   int       `db:"optional_info_mask"`   // int(11) unsigned NOT NULL DEFAULT 0,
	CanDuplicate       int       `db:"_can_duplicate"`       // tinyint(1) NOT NULL DEFAULT 0,
	CrashBug           int       `db:"_crash_bug"`           // tinyint(1) NOT NULL DEFAULT 0,
	TargetInfo         int       `db:"_target_info"`         // tinyint(1) NOT NULL DEFAULT 0,
	CharacterFlags     int       `db:"_character_flags"`     // tinyint(1) NOT NULL DEFAULT 0,
	UnknownValue       int       `db:"_unknown_value"`       // tinyint(1) NOT NULL DEFAULT 0,
	BugReport          string    `db:"bug_report"`           // varchar(1024) NOT NULL DEFAULT '',
	SystemInfo         string    `db:"system_info"`          // varchar(1024) NOT NULL DEFAULT '',
	ReportDatetime     time.Time `db:"report_datetime"`      // datetime NOT NULL DEFAULT current_timestamp(),
	BugStatus          int       `db:"bug_status"`           // tinyint(3) unsigned NOT NULL DEFAULT 0,
	LastReview         time.Time `db:"last_review"`          // datetime NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
	LastReviewer       string    `db:"last_reviewer"`        // varchar(64) NOT NULL DEFAULT 'None',
	ReviewerNotes      string    `db:"reviewer_notes"`       // varchar(1024) NOT NULL DEFAULT '',
	GithubIssueID      int       `db:"github_issue_id"`      // int(11) unsigned NOT NULL DEFAULT 0,
	GithubSyncDatetime time.Time `db:"github_sync_datetime"` // datetime NOT NULL DEFAULT current_timestamp()
}
