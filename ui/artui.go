package ui

import (
	"fmt"
	"log"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/wrap"

	"github.com/larntz/artui/models"
	"github.com/larntz/artui/utils"
)

// ArTUIModel is the bubbletea app model
type ArTUIModel struct {
	Ready    bool
	Apps     list.Model
	Viewport viewport.Model
}

// InitialModel creates the initial model struct
func InitialModel(apps []models.Application) ArTUIModel {
	// apps := make([]list.Item, 0, 500)
	// apps = append(apps, models.Application{Name: "apps", Status: "Synced / Healthy"})
	// apps = append(apps, models.Application{Name: "argocd", Status: "Synced / Progressing"})
	// apps = append(apps, models.Application{Name: "prometheus", Status: "Synced / Progressing"})
	// apps = append(apps, models.Application{Name: "traefik", Status: "Synced / Healthy"})
	// apps = append(apps, models.Application{Name: "webapp", Status: "OutOfSync / Missing"})

	var appsListModel []list.Item
	for _, v := range apps {
		appsListModel = append(appsListModel, v)
	}

	appList := list.New(appsListModel, list.NewDefaultDelegate(), 0, 25)
	//appList.Title = "ArgoCD Applications on <cluster>"
	appList.SetShowTitle(false)
	appList.SetShowPagination(true)
	appList.SetShowHelp(false)
	appList.SetShowStatusBar(false)

	return ArTUIModel{
		Ready: false,
		Apps:  appList,
	}
}

// Init the app model
func (m ArTUIModel) Init() tea.Cmd {
	// return tea.Batch(tick(), tea.EnterAltScreen)
	return tea.EnterAltScreen
}

// Update the app model
func (m ArTUIModel) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := message.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			return m, tea.Quit
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
			m.Viewport.Style.Border(lipgloss.ThickBorder())
			m.Viewport.Style.BorderForeground(lipgloss.Color("198"))
			m.Viewport.SetContent(
				wrap.String(AppsJSON, msg.Width-25),
			)
			log.Printf("m.Ready, msg.Width %d, viewport.Width %d, appList.Width %d", msg.Width, m.Viewport.Width, m.Apps.Width())

			// This is only necessary for high performance rendering, which in
			// most cases you won't need.
			//
			// Render the viewport one line below the header.
			m.Viewport.YPosition = headerHeight + 1
			m.Ready = true
		} else {
			m.Viewport.Width = msg.Width - m.Apps.Width()
			m.Viewport.Height = msg.Height - verticalMarginHeight
			m.Viewport.SetContent(
				wrap.String(AppsJSON, m.Viewport.Width-25),
			)
			log.Printf("m.Ready, msg.Width %d, viewport.Width %d, appList.Width %d", msg.Width, m.Viewport.Width, m.Apps.Width())
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
		lipgloss.JoinHorizontal(lipgloss.Top, appListStyle.Render(m.Apps.View()), viewportStyle.Render(m.Viewport.View())),
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
