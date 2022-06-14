package models

import (
	"log"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

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
		m.Viewport.Width = int(float32(m.WindowWidth) * 0.70)
		m.Viewport.Height = contentHeight
		m.Viewport, cmd = m.Viewport.Update(message)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}
