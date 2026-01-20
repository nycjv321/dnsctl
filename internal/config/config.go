package config

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"gopkg.in/yaml.v3"
)

// Profile represents a DNS profile with a name and servers.
type Profile struct {
	Description string   `yaml:"description"`
	Servers     []string `yaml:"servers,omitempty"`
	DHCP        bool     `yaml:"dhcp,omitempty"`
}

// IsDHCP returns true if this profile clears DNS to use DHCP.
func (p Profile) IsDHCP() bool {
	return p.DHCP || len(p.Servers) == 0
}

// Settings contains application settings.
type Settings struct {
	FlushCache bool `yaml:"flush_cache"`
}

// Config represents the application configuration.
type Config struct {
	Version        int                `yaml:"version"`
	DefaultService string             `yaml:"default_service"`
	Profiles       map[string]Profile `yaml:"profiles"`
	Settings       Settings           `yaml:"settings"`
}

// DefaultConfigPath returns the default configuration file path.
func DefaultConfigPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".config", "dnsctl", "config.yaml")
}

// Load loads the configuration from the specified path.
func Load(path string) (*Config, error) {
	if path == "" {
		path = DefaultConfigPath()
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return DefaultConfig(), nil
		}
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Apply defaults
	if cfg.DefaultService == "" {
		cfg.DefaultService = defaultServiceName()
	}
	if cfg.Profiles == nil {
		cfg.Profiles = make(map[string]Profile)
	}

	return &cfg, nil
}

// DefaultConfig returns a default configuration.
func DefaultConfig() *Config {
	return &Config{
		Version:        1,
		DefaultService: defaultServiceName(),
		Profiles: map[string]Profile{
			"cloudflare": {
				Description: "Cloudflare DNS",
				Servers:     []string{"1.1.1.1", "1.0.0.1"},
			},
			"google": {
				Description: "Google Public DNS",
				Servers:     []string{"8.8.8.8", "8.8.4.4"},
			},
		},
		Settings: Settings{
			FlushCache: true,
		},
	}
}

// Save saves the configuration to the specified path.
func (c *Config) Save(path string) error {
	if path == "" {
		path = DefaultConfigPath()
	}

	// Ensure directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	data, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// ProfileNames returns a sorted list of profile names.
func (c *Config) ProfileNames() []string {
	names := make([]string, 0, len(c.Profiles))
	for name := range c.Profiles {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// GetProfile returns a profile by name.
func (c *Config) GetProfile(name string) (Profile, bool) {
	p, ok := c.Profiles[name]
	return p, ok
}
