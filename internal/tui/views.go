package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/nycjv321/dnsctl/internal/config"
)

// View represents the current view state.
type View int

const (
	ViewMain View = iota
	ViewProfiles
	ViewServices
)

// renderMainView renders the main dashboard view.
func (m Model) renderMainView() string {
	var b strings.Builder

	// Title with current service indicator and inline icon
	title := titleStyle.Render("DNS Profile Switcher")
	icon := GetServiceIcon(m.currentService)
	serviceInfo := dimStyle.Render(fmt.Sprintf("Service: %s %s", m.currentService, icon))
	b.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, title, "    ", serviceInfo))
	b.WriteString("\n\n")

	// Current DNS servers
	b.WriteString("DNS: ")
	if len(m.currentDNS) == 0 {
		b.WriteString(dimStyle.Render("DHCP (automatic)"))
	} else {
		b.WriteString(normalStyle.Render(strings.Join(m.currentDNS, ", ")))
	}
	b.WriteString("\n")

	// Status message
	if m.statusMsg != "" {
		b.WriteString("\n")
		style := successStyle
		if m.statusIsError {
			style = errorStyle
		}
		b.WriteString(style.Render(m.statusMsg))
		b.WriteString("\n")
	}

	// Help
	b.WriteString("\n")
	b.WriteString(helpStyle.Render(m.renderMainHelp()))

	return b.String()
}

// renderMainHelp renders the help text for the main view.
func (m Model) renderMainHelp() string {
	return fmt.Sprintf(
		"%s switch profile  %s clear DNS  %s change service  %s refresh  %s quit",
		keyStyle.Render("[p]"),
		keyStyle.Render("[c]"),
		keyStyle.Render("[s]"),
		keyStyle.Render("[r]"),
		keyStyle.Render("[q]"),
	)
}

// renderProfilesView renders the profile selection view.
func (m Model) renderProfilesView() string {
	var b strings.Builder

	// Title with current service indicator and inline icon
	title := titleStyle.Render("Select DNS Profile")
	icon := GetServiceIcon(m.currentService)
	serviceInfo := dimStyle.Render(fmt.Sprintf("Service: %s %s", m.currentService, icon))
	b.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, title, "    ", serviceInfo))
	b.WriteString("\n\n")

	// Calculate max width for text wrapping
	profileNames := m.config.ProfileNames()
	const minListWidth = 40
	const maxListWidth = 50
	maxWidth := minListWidth

	// Check all profiles for max content width
	for _, name := range profileNames {
		profile := m.config.Profiles[name]
		// Check name width (with cursor prefix "  " or "> ")
		if w := len(name) + 2; w > maxWidth {
			maxWidth = w
		}
		// Check description width (with "    " indent)
		if w := len(profile.Description) + 4; w > maxWidth {
			maxWidth = w
		}
		// Check servers width (with "    Servers: " prefix)
		if !profile.IsDHCP() {
			serverStr := "Servers: " + strings.Join(profile.Servers, ", ")
			if w := len(serverStr) + 4; w > maxWidth {
				maxWidth = w
			}
		}
	}

	// Cap at reasonable max to force wrapping for very long lists
	if maxWidth > maxListWidth {
		maxWidth = maxListWidth
	}

	// Profile list
	for i, name := range profileNames {
		profile := m.config.Profiles[name]
		cursor := "  "
		style := normalStyle

		if i == m.selectedIndex {
			cursor = "> "
			style = selectedStyle
		}

		b.WriteString(cursor)
		b.WriteString(style.Render(name))
		b.WriteString("\n")

		// Show description and servers for selected item
		if i == m.selectedIndex {
			if profile.Description != "" {
				b.WriteString(fmt.Sprintf("    %s\n", descStyle.Render(profile.Description)))
			}
			if profile.IsDHCP() {
				b.WriteString(fmt.Sprintf("    %s\n", dimStyle.Render("DNS: DHCP (automatic)")))
			} else {
				serverStr := "Servers: " + strings.Join(profile.Servers, ", ")
				// Wrap server string if it exceeds available width
				availableWidth := maxWidth - 4 // account for indent
				wrappedServers := wrapText(serverStr, availableWidth)
				for _, line := range wrappedServers {
					b.WriteString(fmt.Sprintf("    %s\n", dimStyle.Render(line)))
				}
			}
		}
	}

	// Status message
	if m.statusMsg != "" {
		b.WriteString("\n")
		style := successStyle
		if m.statusIsError {
			style = errorStyle
		}
		b.WriteString(style.Render(m.statusMsg))
	}

	// Help
	b.WriteString("\n")
	b.WriteString(helpStyle.Render(m.renderListHelp()))

	return b.String()
}

