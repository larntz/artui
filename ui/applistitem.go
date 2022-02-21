package ui

type appListItem struct {
	name        string
	description string
}

func (a appListItem) Title() string       { return a.name }
func (a appListItem) FilterValue() string { return a.name }
func (a appListItem) Description() string { return a.description }
