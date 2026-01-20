package config

import (
	"os"
	"path/filepath"
	"testing"
)

// TestDefaultConfig_HasRequiredProfiles tests that default config has expected profiles.
func TestDefaultConfig_HasRequiredProfiles(t *testing.T) {
	cfg := DefaultConfig()

	if cfg.Profiles == nil {
		t.Fatal("expected profiles to be non-nil")
	}

	if _, ok := cfg.Profiles["cloudflare"]; !ok {
		t.Error("expected cloudflare profile")
	}
	if _, ok := cfg.Profiles["google"]; !ok {
		t.Error("expected google profile")
	}
}

// TestDefaultConfig_HasVersion tests that default config has version set.
func TestDefaultConfig_HasVersion(t *testing.T) {
	cfg := DefaultConfig()

	if cfg.Version != 1 {
		t.Errorf("expected version 1, got %d", cfg.Version)
	}
}

// TestDefaultConfig_HasDefaultService tests that default config has default service.
func TestDefaultConfig_HasDefaultService(t *testing.T) {
	cfg := DefaultConfig()

	if cfg.DefaultService == "" {
		t.Error("expected DefaultService to be set")
	}
}

// TestDefaultConfig_HasSettings tests that default config has settings.
func TestDefaultConfig_HasSettings(t *testing.T) {
	cfg := DefaultConfig()

	if !cfg.Settings.FlushCache {
		t.Error("expected FlushCache to be true by default")
	}
}

// TestProfileNames_SortedAlphabetically tests that profile names are sorted.
func TestProfileNames_SortedAlphabetically(t *testing.T) {
	cfg := &Config{
		Profiles: map[string]Profile{
			"zebra":  {Description: "Z"},
			"alpha":  {Description: "A"},
			"middle": {Description: "M"},
			"bravo":  {Description: "B"},
		},
	}

	names := cfg.ProfileNames()

	expected := []string{"alpha", "bravo", "middle", "zebra"}
	if len(names) != len(expected) {
		t.Fatalf("expected %d names, got %d", len(expected), len(names))
	}
	for i, name := range names {
		if name != expected[i] {
			t.Errorf("expected names[%d] to be %s, got %s", i, expected[i], name)
		}
	}
}

// TestProfileNames_EmptyProfiles tests ProfileNames with no profiles.
func TestProfileNames_EmptyProfiles(t *testing.T) {
	cfg := &Config{
		Profiles: map[string]Profile{},
	}

	names := cfg.ProfileNames()

	if len(names) != 0 {
		t.Errorf("expected 0 names, got %d", len(names))
	}
}

// TestGetProfile_Found tests getting an existing profile.
func TestGetProfile_Found(t *testing.T) {
	cfg := &Config{
		Profiles: map[string]Profile{
			"test": {
				Description: "Test Profile",
				Servers:     []string{"9.9.9.9"},
			},
		},
	}

	profile, ok := cfg.GetProfile("test")

	if !ok {
		t.Fatal("expected profile to be found")
	}
	if profile.Description != "Test Profile" {
		t.Errorf("expected 'Test Profile', got '%s'", profile.Description)
	}
	if len(profile.Servers) != 1 || profile.Servers[0] != "9.9.9.9" {
		t.Errorf("expected [9.9.9.9], got %v", profile.Servers)
	}
}

// TestGetProfile_NotFound tests getting a non-existent profile.
func TestGetProfile_NotFound(t *testing.T) {
	cfg := &Config{
		Profiles: map[string]Profile{},
	}

	_, ok := cfg.GetProfile("nonexistent")

	if ok {
		t.Error("expected profile to not be found")
	}
}

// TestProfile_IsDHCP_WhenFlagSet tests IsDHCP returns true when flag is set.
func TestProfile_IsDHCP_WhenFlagSet(t *testing.T) {
	profile := Profile{
		Description: "DHCP Profile",
		DHCP:        true,
		Servers:     []string{"1.1.1.1"}, // Even with servers, DHCP flag takes precedence
	}

	if !profile.IsDHCP() {
		t.Error("expected IsDHCP to be true when DHCP flag is set")
	}
}

