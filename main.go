package main

import (
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/larntz/artui/ui"
)

func main() {
	// setup loggin
	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		fmt.Println("fatal:", err)
		os.Exit(1)
	}
	defer f.Close()
	log.Println("Application Start")

	// start application
	p := tea.NewProgram(ui.InitialModel(), tea.WithAltScreen())
	if err := p.Start(); err != nil {
		log.Fatal(err)
	}
	log.Println("Application Exit")
}
