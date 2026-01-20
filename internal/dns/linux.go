//go:build linux

package dns

import (
	"os/exec"
	"strings"
)

// NewClient creates a new DNS client for Linux.
// It auto-detects the available DNS management system.
// Priority: systemd-resolved > NetworkManager
func NewClient() (Client, error) {
	// Check for systemd-resolved
	if isResolvedActive() {
		return &resolvedClient{}, nil
	}

	// Check for NetworkManager
	if isNetworkManagerActive() {
		return &nmClient{}, nil
	}

	return nil, ErrNoDNSBackend
}

// isResolvedActive checks if systemd-resolved is active.
func isResolvedActive() bool {
	// Check if resolvectl exists
	if _, err := exec.LookPath("resolvectl"); err != nil {
		return false
	}

	// Check if systemd-resolved service is active
	cmd := exec.Command("systemctl", "is-active", "systemd-resolved")
	output, err := cmd.Output()
	if err != nil {
		return false
	}

	return strings.TrimSpace(string(output)) == "active"
}

// isNetworkManagerActive checks if NetworkManager is active.
func isNetworkManagerActive() bool {
	// Check if nmcli exists
	if _, err := exec.LookPath("nmcli"); err != nil {
		return false
	}

	// Check if NetworkManager service is active
	cmd := exec.Command("systemctl", "is-active", "NetworkManager")
	output, err := cmd.Output()
	if err != nil {
		return false
	}

	return strings.TrimSpace(string(output)) == "active"
}
