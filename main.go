package main

import (
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/larntz/artui/argo"
	"github.com/larntz/artui/ui"
)

func main() {
	// setup logging
	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		fmt.Println("fatal:", err)
		os.Exit(1)
	}
	defer f.Close()
	log.Println("Application Start")
	apps := argo.GetApplications()

	log.Println("Got Applications")

	// start application
	log.Println("UI Start")
	p := tea.NewProgram(ui.InitialModel(apps), tea.WithAltScreen(), tea.WithMouseCellMotion())
	if err := p.Start(); err != nil {
		panic(err)
	}
	log.Println("Application Exit")
}
