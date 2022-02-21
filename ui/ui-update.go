package ui

import (
	"bytes"
	"log"
	"reflect"
	"text/template"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"

	"github.com/larntz/artui/ui/templates"
)

// Update the app model
func (m ArTUIModel) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	log.Printf("message type = %s, message = %s, activity = %d", reflect.TypeOf(message), message, m.Activity)
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)
	tpl, err := template.New("status").Parse(templates.AppOverviewTemplate)
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
			log.Printf("ERROR: %s", msg.Error())
			return m, nil
		}
		m.Textinput, cmd = m.Textinput.Update(message)
		return m, cmd
	} else if m.Activity == view {

		switch msg := message.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case ":":
				m.Textinput.Focus()
				m.Textinput.Prompt = ":"
				m.Activity = input

				return m, nil

			case "ctrl+c":
				return m, tea.Quit

			case "tab", "n":
				m.List.CursorDown()
				// find and update content view
				for _, v := range m.Applications.Items {
					log.Printf("Searching for %s, got %s", m.List.SelectedItem().FilterValue(), v.Name)
					if v.Name == m.List.SelectedItem().FilterValue() {
						buf := new(bytes.Buffer)
						tpl.Execute(buf, v)

						content, err := m.Glamour.Render(buf.String())
						if err != nil {
							panic(err)
						}
						m.Viewport.SetContent(content)
						m.Viewport.YOffset = 1
						break
					}
				}
				return m, nil

			case "shift+tab", "p":
				m.List.CursorUp()

				// find and update content view
				for _, v := range m.Applications.Items {
					if v.Name == m.List.SelectedItem().FilterValue() {
						buf := new(bytes.Buffer)
						tpl.Execute(buf, v)
						content, err := m.Glamour.Render(buf.String())
						if err != nil {
							panic(err)
						}
						m.Viewport.SetContent(content)
						m.Viewport.YOffset = 1
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

				log.Printf("Got WindowSizeMsg, !m.Ready")
				m.Viewport = viewport.New(msg.Width, msg.Height-verticalMarginHeight-1)
				m.Viewport.YPosition = headerHeight + 1
				m.Viewport.YOffset = 1
				m.Viewport.KeyMap.Up.SetKeys("up")
				m.Viewport.KeyMap.Down.SetKeys("down")
				m.Viewport.MouseWheelEnabled = true

				m.Glamour, err = glamour.NewTermRenderer(
					glamour.WithStandardStyle("dark"),
					glamour.WithWordWrap(m.Viewport.Width-5))
				if err != nil {
					log.Panicf("glamour problem: %s", err.Error())
				}
				log.Printf("Re-wide glamour 1: m.Viewport.Width-5=%d", m.Viewport.Width-5)

				m.List.SetHeight(msg.Height - verticalMarginHeight - 1)
				for _, v := range m.Applications.Items {
					if v.Name == m.List.SelectedItem().FilterValue() {
						buf := new(bytes.Buffer)
						tpl.Execute(buf, v)
						// content, err := m.Glamour.Render(buf.String())
						content, err := glamour.Render(buf.String(), "dark")
						if err != nil {
							panic(err)
						}

						m.Viewport.SetContent(content)
						break
					}
				}

				// This is only necessary for high performance rendering, which in
				// most cases you won't need.
				//
				// Render the viewport one line below the header.
				// m.Viewport.YPosition = headerHeight + 1
				m.Ready = true
			} else {
				log.Printf("Got WindowSizeMsg, m.Ready")
				m.Viewport.Width = msg.Width - m.List.Width()
				m.Viewport.Height = msg.Height - verticalMarginHeight - 1
				m.List.SetHeight(msg.Height - verticalMarginHeight - 1)

				m.Glamour, err = glamour.NewTermRenderer(
					glamour.WithStandardStyle("dark"),
					glamour.WithWordWrap(m.Viewport.Width-5))
				if err != nil {
					log.Panicf("glamour problem: %s", err.Error())
				}
				log.Printf("Re-wide glamour 2: m.Viewport.Width-5=%d", m.Viewport.Width-5)

				for _, v := range m.Applications.Items {
					if v.Name == m.List.SelectedItem().FilterValue() {
						buf := new(bytes.Buffer)
						tpl.Execute(buf, v)
						content, err := m.Glamour.Render(buf.String())
						if err != nil {
							panic(err)
						}
						m.Viewport.SetContent(content)
						break
					}
				}
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
	return nil, nil
}
