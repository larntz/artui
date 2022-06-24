package models

import (
	"log"

	"github.com/charmbracelet/bubbles/list"
	"github.com/larntz/artui/ui/keys"
)

// AppListItem is used as bubbletea list.Items
type AppListItem struct {
	Name            string
	ItemDescription string
}

// Title returns the item's name
func (a AppListItem) Title() string { return a.Name }

// FilterValue returns the string to filter on
func (a AppListItem) FilterValue() string { return a.Name }

// Description returns the ItemDescription
func (a AppListItem) Description() string { return a.ItemDescription }

// Return a new appList full of AppListItems
func (m ArTUIModel) updateAppList() list.Model {
	log.Printf("updateAppList: got %d apps", len(m.Applications.Items))
	var appListItems []list.Item
	for _, app := range m.Applications.Items {
		description := string(app.Status.Health.Status) + "/" + string(app.Status.Sync.Status)
		appListItems = append(appListItems, AppListItem{
			Name:            app.Name,
			ItemDescription: description,
		})
	}

	appList := list.New(appListItems, list.NewDefaultDelegate(), 0, 0)
	appList.Title = "App List"
	appList.KeyMap = keys.AppListKeyBinding
	appList.SetShowTitle(true)
	appList.SetShowPagination(true)
	appList.SetShowHelp(true)
	appList.SetShowFilter(true)
	appList.SetFilteringEnabled(true)
	appList.SetSize(int(float32(m.WindowWidth)*0.25), getContentHeight(m))

	return appList
}
