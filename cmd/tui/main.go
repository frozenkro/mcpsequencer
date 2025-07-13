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

	logger, closeLogFile, err := setupLogger()
	if err != nil {
		fmt.Printf("Error creating logger: %v\n", err.Error())
	}
	defer closeLogFile()

	db.Init()

	m, err := tui.InitialModel(logger)
	if err != nil {
		log.Fatalf("Error initializing state: %v", err.Error())
	}
	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}

func setupLogger() (*log.Logger, func() error, error) {
	file, err := os.OpenFile("mcpsequencertui.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, nil, err
	}

	logger := log.New(file, "", log.LstdFlags)

	return logger, file.Close, nil
}
