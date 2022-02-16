package ui

import (
	"bytes"
	"fmt"
	"log"
	"strings"
	"text/template"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/wrap"

	"github.com/larntz/artui/models"
	"github.com/larntz/artui/utils"
)

// ArTUIModel is the bubbletea app model
type ArTUIModel struct {
	Ready        bool
	Applications []models.Application
	List         list.Model
	Viewport     viewport.Model
}

// InitialModel creates the initial model struct
func InitialModel(apps []models.Application) ArTUIModel {
	var appsListModel []list.Item
	for _, v := range apps {
		appsListModel = append(appsListModel, v)
	}

	appList := list.New(appsListModel, list.NewDefaultDelegate(), 0, 25)
	appList.Title = "App List"
	appList.SetShowTitle(true)
	appList.SetShowPagination(true)
	//appList.SetShowHelp(false)
	appList.SetShowStatusBar(true)
	appList.SetShowFilter(true)
	appList.SetFilteringEnabled(true)

	return ArTUIModel{
		Ready:        false,
		List:         appList,
		Applications: apps,
	}
}

// Init the app model
func (m ArTUIModel) Init() tea.Cmd {
	return tea.EnterAltScreen
}

// Update the app model
func (m ArTUIModel) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)
	tpl, err := template.New("status").Parse(statusTemplate)
	if err != nil {
		panic("can't create template")
	}

	switch msg := message.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case "f":
			log.Printf("got key '%s' : show filter?", msg.String())
			m.List.SetShowFilter(true)
			m.List.FilterInput.Focus()
			return m, nil

		case "q", "esc", "ctrl+c":
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
					m.Viewport.SetContent(
						content,
					)
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
					m.Viewport.SetContent(
						content,
					)
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
			m.Viewport.HighPerformanceRendering = true
			m.Viewport = viewport.New(msg.Width, msg.Height-verticalMarginHeight-1)
			m.Viewport.YPosition = headerHeight
			m.List.SetHeight(msg.Height - verticalMarginHeight - 1)
			// m.Viewport.Style.Border(lipgloss.ThickBorder())
			// m.Viewport.Style.BorderForeground(lipgloss.Color("198"))
			for _, v := range m.Applications {
				if v.Name == m.List.SelectedItem().FilterValue() {
					buf := new(bytes.Buffer)
					tpl.Execute(buf, v)
					content, err := glamour.Render(buf.String(), "dark")
					if err != nil {
						panic(err)
					}

					m.Viewport.SetContent(
						wrap.String(content, msg.Width-25),
					)
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
					m.Viewport.SetContent(
						wrap.String(content, m.Viewport.Width-25),
					)
					break
				}
			}
			log.Printf("m.Ready, msg.Width %d, viewport.Width %d, appList.Width %d", msg.Width, m.Viewport.Width, m.List.Width())
		}

	}
	m.Viewport, cmd = m.Viewport.Update(message)
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
	message := fmt.Sprintf("https://github.com/larntz")
	line := strings.Repeat(" ", utils.Max(0, m.Viewport.Width-lipgloss.Width(footerStyle.Render(message))))
	return footerStyle.Render(lipgloss.JoinHorizontal(lipgloss.Center, line, message))
}
