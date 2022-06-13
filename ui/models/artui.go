package models

import (
	"bytes"
	"fmt"
	"log"
	"strings"
	"text/template"

	"github.com/argoproj/argo-cd/v2/pkg/apiclient"
	"github.com/argoproj/argo-cd/v2/pkg/apiclient/session"
	"github.com/argoproj/argo-cd/v2/pkg/apis/application/v1alpha1"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"

	styles "github.com/larntz/artui/ui/styles"
	"github.com/larntz/artui/utils"
)

type errMsg error
type mode int

// Application Modes
const (
	View mode = iota
	Filter
	Input
)

// ArTUIModel is the bubbletea app model
type ArTUIModel struct {
	APIClient          apiclient.ClientOptions
	ArgoSessionRequest session.SessionCreateRequest
	Ready              bool
	Activity           mode
	Applications       v1alpha1.ApplicationList
	List               list.Model
	Viewport           viewport.Model
	Textinput          textinput.Model
	Glamour            *glamour.TermRenderer
	Templates          *template.Template
	WindowHeight       int
	WindowWidth        int
}

// Init the app model
func (m ArTUIModel) Init() tea.Cmd {
	return tea.Batch(GetApplications(m))
}

// View the model
func (m ArTUIModel) View() string {
	return lipgloss.JoinVertical(lipgloss.Left,
		m.headerView(),
		lipgloss.JoinHorizontal(lipgloss.Top, styles.AppListStyle.Render(m.List.View()), styles.ViewportStyle.Render(m.Viewport.View())),
		m.footerView())
}

func (m ArTUIModel) headerView() string {
	title := fmt.Sprintf("ArTUI: Managing ArgoCD Apps")
	line := strings.Repeat(" ", utils.Max(0, m.WindowWidth-lipgloss.Width(styles.HeaderStyle.Render(title))))

	return styles.HeaderStyle.Render(lipgloss.JoinHorizontal(lipgloss.Center, title, line))
}

func (m ArTUIModel) footerView() string {
	textInput := m.Textinput.View()
	message := fmt.Sprintf("https://github.com/larntz/artui")
	line := strings.Repeat(" ", utils.Max(0,
		m.WindowWidth-lipgloss.Width(styles.FooterStyle.Render(message))-lipgloss.Width(styles.FooterStyle.Render(textInput))))

	return styles.FooterStyle.Render(lipgloss.JoinHorizontal(lipgloss.Center, textInput, line, message))
}

// Update Viewport Content
func (m ArTUIModel) renderTemplate(templateName string) (string, error) {
	log.Printf("renderTemplate: %s", templateName)
	app, err := getApplication(m)
	if err != nil {
		return m.Glamour.Render("")
	}

	buf := new(bytes.Buffer)
	if err := m.Templates.ExecuteTemplate(buf, templateName, app); err != nil {
		log.Panicf("templateRender failed\n:\t%s", err.Error())
	}
	return m.Glamour.Render(buf.String())
}
