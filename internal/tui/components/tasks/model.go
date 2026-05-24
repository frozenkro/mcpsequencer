package tasks

import (
	"context"
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/frozenkro/mcpsequencer/internal/services"
	"github.com/frozenkro/mcpsequencer/internal/tui/constants"
	"github.com/frozenkro/mcpsequencer/internal/tui/logger"
	"github.com/frozenkro/mcpsequencer/internal/tui/viewmodels"
)

type Model struct {
	List     list.Model
	Selected *viewmodels.ProjectItem
	svc      services.Services
}

func NewModel(svc services.Services, width, height int) Model {

	list := list.New([]list.Item{}, createDelegate(), width, height)
	list.Title = "Tasks"
	list.SetShowHelp(false)
	return Model{
		List: list,
		svc:  svc,
	}
}

func createDelegate() TaskItemDelegate {
	return TaskItemDelegate{}
}

// func createDelegate() list.DefaultDelegate {
// 	delegate := list.NewDefaultDelegate()
// 	delegate.ShowDescription = true
// 	delegate.Styles.SelectedTitle = delegate.Styles.SelectedTitle.Foreground(lipgloss.Color("5")).Bold(true)
// 	delegate.Styles.SelectedDesc = delegate.Styles.SelectedDesc.Foreground(lipgloss.Color("12"))
// 	delegate.Styles.NormalTitle = delegate.Styles.NormalTitle.Foreground(lipgloss.Color("7"))
// 	delegate.Styles.NormalDesc = delegate.Styles.NormalDesc.Foreground(lipgloss.Color("8"))
// 	return delegate
// }

func (m *Model) HandleProjectSelected(ctx context.Context, msg constants.ProjectSelectedMsg) error {
	tasksData, err := m.svc.GetTasksByProject(ctx, int64(msg.ProjectID))
	if err != nil {
		return fmt.Errorf("Error retrieving tasks data for project %v: %w\n", msg.ProjectID, err)
	}

	l := []list.Item{}
	for _, t := range tasksData {
		taskItem, err := viewmodels.NewTaskItem(t)
		if err != nil {
			logger.Logger.Printf("WARN: Error during initialization of task view model for task '%v'\nError: '%v'\n", t.Name, err.Error())
		}
		l = append(l, taskItem)
	}
	m.List.SetItems(l)
	return nil
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) View() string {
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

		if constants.KeyMatch(msg, constants.KeySelect1, constants.KeySelect2) {
			return m, func() tea.Msg {
				if len(m.List.Items()) == 0 {
					return nil
				}
				item := m.List.Items()[m.List.Index()]
				if task, ok := item.(viewmodels.TaskItem); ok {
					logger.Logger.Printf("task selected: '%v'", task.TaskID)
					return constants.TaskSelectedMsg{TaskID: task.TaskID}
				}
				logger.Logger.Println("Failed to parse Item to TaskItem")
				return nil
			}
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
