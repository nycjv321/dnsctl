package tui

import (
	"fmt"
	"strings"

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

	// Title
	b.WriteString(titleStyle.Render("DNS Profile Switcher"))
	b.WriteString("\n\n")

	// Current service
	b.WriteString(fmt.Sprintf("Service: %s\n", selectedStyle.Render(m.currentService)))

	// Current DNS servers
	b.WriteString("DNS:     ")
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

	// Title
	b.WriteString(titleStyle.Render("Select DNS Profile"))
	b.WriteString("\n\n")

	// Profile list
	profileNames := m.config.ProfileNames()
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
				b.WriteString(fmt.Sprintf("    %s\n", dimStyle.Render("Servers: "+strings.Join(profile.Servers, ", "))))
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
	var b strings.Builder

	// Title
	b.WriteString(titleStyle.Render("Select Network Service"))
	b.WriteString("\n\n")

	// Service list
	for i, service := range m.services {
		cursor := "  "
		style := normalStyle

		if i == m.selectedIndex {
			cursor = "> "
			style = selectedStyle
		}

		// Mark current service
		suffix := ""
		if service == m.currentService {
			suffix = dimStyle.Render(" (current)")
		}

		b.WriteString(cursor)
		b.WriteString(style.Render(service))
		b.WriteString(suffix)
		b.WriteString("\n")
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
