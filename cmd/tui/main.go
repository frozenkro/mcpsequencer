package main

import (
	"fmt"
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

	p := tea.NewProgram(tui.InitialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
