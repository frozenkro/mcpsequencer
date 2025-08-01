package taskdetail

import (
	"context"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/frozenkro/mcpsequencer/internal/models"
	"github.com/frozenkro/mcpsequencer/internal/services"
	"github.com/frozenkro/mcpsequencer/internal/tui/logger"
	"github.com/frozenkro/mcpsequencer/internal/tui/viewmodels"
)

type Mode int

const (
	ViewMode Mode = iota
	EditMode
)

type Model struct {
	Task         *viewmodels.TaskItem
	Mode         Mode
	Width        int
	Height       int
	svc          services.Services
	nameInput    textinput.Model
	descTextarea textarea.Model
	focusedField int // 0: name, 1: description
}

func NewModel(svc services.Services, width, height int) Model {
	nameInput := textinput.New()
	nameInput.Placeholder = "Task name..."
	nameInput.CharLimit = 100

	descTextarea := textarea.New()
	descTextarea.Placeholder = "Task description..."
	descTextarea.SetWidth(width - 4)
	descTextarea.SetHeight(5)

	return Model{
		svc:          svc,
		Width:        width,
		Height:       height,
		Mode:         ViewMode,
		nameInput:    nameInput,
		descTextarea: descTextarea,
		focusedField: 0,
	}
}

func (m *Model) HandleTaskSelected(task *viewmodels.TaskItem) {
	m.Task = task
	if task != nil {
		m.nameInput.SetValue(task.Name)
		m.descTextarea.SetValue(task.DescProp)
	} else {
		m.nameInput.SetValue("")
		m.descTextarea.SetValue("")
	}
	m.Mode = ViewMode
}

func (m *Model) enterEditMode() {
	if m.Task == nil {
		return
	}
	m.Mode = EditMode
	m.nameInput.Focus()
	m.focusedField = 0
}

func (m *Model) exitEditMode() {
	m.Mode = ViewMode
	m.nameInput.Blur()
	m.descTextarea.Blur()
}

func (m *Model) saveChanges(ctx context.Context) error {
	if m.Task == nil {
		return fmt.Errorf("no task selected")
	}

	// Here you would call the service to update the task
	// For now, just update the local model
	m.Task.Name = m.nameInput.Value()
	m.Task.DescProp = m.descTextarea.Value()

	logger.Logger.Printf("Saving task changes: %s", m.Task.Name)
	return nil
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) View() string {
	if m.Task == nil {
		return m.renderEmpty()
	}

	if m.Mode == EditMode {
		return m.renderEditMode()
	}
	return m.renderViewMode()
}

func (m Model) renderEmpty() string {
	style := lipgloss.NewStyle().
		Padding(2).
		Foreground(lipgloss.Color("8"))

	return style.Render("Select a task to view details")
}

func (m Model) renderViewMode() string {
	var content strings.Builder

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("5")).
		MarginBottom(1)

	labelStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("6"))

	valueStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("7"))

	statusStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("2"))

	if m.Task.Status == models.Completed {
		statusStyle = statusStyle.Foreground(lipgloss.Color("2")) // Green
	} else if m.Task.Status == models.InProgress {
		statusStyle = statusStyle.Foreground(lipgloss.Color("3")) // Yellow
	} else if m.Task.Status == models.Failed {
		statusStyle = statusStyle.Foreground(lipgloss.Color("1")) // Red
	} else {
		statusStyle = statusStyle.Foreground(lipgloss.Color("8")) // Gray
	}

	content.WriteString(titleStyle.Render("Task Details"))
	content.WriteString("\n\n")

	content.WriteString(labelStyle.Render("Name: "))
	content.WriteString(valueStyle.Render(m.Task.Name))
	content.WriteString("\n\n")

	content.WriteString(labelStyle.Render("Status: "))
	content.WriteString(statusStyle.Render(string(m.Task.Status)))
	content.WriteString("\n\n")

	content.WriteString(labelStyle.Render("Description:"))
	content.WriteString("\n")
	if m.Task.DescProp != "" {
		content.WriteString(valueStyle.Render(m.Task.DescProp))
	} else {
		content.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("8")).Render("No description"))
	}
	content.WriteString("\n\n")

	content.WriteString(labelStyle.Render("Task ID: "))
	content.WriteString(valueStyle.Render(fmt.Sprintf("%d", m.Task.TaskID)))
	content.WriteString("\n")

	content.WriteString(labelStyle.Render("Sort Order: "))
	content.WriteString(valueStyle.Render(fmt.Sprintf("%d", m.Task.Sort)))
	content.WriteString("\n\n")

	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("8")).
		Italic(true)

	content.WriteString(helpStyle.Render("Press 'e' to edit"))

	return lipgloss.NewStyle().
		Padding(1).
		Width(m.Width - 2).
		Height(m.Height - 2).
		Render(content.String())
}

func (m Model) renderEditMode() string {
	var content strings.Builder

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("5")).
		MarginBottom(1)

	labelStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("6"))

	content.WriteString(titleStyle.Render("Edit Task"))
	content.WriteString("\n\n")

	content.WriteString(labelStyle.Render("Name:"))
	content.WriteString("\n")
	content.WriteString(m.nameInput.View())
	content.WriteString("\n\n")

	content.WriteString(labelStyle.Render("Description:"))
	content.WriteString("\n")
	content.WriteString(m.descTextarea.View())
	content.WriteString("\n\n")

	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("8")).
		Italic(true)

	content.WriteString(helpStyle.Render("Tab: Next field • Enter: Save • Esc: Cancel"))

	return lipgloss.NewStyle().
		Padding(1).
		Width(m.Width - 2).
		Height(m.Height - 2).
		Render(content.String())
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.Mode == EditMode {
			return m.handleEditModeKeys(msg)
		}
		return m.handleViewModeKeys(msg)
	}

	if m.Mode == EditMode {
		if m.focusedField == 0 {
			m.nameInput, cmd = m.nameInput.Update(msg)
			cmds = append(cmds, cmd)
		} else {
			m.descTextarea, cmd = m.descTextarea.Update(msg)
			cmds = append(cmds, cmd)
		}
	}

	return m, tea.Batch(cmds...)
}

func (m Model) handleViewModeKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "e":
		if m.Task != nil {
			m.enterEditMode()
		}
	}
	return m, nil
}

func (m Model) handleEditModeKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg.String() {
	case "esc":
		m.exitEditMode()
		return m, nil
	case "enter":
		if m.focusedField == 0 { // Only save on enter in name field
			// Save changes
			ctx := context.Background()
			if err := m.saveChanges(ctx); err != nil {
				logger.Logger.Printf("Error saving task: %v", err)
			}
			m.exitEditMode()
			return m, nil
		}
	case "tab":
		if m.focusedField == 0 {
			m.nameInput.Blur()
			m.descTextarea.Focus()
			m.focusedField = 1
		} else {
			m.descTextarea.Blur()
			m.nameInput.Focus()
			m.focusedField = 0
		}
		return m, nil
	}

	// Update the focused field
	if m.focusedField == 0 {
		m.nameInput, cmd = m.nameInput.Update(msg)
		cmds = append(cmds, cmd)
	} else {
		m.descTextarea, cmd = m.descTextarea.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m Model) ResizeView(width, height int) {
	m.Width = width
	m.Height = height
	m.descTextarea.SetWidth(width - 4)
}
