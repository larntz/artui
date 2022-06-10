package models

// AppListItem is used as bubbletea list.Items
type AppListItem struct {
	Name            string
	ItemDescription string
}

func (a AppListItem) Title() string       { return a.Name }
func (a AppListItem) FilterValue() string { return a.Name }
func (a AppListItem) Description() string { return a.ItemDescription }
