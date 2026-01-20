//go:build linux

package dns

import (
	"bufio"
	"fmt"
	"os/exec"
	"strings"
)

// resolvedClient provides DNS management via systemd-resolved.
type resolvedClient struct{}

// Name returns the backend name for display purposes.
func (c *resolvedClient) Name() string {
	return "systemd-resolved"
}

// ListNetworkServices returns all available network interfaces.
func (c *resolvedClient) ListNetworkServices() ([]string, error) {
	cmd := exec.Command("resolvectl", "status")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to list network services: %w", err)
	}

	var interfaces []string
	scanner := bufio.NewScanner(strings.NewReader(string(output)))

	for scanner.Scan() {
		line := scanner.Text()
		// Lines like "Link 2 (eth0)" or "Link 3 (wlan0)"
		if strings.HasPrefix(line, "Link ") {
			// Extract interface name from parentheses
			start := strings.Index(line, "(")
			end := strings.Index(line, ")")
			if start != -1 && end != -1 && end > start {
				iface := line[start+1 : end]
				interfaces = append(interfaces, iface)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to parse resolvectl output: %w", err)
	}

	return interfaces, nil
}

// GetDNSServers returns the current DNS servers for an interface.
func (c *resolvedClient) GetDNSServers(service string) ([]string, error) {
	cmd := exec.Command("resolvectl", "dns", service)
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get DNS servers for %s: %w", service, err)
	}

	text := strings.TrimSpace(string(output))
	if text == "" {
		return nil, nil
	}

	// Output format: "Link 2 (eth0): 8.8.8.8 8.8.4.4"
	// Find the colon and parse servers after it
	colonIdx := strings.Index(text, ":")
	if colonIdx == -1 {
		return nil, nil
	}

	serverPart := strings.TrimSpace(text[colonIdx+1:])
	if serverPart == "" {
		return nil, nil
	}

	servers := strings.Fields(serverPart)
	return servers, nil
}

// SetDNSServers sets the DNS servers for an interface.
func (c *resolvedClient) SetDNSServers(service string, servers []string) error {
	args := []string{"dns", service}
	args = append(args, servers...)

	cmd := exec.Command("resolvectl", args...)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to set DNS servers: %s: %w", string(output), err)
	}

	return nil
}

// ClearDNSServers clears DNS servers, reverting to defaults.
func (c *resolvedClient) ClearDNSServers(service string) error {
	cmd := exec.Command("resolvectl", "revert", service)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to clear DNS servers: %s: %w", string(output), err)
	}

	return nil
}

// FlushCache flushes the DNS cache.
func (c *resolvedClient) FlushCache() error {
	cmd := exec.Command("resolvectl", "flush-caches")
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to flush DNS cache: %s: %w", string(output), err)
	}

	return nil
}
