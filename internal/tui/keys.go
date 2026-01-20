package tui

import "github.com/charmbracelet/bubbles/key"

// KeyMap defines the keybindings for the application.
type KeyMap struct {
	Up            key.Binding
	Down          key.Binding
	Select        key.Binding
	Back          key.Binding
	Quit          key.Binding
	SwitchProfile key.Binding
	ClearDNS      key.Binding
	ChangeService key.Binding
	Refresh       key.Binding
}

// DefaultKeyMap returns the default keybindings.
func DefaultKeyMap() KeyMap {
	return KeyMap{
		Up: key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("↑/k", "up"),
		),
		Down: key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("↓/j", "down"),
		),
		Select: key.NewBinding(
			key.WithKeys("enter", " "),
			key.WithHelp("enter", "select"),
		),
		Back: key.NewBinding(
			key.WithKeys("esc", "backspace"),
			key.WithHelp("esc", "back"),
		),
		Quit: key.NewBinding(
			key.WithKeys("q", "ctrl+c"),
			key.WithHelp("q", "quit"),
		),
		SwitchProfile: key.NewBinding(
			key.WithKeys("p"),
			key.WithHelp("p", "switch profile"),
		),
		ClearDNS: key.NewBinding(
			key.WithKeys("c"),
			key.WithHelp("c", "clear DNS"),
		),
		ChangeService: key.NewBinding(
			key.WithKeys("s"),
			key.WithHelp("s", "change service"),
		),
		Refresh: key.NewBinding(
			key.WithKeys("r"),
			key.WithHelp("r", "refresh"),
		),
	}
}
