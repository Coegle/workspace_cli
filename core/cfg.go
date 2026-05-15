package core

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

const (
	ConfigPath      = ".ws"
	ConfigFile      = "config"
	ConfigType      = "yaml"
	EnvPrefix       = "WS"
	UpdateCacheFile = "update_cache.json"
)

type Symlink struct {
	Source      string `mapstructure:"src"`
	Destination string `mapstructure:"dest"`
}

type Config struct {
	BaseReposPath     string    `mapstructure:"base_repos"`
	BaseWorkspacePath string    `mapstructure:"base_ws"`
	ReplaceSlashInDir bool      `mapstructure:"replace_slash"`
	Symlinks          []Symlink `mapstructure:"symlinks"`
}

func GetConfig() (*Config, error) {
	var cfg Config

	// Set defaults
	home, err := os.UserHomeDir()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Warning: failed to get user home directory: %v\n", err)
		home = "." // fallback to current directory
	}
	viper.SetDefault("base_repos", filepath.Join(home, "repos"))
	viper.SetDefault("base_ws", filepath.Join(home, "workspace"))
	viper.SetDefault("replace_slash", true)

	if err = viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("error parsing config: %w", err)
	}
	return &cfg, nil
}

func (c *Config) GetReposPath() string {
	return c.BaseReposPath
}

func (c *Config) GetWorkspacePath() string {
	return c.BaseWorkspacePath
}

func (c *Config) GetSafeDirName(branchName string) string {
	if c.ReplaceSlashInDir {
		if !strings.Contains(branchName, "/") && !strings.Contains(branchName, "\\") && branchName != "" && branchName != "." {
			return branchName
		}
		return strings.ReplaceAll(branchName, "/", "_")
	}
	return branchName
}
