// internal/config.go
package internal

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type GitlabConfig struct {
	BaseURL   string `json:"base_url"`
	Token     string `json:"token"`
	UserID    string `json:"user_id"`
	ProjectID string `json:"project_id"`
	GroupID   string `json:"group_id"`
}

type Config struct {
	Model  string       `json:"model"`
	Gitlab GitlabConfig `json:"gitlab"`
	// Futuramente: Github GithubConfig `json:"github"`
}

func configPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".issue_writer_cli_config.json")
}

func SaveConfig(cfg Config) error {
	f, err := os.Create(configPath())
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(cfg)
}

func LoadConfig() (Config, error) {
	var cfg Config
	f, err := os.Open(configPath())
	if err != nil {
		return cfg, err
	}
	defer f.Close()
	err = json.NewDecoder(f).Decode(&cfg)

	if cfg.Model == "" {
		cfg.Model = "gpt-4-turbo"
	}

	return cfg, err
}
