package models

import (
	"bytes"
	"fmt"
	"log"
	"reflect"
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

// Update the app model
func (m ArTUIModel) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	log.Printf("message type = %s, message = %s, activity = %d", reflect.TypeOf(message), message, m.Activity)
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	if !m.Ready {
		// Don't do anything until we have our first WindowSizeMsg.
		switch msg := message.(type) {
		case tea.WindowSizeMsg:
			return handleWindowSizeMsg(m, msg)
		default:
			return m, nil
		}
	} else if m.Activity == Input {
		return inputUpdate(m, message)
	} else if m.Activity == View {
		switch msg := message.(type) {

		case GetApplicationMsg:
			log.Printf("GetApplicationMsg recieved. len(msg.applications) = %d", len(msg.applications.Items))
			selected := m.List.Index()
			m.Applications = msg.applications
			m.List = m.updateAppList()
			m.List.Select(selected)
			m.List.Update(msg)
			cmds = append(cmds, cmd)
			markdown, err := m.renderTemplate("AppOverviewTemplate")
			if err != nil {
				log.Panicf("144: %s", err.Error())
			}
			m.Viewport.SetContent(markdown)
			m.Viewport, cmd = m.Viewport.Update(message)
			cmds = append(cmds, cmd)

		case tea.KeyMsg:
			switch msg.String() {
			case ":":
				m.Textinput.Focus()
				m.Textinput.Prompt = ":"
				m.Activity = Input

				return m, nil

			case "ctrl+c":
				return m, tea.Quit

			case "j":
				m.List, cmd = m.List.Update(message)
				cmds = append(cmds, cmd)
				markdown, err := m.renderTemplate("AppOverviewTemplate")
				if err != nil {
					log.Panicf("144: %s", err.Error())
				}
				m.Viewport.SetContent(markdown)
				m.Viewport.SetYOffset(0)
				m.Viewport, cmd = m.Viewport.Update(message)
				cmds = append(cmds, cmd)
				return m, tea.Batch(cmds...)
			case "k":
				m.List, cmd = m.List.Update(message)
				cmds = append(cmds, cmd)
				markdown, err := m.renderTemplate("AppOverviewTemplate")
				if err != nil {
					log.Panicf("144: %s", err.Error())
				}
				m.Viewport.SetContent(markdown)
				m.Viewport.SetYOffset(0)
				m.Viewport, cmd = m.Viewport.Update(message)
				cmds = append(cmds, cmd)
				return m, tea.Batch(cmds...)

			} // end inner switch

		case tea.WindowSizeMsg:
			return handleWindowSizeMsg(m, msg)
		}

		if m.List.ShowFilter() && m.Ready {
			m.List.SetWidth(int(float32(m.WindowWidth) * 0.50))
			m.Viewport.Width = int(float32(m.WindowWidth) * 0.50)
		} else {
			m.List.SetWidth(int(float32(m.WindowWidth) * 0.25))
			m.Viewport.Width = int(float32(m.WindowWidth) * 0.70)
		}

		m.List, cmd = m.List.Update(message)
		cmds = append(cmds, cmd)

		m.Viewport, cmd = m.Viewport.Update(message)
		cmds = append(cmds, cmd)

		m.Textinput, cmd = m.Textinput.Update(message)
		cmds = append(cmds, cmd)
	}
	return m, tea.Batch(cmds...)
}

func (m ArTUIModel) headerView() string {
	title := fmt.Sprintf(" ArTUI: Managing ArgoCD Apps")
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

func (m ArTUIModel) asdf() {

}
