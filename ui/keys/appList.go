package keys

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
)

// AppListKeyBinding are the keys bound to the appList
// tui widget
var AppListKeyBinding = list.KeyMap{
	// Keybindings used when browsing the list.
	CursorUp:    key.NewBinding(key.WithKeys("k")),
	CursorDown:  key.NewBinding(key.WithKeys("j")),
	NextPage:    key.NewBinding(key.WithDisabled()),
	PrevPage:    key.NewBinding(key.WithDisabled()),
	GoToStart:   key.NewBinding(key.WithDisabled()),
	GoToEnd:     key.NewBinding(key.WithDisabled()),
	Filter:      key.NewBinding(key.WithKeys("/")),
	ClearFilter: key.NewBinding(key.WithKeys("esc")),

	// Keybindings used when setting a filter.
	CancelWhileFiltering: key.NewBinding(key.WithKeys("esc")),
	AcceptWhileFiltering: key.NewBinding(key.WithKeys("enter")),

	// Help toggle keybindings.
	ShowFullHelp:  key.NewBinding(key.WithDisabled()),
	CloseFullHelp: key.NewBinding(key.WithDisabled()),

	// The quit keybinding. This won't be caught when filtering.
	Quit: key.NewBinding(key.WithKeys("ctrl+c")),

	// The quit-no-matter-what keybinding. This will be caught when filtering.
	ForceQuit: key.NewBinding(key.WithKeys("ctrl+c")),
}
