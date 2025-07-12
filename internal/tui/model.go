package tui

import (
	"context"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/frozenkro/mcpsequencer/internal/services"
	"github.com/frozenkro/mcpsequencer/internal/tui/components/projects"
	"github.com/frozenkro/mcpsequencer/internal/tui/components/tasks"
	"github.com/frozenkro/mcpsequencer/internal/tui/constants"
)

type ActivePane int

const (
	ProjectPane ActivePane = iota
	TasksPane
)

type Model struct {
	Projects   projects.Model
	Tasks      tasks.Model
	ActivePane ActivePane
	Width      int
	Height     int
	Services   services.Services
	Context    context.Context
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		return m, m.ResizeApp(msg)

	case constants.ProjectSelectedMsg:
		if err := m.Tasks.HandleProjectSelected(m.Context, msg); err != nil {
			fmt.Println(err.Error())
		}
		return m, nil

	case tea.KeyMsg:

		// Handle any global keybinds, ex: Quit
		if constants.KeyMatch(msg, constants.KeyQuit1, constants.KeyQuit2) {
			return m, tea.Quit
		}

		if constants.KeyMatch(msg, constants.KeyLeft1, constants.KeyLeft2) {
			return m.handleLeft()
		}
		if constants.KeyMatch(msg, constants.KeyRight1, constants.KeyRight2) {
			return m.handleRight()
		}

		// Handle keybinds specific to components
		if m.ActivePane == ProjectPane {
			return m.Projects.Update(msg)
		} else if m.ActivePane == TasksPane {
			return m.Tasks.Update(msg)
		}

	}

	return nil, nil
}

// View renders the complete UI
func (m Model) View() string {
	// Render components with appropriate styling based on active pane
	return ""
}

func (m Model) ResizeApp(msg tea.WindowSizeMsg) tea.Cmd {
	width := msg.Width / 2
	height := msg.Height
	pCmd := m.Projects.ResizeList(msg, width, height)
	tCmd := m.Tasks.ResizeList(msg, width, height)
	return tea.Batch(pCmd, tCmd)
}

func (m Model) handleLeft() (tea.Model, tea.Cmd) {
	if m.ActivePane == TasksPane {
		m.ActivePane = ProjectPane
	}
	return m, nil
}

func (m Model) handleRight() (tea.Model, tea.Cmd) {
	if m.ActivePane == ProjectPane {
		m.ActivePane = TasksPane
	}
	return m, nil
}
