package config

import (
	"context"
	"fmt"
	"os"
	"runtime"

	"github.com/jbsmith7741/toml"
	"github.com/rs/zerolog"
)

// Config represents a configuration parse
type Config struct {
	Debug                bool `toml:"debug" desc:"GithubEQ Configuration\n\n# Debug messages are displayed. This will cause console to be more verbose, but also more informative"`
	SyncFrequencyMinutes int  `toml:"sync_frequency_minutes" desc:"How often should stale bugs/issues be checked for state changes."`
	Github               ConfigGithub
}

type ConfigGithub struct {
	PersonalAccessToken string `toml:"personal_access_token" desc:"Personal access token for github"`
	Repository          string `toml:"repository" desc:"Repository name, e.g. githubeq in jamfesteq/githubeq"`
	User                string `toml:"user" desc:"User name the repo is in, e.g. jamfesteq in jamfesteq/githubeq"`
	FallbackLabel       string `toml:"fallback_label" desc:"Optional label to use if no other label is found"`
	BugLabel            string `toml:"bug_label" desc:"Optional label, all reported bugs use this label if set"`
	OtherLabel          string `toml:"other_label" desc:"Optional label to use for bugs with other report"`
	VideoLabel          string `toml:"video_label" desc:"Optional label to use for bugs with video report"`
	AudioLabel          string `toml:"audio_label" desc:"Optional label to use for bugs with audio report"`
	PathingLabel        string `toml:"pathing_label" desc:"Optional label to use for bugs with pathing report"`
	QuestLabel          string `toml:"quest_label" desc:"Optional label to use for bugs with quest report"`
	TradeskillsLabel    string `toml:"tradeskills_label" desc:"Optional label to use for bugs with tradeskills report"`
	SpellStackingLabel  string `toml:"spell_stacking_label" desc:"Optional label to use for bugs with spell stacking report"`
	DoorsPortalLabel    string `toml:"doors_portal_label" desc:"Optional label to use for bugs with doors/portal report"`
	ItemsLabel          string `toml:"items_label" desc:"Optional label to use for bugs with items report"`
	NPCLabel            string `toml:"npc_label" desc:"Optional label to use for bugs with npc report"`
	DialogsLabel        string `toml:"dialogs_label" desc:"Optional label to use for bugs with dialogs report"`
	LoNLabel            string `toml:"lon_label" desc:"Optional label to use for bugs with LoN report"`
	MercenariesLabel    string `toml:"mercenaries_label" desc:"Optional label to use for bugs with mercenaries report"`
}

// NewConfig creates a new configuration
func NewConfig(ctx context.Context) (*Config, error) {
	var f *os.File
	cfg := Config{}
	path := "githubeq.conf"

	isNewConfig := false
	fi, err := os.Stat(path)
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, fmt.Errorf("config info: %w", err)
		}
		f, err = os.Create(path)
		if err != nil {
			return nil, fmt.Errorf("create githubeq.conf: %w", err)
		}
		fi, err = os.Stat(path)
		if err != nil {
			return nil, fmt.Errorf("new config info: %w", err)
		}
		isNewConfig = true
	}
	if !isNewConfig {
		f, err = os.Open(path)
		if err != nil {
			return nil, fmt.Errorf("open config: %w", err)
		}
	}

	defer f.Close()
	if fi.IsDir() {
		return nil, fmt.Errorf("githubeq.conf is a directory, should be a file")
	}

	if isNewConfig {
		enc := toml.NewEncoder(f)
		enc.Encode(defaultLabel())

		fmt.Println("a new githubeq.conf file was created. Please open this file and configure githubeq, then run it again.")
		if runtime.GOOS == "windows" {
			option := ""
			fmt.Println("press a key then enter to exit.")
			fmt.Scan(&option)
		}
		os.Exit(0)
	}

	_, err = toml.DecodeReader(f, &cfg)
	if err != nil {
		return nil, fmt.Errorf("decode githubeq.conf: %w", err)
	}

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if cfg.Debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	err = cfg.Verify()
	if err != nil {
		return nil, fmt.Errorf("verify: %w", err)
	}

	return &cfg, nil
}

// Verify returns an error if configuration appears off
func (c *Config) Verify() error {
	if c.SyncFrequencyMinutes < 1 {
		c.SyncFrequencyMinutes = 1
	}

	return nil
}

func defaultLabel() Config {
	cfg := Config{
		Debug:                true,
		SyncFrequencyMinutes: 1,
		Github: ConfigGithub{
			BugLabel: "bug",
		},
	}

	return cfg
}
