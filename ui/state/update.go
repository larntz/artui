package state

import (
	"log"
	"reflect"

	"github.com/argoproj/argo-cd/v2/pkg/apis/application/v1alpha1"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

// Find correct application
func getApplication(m ArTUIModel) v1alpha1.Application {
	for _, v := range m.Applications.Items {
		if v.Name == m.List.SelectedItem().FilterValue() {
			return v
		}
	}
	log.Printf("failed to find application")
	return v1alpha1.Application{}
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
				return m, tea.Quit
			}
			return m, nil
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

// Handle WindowSizeMsg
func handleWindowSizeMsg(m ArTUIModel, message tea.WindowSizeMsg) (tea.Model, tea.Cmd) {
	headerHeight := lipgloss.Height(m.headerView())
	footerHeight := lipgloss.Height(m.footerView())
	verticalMarginHeight := headerHeight + footerHeight
	var cmd tea.Cmd
	var cmds []tea.Cmd

	if !m.Ready {
		// Since this program is using the full size of the viewport we
		// need to wait until we've received the window dimensions before
		// we can initialize the viewport. The initial dimensions come in
		// quickly, though asynchronously, which is why we wait for them
		// here.

		log.Printf("Got WindowSizeMsg, !m.Ready")
		m.Viewport = viewport.New(message.Width, message.Height-verticalMarginHeight-1)
		m.Viewport.YPosition = headerHeight + 1
		m.Viewport.YOffset = 1
		m.Viewport.KeyMap.Up.SetKeys("up")
		m.Viewport.KeyMap.Down.SetKeys("down")
		m.Viewport.MouseWheelEnabled = true

		m.Viewport, cmd = m.Viewport.Update(message)
		cmds = append(cmds, cmd)

		var err error
		m.Glamour, err = glamour.NewTermRenderer(
			glamour.WithStandardStyle("dark"),
			glamour.WithWordWrap(m.Viewport.Width-5))
		if err != nil {
			log.Panicf("glamour problem: %s", err.Error())
		}
		log.Printf("Re-wide glamour 1: m.Viewport.Width-5=%d", m.Viewport.Width-5)

		m.List.SetHeight(message.Height - verticalMarginHeight - 1)
		m.List, cmd = m.List.Update(message)
		cmds = append(cmds, cmd)

		markdown, err := m.renderTemplate("AppOverviewTemplate")
		if err != nil {
			log.Panicf("86: %s", err.Error())
		}
		m.Viewport.SetContent(markdown)
		m.Ready = true

	} else {
		log.Printf("Got WindowSizeMsg, m.Ready")
		m.Viewport.Width = message.Width - m.List.Width()
		m.Viewport.Height = message.Height - verticalMarginHeight - 1
		m.Viewport, cmd = m.Viewport.Update(message)
		cmds = append(cmds, cmd)

		m.List.SetHeight(message.Height - verticalMarginHeight - 1)
		m.List, cmd = m.List.Update(message)
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
				m.Viewport.YOffset = 1
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
				m.Viewport.YOffset = 1
				m.Viewport, cmd = m.Viewport.Update(message)
				cmds = append(cmds, cmd)
				return m, tea.Batch(cmds...)

			} // end inner switch

		case tea.WindowSizeMsg:
			return handleWindowSizeMsg(m, msg)
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
