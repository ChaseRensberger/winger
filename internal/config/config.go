package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

const (
	DefaultRelayURL = "https://winger.dev"
	DirName         = ".winger"
	FileName        = "config.toml"
	PlanFileName    = ".plan"
)

type Config struct {
	RelayURL string `toml:"relay_url"`
	Handle   string `toml:"handle"`
}

func Dir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("could not find home directory: %w", err)
	}
	return filepath.Join(home, DirName), nil
}

func PlanPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("could not find home directory: %w", err)
	}
	return filepath.Join(home, PlanFileName), nil
}

func Path() (string, error) {
	dir, err := Dir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, FileName), nil
}

func Load() (*Config, error) {
	p, err := Path()
	if err != nil {
		return nil, err
	}

	var cfg Config
	if _, err := toml.DecodeFile(p, &cfg); err != nil {
		return nil, fmt.Errorf("could not read config: %w", err)
	}
	return &cfg, nil
}

func Save(cfg *Config) error {
	dir, err := Dir()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("could not create config directory: %w", err)
	}

	p := filepath.Join(dir, FileName)
	f, err := os.OpenFile(p, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("could not write config: %w", err)
	}
	defer f.Close()

	return toml.NewEncoder(f).Encode(cfg)
}
