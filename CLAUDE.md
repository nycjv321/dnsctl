# CLAUDE.md

This file provides context for Claude Code when working on this project.

## Project Overview

**dnsctl** is a macOS CLI tool with a TUI for switching between DNS server profiles. It wraps macOS `networksetup` commands with a user-friendly Bubble Tea interface.

## Tech Stack

- **Language**: Go 1.21+
- **TUI Framework**: Bubble Tea (github.com/charmbracelet/bubbletea)
- **Styling**: Lip Gloss (github.com/charmbracelet/lipgloss)
- **Config**: YAML via gopkg.in/yaml.v3
- **Platform**: macOS only (uses `networksetup` and `dscacheutil`)

## Architecture

```
cmd/dnsctl/main.go       # Entry point - loads config, creates DNS client, runs TUI
internal/
├── config/
│   ├── config.go        # Config loading from ~/.config/dnsctl/config.yaml
│   ├── config_darwin.go # macOS-specific defaults
│   ├── config_linux.go  # Linux-specific defaults
│   └── config_test.go   # Config tests
├── dns/
│   ├── client.go        # Client interface for DNS operations
│   ├── macos.go         # DNS operations wrapper around networksetup
│   └── mock.go          # Mock client for testing
└── tui/
    ├── app.go           # Bubble Tea Model with Init/Update/View
    ├── app_test.go      # TUI logic tests
    ├── keys.go          # KeyMap for all keybindings
    ├── styles.go        # Lip Gloss style definitions
    ├── views.go         # View rendering functions and View enum
    └── views_test.go    # View rendering tests
```

## Key Patterns

### Bubble Tea Model

The TUI follows the Elm architecture:
- `Model` struct holds all state (current view, selected index, DNS status, etc.)
- `Init()` returns initial command to fetch DNS status
- `Update()` handles messages and returns new model + commands
- `View()` renders current state to string

### View State Machine

```go
type View int
const (
    ViewMain View = iota      // Main dashboard
    ViewProfiles              // Profile selection list
    ViewServices              // Network service selection list
)
```

### Profile Configuration

Profiles are defined in `internal/config/config.go`:

```go
type Profile struct {
    Description string   `yaml:"description"`
    Servers     []string `yaml:"servers,omitempty"`
    DHCP        bool     `yaml:"dhcp,omitempty"`
}
```

- `IsDHCP()` returns true if `DHCP` is set or `Servers` is empty
- DHCP profiles clear DNS settings instead of setting specific servers
- Useful for "traveling" profiles where you want to use the network's DNS

### DNS Operations

All DNS operations are in `internal/dns/macos.go`:
- `ListNetworkServices()` - Get available network interfaces
- `GetDNSServers(service)` - Get current DNS for a service
- `SetDNSServers(service, servers)` - Set DNS servers
- `ClearDNSServers(service)` - Clear DNS (revert to DHCP)
- `FlushCache()` - Flush DNS cache

## Build Commands

```bash
make build      # Build to bin/dnsctl
make run        # Build and run
make install    # Install to /usr/local/bin (needs sudo)
make config     # Create config from example
make clean      # Remove build artifacts
make test       # Run tests
make fmt        # Format code
```

## Config Location

User config: `~/.config/dnsctl/config.yaml`

Default config is generated if file doesn't exist (see `config.DefaultConfig()`).

## Common Tasks

### Adding a new keybinding

1. Add to `KeyMap` struct in `internal/tui/keys.go`
2. Add to `DefaultKeyMap()` function
3. Handle in appropriate `handle*Keys()` function in `internal/tui/app.go`
4. Update help text in `internal/tui/views.go`

### Adding a new view

1. Add constant to `View` enum in `internal/tui/views.go`
2. Create `render*View()` function in `internal/tui/views.go`
3. Add case to `View()` method in `internal/tui/app.go`
4. Create `handle*Keys()` function in `internal/tui/app.go`
5. Add case to `handleKeyPress()` in `internal/tui/app.go`

### Adding a new DNS operation

1. Add method to `Client` in `internal/dns/macos.go`
2. Create command function in `internal/tui/app.go` that returns `tea.Cmd`
3. Handle result message in `Update()`

## Testing

### Running Tests

```bash
make test              # Run all tests
go test -v ./...       # Verbose output
go test -cover ./...   # With coverage report
go test -v ./internal/tui/...   # Run specific package

# View coverage in browser (opens interactive HTML report)
go test -coverprofile=coverage.out ./... && go tool cover -html=coverage.out
```

### Test Architecture

The project uses a mock DNS client for testing, avoiding system command dependencies:

```go
// internal/dns/mock.go
type MockClient struct {
    Services   []string              // Configurable service list
    DNSServers map[string][]string   // DNS servers by service
    SetError   error                 // Inject errors for testing
    ClearError error
    FlushError error
    SetCalls   []SetDNSCall          // Record calls for assertions
    ClearCalls []string
    FlushCalls int
}
```

### Test Files

| File | Coverage | Description |
|------|----------|-------------|
| `internal/config/config_test.go` | ~85% | Config loading, parsing, defaults, profile helpers |
| `internal/tui/app_test.go` | ~89% | Model init, Update(), key handlers, DNS operations |
| `internal/tui/views_test.go` | ~89% | View rendering, helper functions |

### Key Test Patterns

**Testing Update() with messages:**
```go
func TestUpdate_StatusMsg_Success(t *testing.T) {
    model, _ := testModel()
    msg := statusMsg{services: []string{"Wi-Fi"}, dnsServers: []string{"1.1.1.1"}}
    newModel, cmd := model.Update(msg)
    m := newModel.(Model)
    // Assert on m.services, m.currentDNS, etc.
}
```

**Testing key navigation:**
```go
func TestProfilesView_NavigateDown(t *testing.T) {
    model, _ := testModel()
    model.currentView = ViewProfiles
    msg := tea.KeyMsg{Type: tea.KeyDown}
    newModel, _ := model.Update(msg)
    // Assert selectedIndex changed
}
```

**Testing DNS operations with mock:**
```go
func TestApplyProfile_SetsDNS(t *testing.T) {
    model, mock := testModel()
    profile := config.Profile{Servers: []string{"9.9.9.9"}}
    cmd := model.applyProfile("test", profile)
    cmd()  // Execute the command
    // Assert mock.SetCalls contains expected call
}
```

### Writing New Tests

1. Use `testModel()` helper from `app_test.go` for TUI tests
2. Use `testConfig()` helper for consistent test configuration
3. Tests are in the same package (`tui`, `config`) for access to unexported fields
4. Use `t.TempDir()` for file-based config tests

## Permissions Note

Changing DNS requires elevated privileges. Users may need to run with `sudo` or grant terminal Full Disk Access.
