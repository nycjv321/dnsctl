//go:build linux

package dns

import (
	"fmt"
	"os/exec"
	"strings"
)

// nmClient provides DNS management via NetworkManager.
type nmClient struct{}

// Name returns the backend name for display purposes.
func (c *nmClient) Name() string {
	return "NetworkManager"
}

// ListNetworkServices returns all active network connections.
func (c *nmClient) ListNetworkServices() ([]string, error) {
	// Get active connections in terse format
	cmd := exec.Command("nmcli", "-t", "-f", "NAME", "connection", "show", "--active")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to list network services: %w", err)
	}

	text := strings.TrimSpace(string(output))
	if text == "" {
		return nil, nil
	}

	lines := strings.Split(text, "\n")
	var connections []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			connections = append(connections, line)
		}
	}

	return connections, nil
}

// GetDNSServers returns the current DNS servers for a connection.
func (c *nmClient) GetDNSServers(service string) ([]string, error) {
	cmd := exec.Command("nmcli", "-t", "-f", "ipv4.dns", "connection", "show", service)
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get DNS servers for %s: %w", service, err)
	}

	text := strings.TrimSpace(string(output))
	if text == "" {
		return nil, nil
	}

	// Output format: "ipv4.dns:8.8.8.8,8.8.4.4"
	colonIdx := strings.Index(text, ":")
	if colonIdx == -1 {
		return nil, nil
	}

	serverPart := strings.TrimSpace(text[colonIdx+1:])
	if serverPart == "" {
		return nil, nil
	}

	// DNS servers are comma-separated
	servers := strings.Split(serverPart, ",")
	var result []string
	for _, s := range servers {
		s = strings.TrimSpace(s)
		if s != "" {
			result = append(result, s)
		}
	}

	return result, nil
}

// SetDNSServers sets the DNS servers for a connection.
func (c *nmClient) SetDNSServers(service string, servers []string) error {
	dnsValue := strings.Join(servers, ",")

	// Modify the connection
	cmd := exec.Command("nmcli", "connection", "modify", service,
		"ipv4.dns", dnsValue,
		"ipv4.ignore-auto-dns", "yes")
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to set DNS servers: %s: %w", string(output), err)
	}

	// Reactivate the connection to apply changes
	cmd = exec.Command("nmcli", "connection", "up", service)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to reactivate connection: %s: %w", string(output), err)
	}

	return nil
}

// ClearDNSServers clears DNS servers, reverting to DHCP defaults.
func (c *nmClient) ClearDNSServers(service string) error {
	// Clear manual DNS and enable auto DNS
	cmd := exec.Command("nmcli", "connection", "modify", service,
		"ipv4.dns", "",
		"ipv4.ignore-auto-dns", "no")
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to clear DNS servers: %s: %w", string(output), err)
	}

	// Reactivate the connection to apply changes
	cmd = exec.Command("nmcli", "connection", "up", service)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to reactivate connection: %s: %w", string(output), err)
	}

	return nil
}

// FlushCache flushes the DNS cache.
func (c *nmClient) FlushCache() error {
	// Try resolvectl first (if systemd-resolved is being used as a cache)
	cmd := exec.Command("resolvectl", "flush-caches")
	if err := cmd.Run(); err == nil {
		return nil
	}

	// Try nscd if available
	cmd = exec.Command("nscd", "-i", "hosts")
	if err := cmd.Run(); err == nil {
		return nil
	}

	// If neither works, return success anyway since NetworkManager
	// may not have a cache to flush
	return nil
}
