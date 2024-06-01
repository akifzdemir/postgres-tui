package constants

import "github.com/charmbracelet/bubbles/key"

type GeneralKeyMap struct {
	Enter key.Binding
	Back  key.Binding
	Quit  key.Binding
	Up    key.Binding
	Down  key.Binding
}

type TableKeyMap struct {
	GeneralKeyMap
	Create key.Binding
	Delete key.Binding
	Rename key.Binding
}

type LoginKeyMap struct {
	GeneralKeyMap
	Submit key.Binding
}

var GeneralKeys = GeneralKeyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "move down"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "select"),
	),
	Back: key.NewBinding(
		key.WithKeys("backspace"),
		key.WithHelp("backspace", "back"),
	),

	Quit: key.NewBinding(
		key.WithKeys("ctrl+c", "q"),
		key.WithHelp("ctrl+c/q", "quit"),
	),
}

var TableKeys = TableKeyMap{
	GeneralKeyMap: GeneralKeys,
	Create: key.NewBinding(
		key.WithKeys("c"),
		key.WithHelp("c", "create"),
	),
	Rename: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "rename"),
	),
	Delete: key.NewBinding(
		key.WithKeys("d"),
		key.WithHelp("d", "delete"),
	),
}

var LoginKeys = LoginKeyMap{
	GeneralKeyMap: GeneralKeys,
	Submit: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "Submit"),
	),
}

func (k GeneralKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.Enter, k.Quit, k.Back}
}

func (k GeneralKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Enter, k.Quit},
	}
}

func (k TableKeyMap) ShortHelp() []key.Binding {
	return append(k.GeneralKeyMap.ShortHelp(), k.Create, k.Delete, k.Up)
}

func (k TableKeyMap) FullHelp() [][]key.Binding {
	return append(k.GeneralKeyMap.FullHelp(), []key.Binding{k.Create, k.Delete, k.Rename})
}

func (k LoginKeyMap) ShortHelp() []key.Binding {
	return append(k.GeneralKeyMap.ShortHelp(), k.Submit, k.Up, k.Down)
}
