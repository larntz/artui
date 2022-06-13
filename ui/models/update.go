package models

import (
	"errors"
	"log"
	"reflect"

	"github.com/argoproj/argo-cd/v2/pkg/apis/application/v1alpha1"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"

	"github.com/larntz/artui/ui/keys"
)

// Find correct application
func getApplication(m ArTUIModel) (v1alpha1.Application, error) {
	log.Printf("getApplications: len(m.Applications.Items) = %d", len(m.Applications.Items))
	if m.List.SelectedItem() != nil {
		log.Printf("getApplication: SelectedItem = %v", m.List.SelectedItem().FilterValue())
	}
	for _, v := range m.Applications.Items {
		log.Printf("getApplication: want=%s, got=%s", v.Name, m.List.SelectedItem().FilterValue())
		if v.Name == m.List.SelectedItem().FilterValue() {
			return v, nil
		}
	}
	log.Printf("getApplication: failed to find application")
	return v1alpha1.Application{}, errors.New("failed-to-find-app")
}

// Handle updates during input
func inputUpdate(m ArTUIModel, message tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := message.(type) {

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc:
			m.Textinput.Prompt = " "
			m.Textinput.Reset()
			m.Activity = View
			return m, nil
		case tea.KeyEnter:

			switch m.Textinput.Value() {
			case "q", "quit":
				m.Textinput.Prompt = " "
				m.Textinput.Reset()
				return m, tea.Quit

			case "r", "refresh-applications":
				log.Printf("User requested application refresh")
				m.Textinput.Prompt = " "
				m.Textinput.Reset()
				m.Activity = View
				return m, GetApplications(m)

			default:
				m.Textinput.Prompt = " "
				m.Textinput.Reset()
				m.Activity = View
				return m, nil
			}
		}

	// We handle errors just like any other message
	case errMsg:
		log.Printf("ERROR: %s", msg.Error())
		return m, nil
	}

	newModel, cmd := m.Textinput.Update(message)
	m.Textinput = newModel
	return m, cmd
}

// Return windowHeight - (header + footer)
func getContentHeight(m ArTUIModel) int {
	headerHeight := lipgloss.Height(m.headerView())
	footerHeight := lipgloss.Height(m.footerView())

	return m.WindowHeight - (headerHeight + footerHeight + 1)
}

// Handle WindowSizeMsg
func handleWindowSizeMsg(m ArTUIModel, message tea.WindowSizeMsg) (tea.Model, tea.Cmd) {
	m.WindowHeight = message.Height
	m.WindowWidth = message.Width
	contentHeight := getContentHeight(m)

	var cmd tea.Cmd
	var cmds []tea.Cmd

	if !m.Ready {
		// Since this program is using the full size of the viewport we
		// need to wait until we've received the window dimensions before
		// we can initialize the viewport. The initial dimensions come in
		// quickly, though asynchronously, which is why we wait for them
		// here.

		log.Printf("!m.Ready")
		m.Viewport = viewport.New(int(float32(m.WindowWidth)*0.70), contentHeight)
		m.Viewport.YPosition = lipgloss.Height(m.headerView())
		m.Viewport.KeyMap.Up.SetKeys("up")
		m.Viewport.KeyMap.Down.SetKeys("down")
		m.Viewport.MouseWheelEnabled = true

		var err error
		m.Glamour, err = glamour.NewTermRenderer(
			glamour.WithStandardStyle("dark"),
			glamour.WithWordWrap(m.Viewport.Width-5))
		if err != nil {
			log.Panicf("glamour problem: %s", err.Error())
		}

		m.List, cmd = m.List.Update(message)
		m.List.SetWidth(int(float32(m.WindowWidth) * 0.25))
		m.List.SetHeight(getContentHeight(m))
		cmds = append(cmds, cmd)

		markdown, err := m.renderTemplate("AppOverviewTemplate")
		if err != nil {
			log.Panicf("86: %s", err.Error())
		}
		m.Viewport.SetContent(markdown)
		m.Viewport, cmd = m.Viewport.Update(message)
		cmds = append(cmds, cmd)
		m.Ready = true

	} else {
		log.Printf("m.Ready")
		m.Viewport.Width = int(float32(m.WindowWidth) * 0.70)
		m.Viewport.Height = contentHeight
		m.Viewport, cmd = m.Viewport.Update(message)
		cmds = append(cmds, cmd)

		m.List, cmd = m.List.Update(message)
		m.List.SetWidth(int(float32(m.WindowWidth) * 0.25))
		m.List.SetHeight(getContentHeight(m))
		cmds = append(cmds, cmd)

		var err error
		m.Glamour, err = glamour.NewTermRenderer(
			glamour.WithStandardStyle("dark"),
			glamour.WithWordWrap(m.Viewport.Width-5))
		if err != nil {
			log.Panicf("glamour problem: %s", err.Error())
		}
		log.Printf("Re-wide glamour 2: m.Viewport.Width-5=%d", m.Viewport.Width-5)

		markdown, err := m.renderTemplate("AppOverviewTemplate")
		if err != nil {
			log.Panicf("108: %s", err.Error())
		}
		m.Viewport.SetContent(markdown)
	}

	return m, tea.Batch(cmds...)
}

// Update the app model
func (m ArTUIModel) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	log.Printf("message type = %s, message = %s, activity = %d", reflect.TypeOf(message), message, m.Activity)
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	if m.Activity == Input {
		return inputUpdate(m, message)
	} else if m.Activity == View {
		switch msg := message.(type) { // inner switch

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

	appList := list.New(appListItems, list.NewDefaultDelegate(), int(float32(m.WindowWidth)*0.25), getContentHeight(m))
	appList.Title = "App List"
	appList.KeyMap = keys.AppListKeyBinding
	appList.SetShowTitle(true)
	appList.SetShowPagination(true)
	appList.SetShowHelp(false)
	appList.SetShowFilter(true)
	appList.SetFilteringEnabled(true)

	return appList
}
