package tui

import (
	"context"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/frozenkro/mcpsequencer/internal/services"
	"github.com/frozenkro/mcpsequencer/internal/tui/components/projects"
	"github.com/frozenkro/mcpsequencer/internal/tui/components/tasks"
	"github.com/frozenkro/mcpsequencer/internal/tui/constants"
	"github.com/frozenkro/mcpsequencer/internal/tui/logger"
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
	logger.Logger.Printf("DEBUG UPDATE: %v", msg)
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		return m, m.ResizeApp(msg)

	case constants.ProjectSelectedMsg:
		if err := m.Tasks.HandleProjectSelected(m.Context, msg); err != nil {
			logger.Logger.Println(err.Error())
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
			teaModel, cmd := m.Projects.Update(msg)
			if projectsModel, ok := teaModel.(projects.Model); ok {
				m.Projects = projectsModel
			}
			return m, cmd
		} else if m.ActivePane == TasksPane {
			teaModel, cmd := m.Tasks.Update(msg)
			if tasksModel, ok := teaModel.(tasks.Model); ok {
				m.Tasks = tasksModel
			}
			return m, cmd
		}

	}

	return m, nil
}

// View renders the complete UI
func (m Model) View() string {
	// Render components with appropriate styling based on active pane
	var pStyle, tStyle lipgloss.Style

	if m.ActivePane == ProjectPane {
		pStyle = FocusedStyle.Copy().BorderForeground(lipgloss.Color("5"))
		tStyle = UnfocusedStyle
	} else {
		pStyle = UnfocusedStyle
		tStyle = FocusedStyle.Copy().BorderForeground(lipgloss.Color("5"))
	}

	pView := pStyle.Render(m.Projects.View())
	tView := tStyle.Render(m.Tasks.View())

	row := lipgloss.JoinHorizontal(lipgloss.Top, pView, tView)

	helpText := "\nNavigate: ←/h ↑/k ↓/j →/l • Select: Enter/Space • Quit: q or Ctrl+c"

	return AppStyle.Render(row + helpText)
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
