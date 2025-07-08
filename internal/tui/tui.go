package tui

import (
	// "fmt"
	// "os"

	"context"

	tea "github.com/charmbracelet/bubbletea"
	// "github.com/frozenkro/mcpsequencer/internal/db"
	// "github.com/frozenkro/mcpsequencer/internal/globals"
	"github.com/frozenkro/mcpsequencer/internal/projectsdb"
	"github.com/frozenkro/mcpsequencer/internal/services"
	// "github.com/frozenkro/mcpsequencer/internal/utils"
)

var svc services.Services
var ctx context.Context

type model struct {
	projects      []projectsdb.Project
	activeProject *projectsdb.Project
	tasks         []projectsdb.Task
	cursor        int
	activeWindow  ActiveWindow
}

type ActiveWindow int

const (
	ProjectWindow ActiveWindow = iota
	TasksWindow
)

func InitialModel() (model, error) {
	ctx = context.Background()
	svc = services.Services{}
	projects, err := svc.GetProjects(ctx)
	if err != nil {
		return model{}, err
	}

	return model{
		projects: projects,
	}, nil
}

func (m model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// return nil, nil
	switch msg := msg.(type) {

	// Is it a key press?
	case tea.KeyMsg:

		// Cool, what was the actual key pressed?
		switch msg.String() {

		// These keys should exit the program.
		case "ctrl+c", "q":
			return m, tea.Quit

		// The "up" and "k" keys move the cursor up
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		// The "down" and "j" keys move the cursor down
		case "down", "j":
			if m.cursor < m.activeListLen()-1 {
				m.cursor++
			}

		// The "enter" key and the spacebar (a literal space) toggle
		// the selected state for the item that the cursor is pointing at.
		case "enter", " ":
			m.handleSelect()
		}
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}

func (m model) View() string {
	return ""
	// The header
	// s := "What should we buy at the market?\n\n"

	// // Iterate over our choices
	// for i, choice := range m.choices {

	// 	// Is the cursor pointing at this choice?
	// 	cursor := " " // no cursor
	// 	if m.cursor == i {
	// 		cursor = ">" // cursor!
	// 	}

	// 	// Is this choice selected?
	// 	checked := " " // not selected
	// 	if _, ok := m.selected[i]; ok {
	// 		checked = "x" // selected!
	// 	}

	// 	// Render the row
	// 	s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
	// }

	// // The footer
	// s += "\nPress q to quit.\n"

	// // Send the UI for rendering
	// return s
}

func (m model) activeListLen() int {
	if m.activeWindow == ProjectWindow {
		return len(m.projects)
	}
	if m.activeWindow == TasksWindow {
		return len(m.tasks)
	}
	return 0
}

func (m model) handleSelect() {
	if m.activeWindow == ProjectWindow {
		m.activeProject = &m.projects[m.cursor]
	}
}
