package ui

import (
	"bytes"
	"fmt"
	"log"
	"strings"
	"text/template"

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
	Applications []models.Application
	List         list.Model
	Viewport     viewport.Model
	Textinput    textinput.Model
}

// InitialModel creates the initial model struct
func InitialModel(apps []models.Application) ArTUIModel {
	var appsListModel []list.Item
	for _, v := range apps {
		appsListModel = append(appsListModel, v)
	}

	appList := list.New(appsListModel, list.NewDefaultDelegate(), 0, 0)
	appList.Title = "App List"
	appList.SetShowTitle(true)
	appList.SetShowPagination(true)
	appList.SetShowHelp(false)
	// appList.SetShowStatusBar(true)
	// appList.NewStatusMessage("got some apps")
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
	return tea.EnterAltScreen
}

// Update the app model
func (m ArTUIModel) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	log.Printf("message.(type)=%t, message=%+v", message, message)
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)
	tpl, err := template.New("status").Parse(statusTemplate)
	if err != nil {
		panic("can't create template")
	}

	if m.Activity == input {
		switch msg := message.(type) {
		case tea.KeyMsg:
			switch msg.Type {
			case tea.KeyEsc:
				m.Textinput.Prompt = " "
				m.Textinput.Reset()
				m.Activity = view
				return m, nil

			case tea.KeyEnter:
				switch m.Textinput.Value() {
				case "q", "quit":
					return m, tea.Quit
				}

				return m, nil
			}

		// We handle errors just like any other message
		case errMsg:
			log.Printf("Error %s", msg.Error())
			return m, nil
		}
		m.Textinput, cmd = m.Textinput.Update(message)
		return m, cmd
	}

	switch msg := message.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case ":":
			log.Printf("got key '%s' : Textinput", msg.String())
			m.Textinput.Focus()
			m.Textinput.Prompt = ":"
			//m.Textinput.Update(m)
			m.Activity = input

			return m, nil

		case "ctrl+c":
			return m, tea.Quit

		case "tab", "n":
			log.Printf("got key '%s' : m.List.CursorDown()", msg.String())
			m.List.CursorDown()
			// find and update content view
			for _, v := range m.Applications {
				if v.Name == m.List.SelectedItem().FilterValue() {
					buf := new(bytes.Buffer)
					tpl.Execute(buf, v)
					content, err := glamour.Render(buf.String(), "dark")
					if err != nil {
						panic(err)
					}
					m.Viewport.SetContent(content)
					m.Viewport.YOffset = 0
					break
				}
			}
			return m, nil

		case "shift+tab", "p":
			log.Printf("got key '%s' : m.List.CursorUp()", msg.String())
			m.List.CursorUp()

			// find and update content view
			for _, v := range m.Applications {
				if v.Name == m.List.SelectedItem().FilterValue() {
					buf := new(bytes.Buffer)
					tpl.Execute(buf, v)
					content, err := glamour.Render(buf.String(), "dark")
					if err != nil {
						panic(err)
					}
					m.Viewport.SetContent(content)
					m.Viewport.YOffset = 0
					break
				}
			}
			return m, nil

		}

	case tea.WindowSizeMsg:
		headerHeight := lipgloss.Height(m.headerView())
		footerHeight := lipgloss.Height(m.footerView())
		verticalMarginHeight := headerHeight + footerHeight

		if !m.Ready {
			// Since this program is using the full size of the viewport we
			// need to wait until we've received the window dimensions before
			// we can initialize the viewport. The initial dimensions come in
			// quickly, though asynchronously, which is why we wait for them
			// here.

			m.Viewport = viewport.New(msg.Width, msg.Height-verticalMarginHeight-1)
			m.Viewport.YPosition = headerHeight + 1
			m.Viewport.KeyMap.Up.SetKeys("up")
			m.Viewport.KeyMap.Down.SetKeys("down")
			m.Viewport.MouseWheelEnabled = true

			m.List.SetHeight(msg.Height - verticalMarginHeight - 1)
			for _, v := range m.Applications {
				if v.Name == m.List.SelectedItem().FilterValue() {
					buf := new(bytes.Buffer)
					tpl.Execute(buf, v)
					content, err := glamour.Render(buf.String(), "dark")
					if err != nil {
						panic(err)
					}

					m.Viewport.SetContent(content)
					break
				}
			}
			log.Printf("m.Ready, msg.Width %d, viewport.Width %d, appList.Width %d", msg.Width, m.Viewport.Width, m.List.Width())

			// This is only necessary for high performance rendering, which in
			// most cases you won't need.
			//
			// Render the viewport one line below the header.
			m.Viewport.YPosition = headerHeight + 1
			m.Ready = true
		} else {
			m.Viewport.Width = msg.Width - m.List.Width()
			m.Viewport.Height = msg.Height - verticalMarginHeight - 1
			m.List.SetHeight(msg.Height - verticalMarginHeight - 1)
			for _, v := range m.Applications {
				if v.Name == m.List.SelectedItem().FilterValue() {
					content := fmt.Sprintf("# %s\n\n```yaml\n%s\n```\n\n", v.Name, v.LongStatus)
					buf := new(bytes.Buffer)
					tpl.Execute(buf, v)
					content, err := glamour.Render(buf.String(), "dark")
					if err != nil {
						panic(err)
					}
					m.Viewport.SetContent(content)
					// m.Viewport.SetContent(
					// 	wrap.String(content, m.Viewport.Width-10))
					break
				}
			}
			log.Printf("m.Ready, msg.Width %d, viewport.Width %d, appList.Width %d", msg.Width, m.Viewport.Width, m.List.Width())
		}
	}

	m.List, cmd = m.List.Update(message)
	cmds = append(cmds, cmd)

	m.Viewport, cmd = m.Viewport.Update(message)
	cmds = append(cmds, cmd)

	m.Textinput, cmd = m.Textinput.Update(message)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

// View the model?
func (m ArTUIModel) View() string {
	// return style.Render(fmt.Sprintf("ArgoCD Apps"))
	// return style.Render(m.viewport.View())
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
		m.Viewport.Width-lipgloss.Width(footerStyle.Render(message))-lipgloss.Width(footerStyle.Render(textInput)))) // lipgloss.Width(footerStyle.Render(textInput))))

	return footerStyle.Render(lipgloss.JoinHorizontal(lipgloss.Center, textInput, line, message))
}