// renderServicesView renders the network service selection view.
func (m Model) renderServicesView() string {
	var listBuilder strings.Builder

	// Title
	listBuilder.WriteString(titleStyle.Render("Select Network Service"))
	listBuilder.WriteString("\n\n")

	// Service list
	selectedService := ""
	for i, service := range m.services {
		cursor := "  "
		style := normalStyle

		if i == m.selectedIndex {
			cursor = "> "
			style = selectedStyle
			selectedService = service
		}

		// Mark current service
		suffix := ""
		if service == m.currentService {
			suffix = dimStyle.Render(" (current)")
		}

		listBuilder.WriteString(cursor)
		listBuilder.WriteString(style.Render(service))
		listBuilder.WriteString(suffix)
		listBuilder.WriteString("\n")
	}

	// Build the ASCII art panel for the selected service
	artPanel := ""
	if selectedService != "" {
		art := GetServiceArt(selectedService)
		artPanel = artStyle.Render(art)
	}

	// Join list and art horizontally
	content := lipgloss.JoinHorizontal(
		lipgloss.Center,
		listBuilder.String(),
		"    ", // spacing
		artPanel,
	)

	var b strings.Builder
	b.WriteString(content)

	// Status message
	if m.statusMsg != "" {
		b.WriteString("\n")
		style := successStyle
		if m.statusIsError {
			style = errorStyle
		}
		b.WriteString(style.Render(m.statusMsg))
	}

	// Help
	b.WriteString("\n")
	b.WriteString(helpStyle.Render(m.renderListHelp()))

	return b.String()
}

// renderListHelp renders the help text for list views.
func (m Model) renderListHelp() string {
	return fmt.Sprintf(
		"%s navigate  %s select  %s back  %s quit",
		keyStyle.Render("[↑/↓]"),
		keyStyle.Render("[enter]"),
		keyStyle.Render("[esc]"),
		keyStyle.Render("[q]"),
	)
}

// getSelectedProfile returns the currently selected profile name.
func (m Model) getSelectedProfile() (string, config.Profile, bool) {
	names := m.config.ProfileNames()
	if m.selectedIndex >= 0 && m.selectedIndex < len(names) {
		name := names[m.selectedIndex]
		profile, ok := m.config.GetProfile(name)
		return name, profile, ok
	}
	return "", config.Profile{}, false
}

// getSelectedService returns the currently selected service name.
func (m Model) getSelectedService() (string, bool) {
	if m.selectedIndex >= 0 && m.selectedIndex < len(m.services) {
		return m.services[m.selectedIndex], true
	}
	return "", false
}

// wrapText wraps text to fit within the specified width.
func wrapText(text string, width int) []string {
	if width <= 0 || len(text) <= width {
		return []string{text}
	}

	var lines []string
	remaining := text

	for len(remaining) > width {
		// Find a good break point (prefer space)
		breakPoint := width
		for i := width; i > 0; i-- {
			if remaining[i] == ' ' || remaining[i] == ',' {
				breakPoint = i + 1
				break
			}
		}
		// If no good break point found, force break at width
		if breakPoint == width && remaining[width] != ' ' && remaining[width] != ',' {
			breakPoint = width
		}

		lines = append(lines, strings.TrimRight(remaining[:breakPoint], " "))
		remaining = remaining[breakPoint:]
	}

	if len(remaining) > 0 {
		lines = append(lines, remaining)
	}

	return lines
}
