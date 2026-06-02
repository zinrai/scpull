package main

import (
	"fmt"
	"io"
	"os"

	"github.com/goccy/go-yaml"
)

// Profile defines what to fetch from a host: the login user and the
// list of absolute remote paths. The destination is never part of a
// profile; it is always derived from the host argument at runtime.
type Profile struct {
	User  string   `yaml:"user"`
	Paths []string `yaml:"paths"`
}

// Config is the top-level structure of the YAML configuration file.
type Config struct {
	Profiles map[string]Profile `yaml:"profiles"`
}

// loadConfig reads and parses the configuration file at path.
func loadConfig(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open config: %w", err)
	}
	defer f.Close()

	return parseConfig(f)
}

// parseConfig decodes a Config from r. It is separated from loadConfig
// so it can be tested without touching the filesystem.
func parseConfig(r io.Reader) (*Config, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	if len(cfg.Profiles) == 0 {
		return nil, fmt.Errorf("no profiles defined")
	}

	return &cfg, nil
}

// profile returns the named profile after validating it. An unknown
// name, an empty user, or an empty path list are all errors.
func (c *Config) profile(name string) (Profile, error) {
	p, ok := c.Profiles[name]
	if !ok {
		return Profile{}, fmt.Errorf("profile %q not found", name)
	}

	if p.User == "" {
		return Profile{}, fmt.Errorf("profile %q: user is empty", name)
	}

	if len(p.Paths) == 0 {
		return Profile{}, fmt.Errorf("profile %q: paths is empty", name)
	}

	for _, path := range p.Paths {
		if path == "" {
			return Profile{}, fmt.Errorf("profile %q: contains an empty path", name)
		}
	}

	return p, nil
}