// TestProfile_IsDHCP_WhenNoServers tests IsDHCP returns true when no servers.
func TestProfile_IsDHCP_WhenNoServers(t *testing.T) {
	profile := Profile{
		Description: "Empty Servers Profile",
		DHCP:        false,
		Servers:     nil,
	}

	if !profile.IsDHCP() {
		t.Error("expected IsDHCP to be true when servers is nil")
	}
}

// TestProfile_IsDHCP_WhenEmptyServers tests IsDHCP returns true when servers empty.
func TestProfile_IsDHCP_WhenEmptyServers(t *testing.T) {
	profile := Profile{
		Description: "Empty Servers Profile",
		DHCP:        false,
		Servers:     []string{},
	}

	if !profile.IsDHCP() {
		t.Error("expected IsDHCP to be true when servers is empty")
	}
}

// TestProfile_NotDHCP_WithServers tests IsDHCP returns false when servers set.
func TestProfile_NotDHCP_WithServers(t *testing.T) {
	profile := Profile{
		Description: "Normal Profile",
		DHCP:        false,
		Servers:     []string{"1.1.1.1", "1.0.0.1"},
	}

	if profile.IsDHCP() {
		t.Error("expected IsDHCP to be false when servers are set")
	}
}

// TestLoad_DefaultsWhenMissing tests that Load returns defaults when file missing.
func TestLoad_DefaultsWhenMissing(t *testing.T) {
	// Use a path that definitely doesn't exist
	cfg, err := Load("/nonexistent/path/to/config.yaml")

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if cfg == nil {
		t.Fatal("expected config to be non-nil")
	}
	if cfg.Version != 1 {
		t.Errorf("expected version 1, got %d", cfg.Version)
	}
	if len(cfg.Profiles) == 0 {
		t.Error("expected default profiles")
	}
}

// TestLoad_ParsesYAML tests that Load correctly parses YAML.
func TestLoad_ParsesYAML(t *testing.T) {
	// Create a temp directory and config file
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	yamlContent := `version: 2
default_service: Ethernet
profiles:
  custom:
    description: Custom DNS
    servers:
      - 9.9.9.9
      - 149.112.112.112
settings:
  flush_cache: false
`
	if err := os.WriteFile(configPath, []byte(yamlContent), 0644); err != nil {
		t.Fatalf("failed to write test config: %v", err)
	}

	cfg, err := Load(configPath)

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if cfg.Version != 2 {
		t.Errorf("expected version 2, got %d", cfg.Version)
	}
	if cfg.DefaultService != "Ethernet" {
		t.Errorf("expected Ethernet, got %s", cfg.DefaultService)
	}
	if _, ok := cfg.Profiles["custom"]; !ok {
		t.Error("expected custom profile")
	}
	if cfg.Settings.FlushCache {
		t.Error("expected FlushCache to be false")
	}
}

// TestLoad_AppliesDefaults tests that Load applies defaults for missing fields.
func TestLoad_AppliesDefaults(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	// Minimal config without default_service and profiles
	yamlContent := `version: 1
`
	if err := os.WriteFile(configPath, []byte(yamlContent), 0644); err != nil {
		t.Fatalf("failed to write test config: %v", err)
	}

	cfg, err := Load(configPath)

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if cfg.DefaultService == "" {
		t.Error("expected DefaultService to be set to default")
	}
	if cfg.Profiles == nil {
		t.Error("expected Profiles to be initialized")
	}
}

// TestLoad_InvalidYAML tests that Load returns error for invalid YAML.
func TestLoad_InvalidYAML(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	invalidYAML := `version: invalid yaml: ::::`
	if err := os.WriteFile(configPath, []byte(invalidYAML), 0644); err != nil {
		t.Fatalf("failed to write test config: %v", err)
	}

	_, err := Load(configPath)

	if err == nil {
		t.Error("expected error for invalid YAML")
	}
}

