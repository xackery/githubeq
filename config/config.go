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
	Debug                 bool `toml:"debug" desc:"GithubEQ Configuration\n\n# Debug messages are displayed. This will cause console to be more verbose, but also more informative"`
	CheckFrequencyMinutes int  `toml:"check_frequency_minutes" desc:"How often should github be queried for issues."`
	Database              ConfigDatabase
	Github                ConfigGithub
}

type ConfigDatabase struct {
	Host     string `toml:"host" desc:"Database host"`
	Port     int    `toml:"port" desc:"Database port"`
	Username string `toml:"username" desc:"Database username"`
	Password string `toml:"password" desc:"Database password"`
	Database string `toml:"database" desc:"Database name"`
}

type ConfigGithub struct {
	PersonalAccessToken string `toml:"personal_access_token" desc:"Personal access token for github"`
	Repository          string `toml:"repository" desc:"Repository name, e.g. githubeq in xackery/githubeq"`
	User                string `toml:"user" desc:"User name the repo is in, e.g. xackery in xackery/githubeq"`
	CharacterLabel      string `toml:"character_label" desc:"Label to use for character"`
	NPCLabel            string `toml:"npc_label" desc:"Label to use for npc"`
	ItemLabel           string `toml:"item_label" desc:"Label to use for item"`
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
		enc.Encode(getDefaultConfig())

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
	if c.CheckFrequencyMinutes < 1 {
		c.CheckFrequencyMinutes = 1
	}

	return nil
}

func getDefaultConfig() Config {
	cfg := Config{
		Debug:                 true,
		CheckFrequencyMinutes: 10,
	}

	cfg.Database = ConfigDatabase{
		Host:     "127.0.0.1",
		Username: "eqemu",
		Password: "eqemu",
		Database: "peq",
		Port:     3306,
	}
	return cfg
}
