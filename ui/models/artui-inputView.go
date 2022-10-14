package models

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
)

// Handle updates when m.Activity = Input
func inputUpdate(m ArTUIModel, message tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := message.(type) {

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc:
			m.Textinput.Prompt = " "
			m.Textinput.SetValue("")
			m.Textinput.Blur()
			m.Activity = View
			return m, nil
		case tea.KeyEnter:

			switch m.Textinput.Value() {
			case "q", "quit":
				m.Textinput.Prompt = " "
				m.Textinput.SetValue("")
				m.Textinput.Blur()
				m.Activity = View
				return m, tea.Quit

			case "r", "refresh":
				m.Textinput.Prompt = " "
				m.Textinput.SetValue("")
				m.Textinput.Blur()
				m.Activity = View
				app, err := getApplication(m)
				if err != nil {
					log.Printf("Cannot find application to refresh...")
				}
				m.AppWorkerChan <- WorkerCmd{
					Cmd: Refresh,
					App: app,
				}
				return m, nil

			case "hr", "hard-refresh":
				m.Textinput.Prompt = " "
				m.Textinput.SetValue("")
				m.Textinput.Blur()
				m.Activity = View
				app, err := getApplication(m)
				if err != nil {
					log.Printf("Cannot find application to refresh...")
				}
				m.AppWorkerChan <- WorkerCmd{
					Cmd: HardRefresh,
					App: app,
				}
				return m, nil

			case "s", "sync":
				m.Textinput.Prompt = " "
				m.Textinput.SetValue("")
				m.Textinput.Blur()
				m.Activity = View
				app, err := getApplication(m)
				if err != nil {
					log.Printf("Cannot find application to refresh...")
				}
				m.AppWorkerChan <- WorkerCmd{
					Cmd: Sync,
					App: app,
				}
				return m, nil

			default:
				m.Textinput.Prompt = " "
				m.Textinput.SetValue("")
				m.Textinput.Blur()
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
