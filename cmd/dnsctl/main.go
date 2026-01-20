package main

import (
	"errors"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/nycjv321/dnsctl/internal/config"
	"github.com/nycjv321/dnsctl/internal/dns"
	"github.com/nycjv321/dnsctl/internal/tui"
)

func main() {
	// Load configuration
	cfg, err := config.Load("")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}

	// Create DNS client
	dnsClient, err := dns.NewClient()
	if err != nil {
		if errors.Is(err, dns.ErrNoDNSBackend) {
			fmt.Fprintln(os.Stderr, "Error: no supported DNS management system detected")
			fmt.Fprintln(os.Stderr, "")
			fmt.Fprintln(os.Stderr, "Supported systems:")
			fmt.Fprintln(os.Stderr, "  - macOS with networksetup")
			fmt.Fprintln(os.Stderr, "  - Linux with systemd-resolved (resolvectl)")
			fmt.Fprintln(os.Stderr, "  - Linux with NetworkManager (nmcli)")
		} else {
			fmt.Fprintf(os.Stderr, "Error creating DNS client: %v\n", err)
		}
		os.Exit(1)
	}

	// Create and run the TUI
	model := tui.NewModel(cfg, dnsClient)
	p := tea.NewProgram(model, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running program: %v\n", err)
		os.Exit(1)
	}
}