// TestSave_CreatesFile tests that Save creates a config file.
func TestSave_CreatesFile(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "subdir", "config.yaml")

	cfg := DefaultConfig()

	err := cfg.Save(configPath)

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Error("expected config file to be created")
	}

	// Verify we can load it back
	loaded, err := Load(configPath)
	if err != nil {
		t.Fatalf("failed to load saved config: %v", err)
	}
	if loaded.Version != cfg.Version {
		t.Errorf("expected version %d, got %d", cfg.Version, loaded.Version)
	}
}

// TestSave_OverwritesExisting tests that Save overwrites existing file.
func TestSave_OverwritesExisting(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	// Create initial file
	initial := &Config{
		Version: 1,
		Profiles: map[string]Profile{
			"initial": {Description: "Initial"},
		},
	}
	if err := initial.Save(configPath); err != nil {
		t.Fatalf("failed to save initial: %v", err)
	}

	// Overwrite with new config
	updated := &Config{
		Version: 2,
		Profiles: map[string]Profile{
			"updated": {Description: "Updated"},
		},
	}
	if err := updated.Save(configPath); err != nil {
		t.Fatalf("failed to save updated: %v", err)
	}

	// Load and verify
	loaded, err := Load(configPath)
	if err != nil {
		t.Fatalf("failed to load: %v", err)
	}
	if loaded.Version != 2 {
		t.Errorf("expected version 2, got %d", loaded.Version)
	}
	if _, ok := loaded.Profiles["updated"]; !ok {
		t.Error("expected updated profile")
	}
	if _, ok := loaded.Profiles["initial"]; ok {
		t.Error("did not expect initial profile")
	}
}

// TestDefaultConfigPath tests that DefaultConfigPath returns expected path.
func TestDefaultConfigPath(t *testing.T) {
	path := DefaultConfigPath()

	if path == "" {
		t.Skip("skipping test when home directory is not available")
	}

	if !filepath.IsAbs(path) {
		t.Errorf("expected absolute path, got %s", path)
	}
	if !filepath.IsAbs(path) {
		t.Errorf("expected absolute path, got: %s", path)
	}
	if filepath.Base(path) != "config.yaml" {
		t.Errorf("expected config.yaml, got %s", filepath.Base(path))
	}
}

// TestLoad_EmptyPath_UsesDefault tests that empty path uses default.
func TestLoad_EmptyPath_UsesDefault(t *testing.T) {
	// This test may fail if the default config doesn't exist,
	// but should return defaults in that case
	cfg, err := Load("")

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if cfg == nil {
		t.Fatal("expected config to be non-nil")
	}
}

// TestSave_EmptyPath_UsesDefault tests that empty path uses default.
func TestSave_EmptyPath_UsesDefault(t *testing.T) {
	// Skip this test as it would modify the real default config
	t.Skip("skipping to avoid modifying real config")
}

// TestCloudflareProfile_HasCorrectServers tests cloudflare profile DNS servers.
func TestCloudflareProfile_HasCorrectServers(t *testing.T) {
	cfg := DefaultConfig()

	profile, ok := cfg.GetProfile("cloudflare")
	if !ok {
		t.Fatal("expected cloudflare profile")
	}

	if len(profile.Servers) != 2 {
		t.Fatalf("expected 2 servers, got %d", len(profile.Servers))
	}
	if profile.Servers[0] != "1.1.1.1" {
		t.Errorf("expected 1.1.1.1, got %s", profile.Servers[0])
	}
	if profile.Servers[1] != "1.0.0.1" {
		t.Errorf("expected 1.0.0.1, got %s", profile.Servers[1])
	}
}

// TestGoogleProfile_HasCorrectServers tests google profile DNS servers.
func TestGoogleProfile_HasCorrectServers(t *testing.T) {
	cfg := DefaultConfig()

	profile, ok := cfg.GetProfile("google")
	if !ok {
		t.Fatal("expected google profile")
	}

	if len(profile.Servers) != 2 {
		t.Fatalf("expected 2 servers, got %d", len(profile.Servers))
	}
	if profile.Servers[0] != "8.8.8.8" {
		t.Errorf("expected 8.8.8.8, got %s", profile.Servers[0])
	}
	if profile.Servers[1] != "8.8.4.4" {
		t.Errorf("expected 8.8.4.4, got %s", profile.Servers[1])
	}
}
