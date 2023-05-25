// Package ui handles the tui and application state
package ui

import (
	"text/template"
	"time"

	"github.com/argoproj/argo-cd/v2/pkg/apis/application/v1alpha1"
	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"

	"github.com/larntz/artui/ui/keys"
	"github.com/larntz/artui/ui/models"
	"github.com/larntz/artui/ui/templates"
)

// InitializeModel creates the initial model struct
func InitializeModel(cluster string, appEvent <-chan models.AppEvent, workerChan chan<- models.WorkerCmd, darkMode bool) models.ArTUIModel {
	appList := initAppList(v1alpha1.ApplicationList{})
	appList.Title = cluster
	textInput := initTextInput()
	templates := initTemplates()
	refreshDuration, err := time.ParseDuration("15s")
	if err != nil {
		panic(err)
	}

	return models.ArTUIModel{
		Cluster:         cluster,
		Ready:           false,
		Activity:        models.View,
		List:            appList,
		Applications:    v1alpha1.ApplicationList{},
		Textinput:       textInput,
		Templates:       templates,
		LastAppRefresh:  time.Now(),
		RefreshDuration: refreshDuration,
		AppEventChan:    appEvent,
		AppWorkerChan:   workerChan,
		DarkMode:        darkMode,
	}
}

func initTextInput() textinput.Model {
	ti := textinput.New()
	ti.Cursor.SetMode(cursor.CursorHide)
	ti.Prompt = " "
	ti.PromptStyle.PaddingLeft(0)
	ti.CharLimit = 20
	ti.Width = 20

	return ti
}

func initTemplates() *template.Template {
	var tpl *template.Template
	var err error
	for _, tplName := range templates.GetTemplateList() {
		tpl, err = template.New(tplName).Parse(templates.AppOverviewTemplate)
		if err != nil {
			panic(err)
		}
	}
	return tpl
}

func initAppList(apps v1alpha1.ApplicationList) list.Model {
	var appsListModel []list.Item
	for _, item := range apps.Items {
		appsListModel = append(appsListModel, models.AppListItem{
			Name:            item.Name,
			ItemDescription: string(item.Status.Health.Status) + "/" + string(item.Status.Sync.Status),
		})
	}
	appList := list.New(appsListModel, list.NewDefaultDelegate(), 20, 10)
	appList.KeyMap = keys.AppListKeyBinding
	appList.SetShowTitle(true)
	appList.SetShowPagination(true)
	appList.SetShowHelp(false)
	appList.SetShowFilter(true)
	appList.SetFilteringEnabled(true)

	return appList
}
