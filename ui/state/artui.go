package state

import (
	"bytes"
	"fmt"
	"log"
	"strings"
	"text/template"

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
	Ready        bool
	Activity     mode
	Applications v1alpha1.ApplicationList
	List         list.Model
	Viewport     viewport.Model
	Textinput    textinput.Model
	Glamour      *glamour.TermRenderer
	Templates    *template.Template
}

// Init the app model
func (m ArTUIModel) Init() tea.Cmd {
	return nil
}

// View the model?
func (m ArTUIModel) View() string {
	return fmt.Sprintf("%s\n%s\n%s\n",
		m.headerView(),
		lipgloss.JoinHorizontal(lipgloss.Top, styles.AppListStyle.Render(m.List.View()), styles.ViewportStyle.Render(m.Viewport.View())),
		m.footerView())
}

func (m ArTUIModel) headerView() string {
	title := fmt.Sprintf("ArTUI: Managing ArgoCD Apps")
	line := strings.Repeat(" ", utils.Max(0, m.Viewport.Width-lipgloss.Width(styles.HeaderStyle.Render(title))))
	return styles.HeaderStyle.Render(lipgloss.JoinHorizontal(lipgloss.Center, title, line))
}

func (m ArTUIModel) footerView() string {
	message := fmt.Sprintf("https://github.com/larntz/artui")
	textInput := m.Textinput.View()
	line := strings.Repeat(" ", utils.Max(0,
		m.Viewport.Width-lipgloss.Width(styles.FooterStyle.Render(message))-lipgloss.Width(styles.FooterStyle.Render(textInput))))
	return styles.FooterStyle.Render(lipgloss.JoinHorizontal(lipgloss.Center, textInput, line, message))
}

// Update Viewport Content
func (m ArTUIModel) renderTemplate(templateName string) (string, error) {
	app := getApplication(m)

	buf := new(bytes.Buffer)
	log.Printf("len(Templates()) == %d", len(m.Templates.Templates()))
	err := m.Templates.ExecuteTemplate(buf, templateName, app)
	if err != nil {
		log.Panicf("templateRender failed\n:%s", err.Error())
	}

	return m.Glamour.Render(buf.String())
}
