package main

import (
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/frozenkro/mcpsequencer/internal/db"
	"github.com/frozenkro/mcpsequencer/internal/globals"
	"github.com/frozenkro/mcpsequencer/internal/tui"
	"github.com/frozenkro/mcpsequencer/internal/utils"
)

func main() {
	env := globals.Prod
	if utils.IsDev() {
		env = globals.Dev
	}
	globals.Init(env)
	db.Init()

	m, err := tui.InitialModel()
	if err != nil {
		log.Fatalf("Error initializing state: %v", err.Error())
	}
	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
