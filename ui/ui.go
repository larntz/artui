package ui

import (
	"fmt"
	"log"
	"strings"

	"gopkg.in/yaml.v2"

	"github.com/argoproj/argo-cd/v2/pkg/apis/application/v1alpha1"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"

	"github.com/larntz/artui/models"
	"github.com/larntz/artui/utils"
)

type errMsg error
type mode int

const (
	view mode = iota
	filter
	input
)

// ArTUIModel is the bubbletea app model
type ArTUIModel struct {
	Ready        bool
	Activity     mode
	Applications v1alpha1.ApplicationList
	List         list.Model
	Viewport     viewport.Model
	Textinput    textinput.Model
	Glamour      *glamour.TermRenderer
}

// InitialModel creates the initial model struct
func InitialModel(apps v1alpha1.ApplicationList) ArTUIModel {
	var appsListModel []list.Item
	for _, item := range apps.Items {
		longStatus, err := yaml.Marshal(item.Status)
		if err != nil {
			log.Fatal(err)
		}

		appsListModel = append(appsListModel, models.Application{
			Name:       item.Name,
			Status:     string(item.Status.Health.Status) + "/" + string(item.Status.Sync.Status),
			LongStatus: string(longStatus),
		})
	}

	appList := list.New(appsListModel, list.NewDefaultDelegate(), 0, 0)
	appList.Title = "App List"
	appList.SetShowTitle(true)
	appList.SetShowPagination(true)
	appList.SetShowHelp(false)
	appList.SetShowFilter(true)
	appList.SetFilteringEnabled(true)

	ti := textinput.New()
	ti.SetCursorMode(textinput.CursorHide)
	ti.Prompt = " "
	ti.PromptStyle.PaddingLeft(0)
	ti.CharLimit = 20
	ti.Width = 20

	return ArTUIModel{
		Ready:        false,
		Activity:     view,
		List:         appList,
		Applications: apps,
		Textinput:    ti,
	}
}

// Init the app model
func (m ArTUIModel) Init() tea.Cmd {
	return nil
}

// View the model?
func (m ArTUIModel) View() string {
	return fmt.Sprintf("%s\n%s\n%s\n",
		m.headerView(),
		lipgloss.JoinHorizontal(lipgloss.Top, appListStyle.Render(m.List.View()), viewportStyle.Render(m.Viewport.View())),
		m.footerView())
}

func (m ArTUIModel) headerView() string {
	title := fmt.Sprintf("ArTUI: Managing ArgoCD Apps")
	line := strings.Repeat(" ", utils.Max(0, m.Viewport.Width-lipgloss.Width(headerStyle.Render(title))))
	return headerStyle.Render(lipgloss.JoinHorizontal(lipgloss.Center, title, line))
}

func (m ArTUIModel) footerView() string {
	message := fmt.Sprintf("https://github.com/larntz/artui")
	textInput := m.Textinput.View()
	line := strings.Repeat(" ", utils.Max(0,
		m.Viewport.Width-lipgloss.Width(footerStyle.Render(message))-lipgloss.Width(footerStyle.Render(textInput))))
	return footerStyle.Render(lipgloss.JoinHorizontal(lipgloss.Center, textInput, line, message))
}
