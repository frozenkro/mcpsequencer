package tasks

import (
	"context"
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/frozenkro/mcpsequencer/internal/services"
	"github.com/frozenkro/mcpsequencer/internal/tui/constants"
	"github.com/frozenkro/mcpsequencer/internal/tui/viewmodels"
)

type Model struct {
	List     list.Model
	Selected *viewmodels.ProjectView
	svc      services.Services
}

func NewModel(svc services.Services, width, height int) Model {

	list := list.New([]list.Item{}, createDelegate(), width, height)
	list.Title = "Tasks"
	return Model{
		List: list,
		svc:  svc,
	}
}

func createDelegate() list.DefaultDelegate {
	delegate := list.NewDefaultDelegate()
	delegate.ShowDescription = true
	delegate.Styles.SelectedTitle = delegate.Styles.SelectedTitle.Foreground(lipgloss.Color("5")).Bold(true)
	delegate.Styles.SelectedDesc = delegate.Styles.SelectedDesc.Foreground(lipgloss.Color("12"))
	delegate.Styles.NormalTitle = delegate.Styles.NormalTitle.Foreground(lipgloss.Color("7"))
	delegate.Styles.NormalDesc = delegate.Styles.NormalDesc.Foreground(lipgloss.Color("8"))
	return delegate
}

func (m *Model) HandleProjectSelected(ctx context.Context, msg constants.ProjectSelectedMsg) error {
	tasksData, err := m.svc.GetTasksByProject(ctx, int64(msg.ProjectID))
	if err != nil {
		return fmt.Errorf("Error retrieving tasks data for project %v: %w\n", msg.ProjectID, err)
	}

	l := []list.Item{}
	for _, t := range tasksData {
		taskView := viewmodels.NewTaskView(t)
		l = append(l, taskView)
	}
	m.List.SetItems(l)
	return nil
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) View() string {
	// TODO add additional icons and other components here
	return m.List.View()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:

		if constants.KeyMatch(msg, constants.KeyUp1, constants.KeyUp2, constants.KeyDown1, constants.KeyDown2) {
			l, cmd := m.List.Update(msg)
			m.List = l
			return m, cmd
		}
	}

	return m, nil
}

func (m Model) ResizeList(msg tea.WindowSizeMsg, width, height int) tea.Cmd {
	m.List.SetWidth(width)
	m.List.SetHeight(height)

	l, cmd := m.List.Update(msg)
	m.List = l
	return cmd
}
