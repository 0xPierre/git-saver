package config

import (
	"fmt"
	"os"

	"go.yaml.in/yaml/v3"
)

type Source struct {
	Type  string `yaml:"type"`
	URL   string `yaml:"url"`
	Token string `yaml:"token"`
}

type Target struct {
	Type  string `yaml:"type"`
	URL   string `yaml:"url"`
	Token string `yaml:"token"`
	Dir   string `yaml:"dir"`
}

type Config struct {
	Sources              []Source `yaml:"sources"`
	SyncAllRepos         bool     `yaml:"sync-all-repo"`
	SourceRepositories   []string `yaml:"source-repositories"`
	Targets              []Target `yaml:"targets"`
	DefaultPullDirectory string   `yaml:"default-pull-directory"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config: %w", err)
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing config: %w", err)
	}
	return &cfg, nil
}
