package ui

import (
	"text/template"

	"github.com/argoproj/argo-cd/v2/pkg/apiclient"
	"github.com/argoproj/argo-cd/v2/pkg/apiclient/session"
	"github.com/argoproj/argo-cd/v2/pkg/apis/application/v1alpha1"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"

	"github.com/larntz/artui/argo"
	"github.com/larntz/artui/models"
	"github.com/larntz/artui/ui/keys"
	"github.com/larntz/artui/ui/state"
	"github.com/larntz/artui/ui/templates"
)

// InitializeModel creates the initial model struct
func InitializeModel(sessionRequest session.SessionCreateRequest, apiClient apiclient.ClientOptions) state.ArTUIModel {

	//var apps v1alpha1.ApplicationList
	apps := argo.GetApplications(sessionRequest, apiClient)

	// Initialize Application List
	appList := initAppList(apps)

	// Initialize TextInput
	ti := initTextInput()

	// Initialize Templates
	tpl := initTemplates()

	return state.ArTUIModel{
		ArgoSessionRequest: sessionRequest,
		APIClient:          apiClient,
		Ready:              false,
		Activity:           state.View,
		List:               appList,
		Applications:       apps,
		Textinput:          ti,
		Templates:          tpl,
	}
}

func initTextInput() textinput.Model {
	ti := textinput.New()
	ti.SetCursorMode(textinput.CursorHide)
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
	appList := list.New(appsListModel, list.NewDefaultDelegate(), 0, 0)
	appList.Title = "App List"
	appList.KeyMap = keys.AppListKeyBinding
	appList.SetShowTitle(true)
	appList.SetShowPagination(true)
	appList.SetShowHelp(false)
	appList.SetShowFilter(true)
	appList.SetFilteringEnabled(true)

	return appList
}
