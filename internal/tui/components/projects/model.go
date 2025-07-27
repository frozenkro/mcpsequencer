package projects

import (
	"context"
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/frozenkro/mcpsequencer/internal/services"
	"github.com/frozenkro/mcpsequencer/internal/tui/constants"
	"github.com/frozenkro/mcpsequencer/internal/tui/logger"
	"github.com/frozenkro/mcpsequencer/internal/tui/viewmodels"
)

type Model struct {
	List     list.Model
	Selected *viewmodels.ProjectView
	svc      services.Services
}

func NewModel(svc services.Services, ctx context.Context, width int, height int) (Model, error) {
	projectsData, err := svc.GetProjects(ctx)
	if err != nil {
		return Model{}, err
	}

	var items []list.Item
	for _, project := range projectsData {
		items = append(items, viewmodels.NewProjectView(project))
	}

	delegate := createDelegate()

	projects := list.New(items, delegate, width, height)
	projects.Title = "Projects"

	return Model{
		List: projects,
		svc:  svc,
	}, nil
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

// Handle project selection
func (m *Model) SelectProject(ctx context.Context) (*viewmodels.ProjectView, []list.Item, error) {
	selectedItem := m.List.SelectedItem()
	var project *viewmodels.ProjectView

	if p, ok := selectedItem.(viewmodels.ProjectView); ok {
		project = &p
	} else if p, ok := selectedItem.(*viewmodels.ProjectView); ok {
		project = p
	} else {
		return nil, nil, nil
	}

	m.Selected = project

	// Load tasks for the selected project
	tasksData, err := m.svc.GetTasksByProject(ctx, int64(project.ProjectID))
	if err != nil {
		return project, nil, err
	}

	items := []list.Item{}
	for _, task := range tasksData {
		viewItem, err := viewmodels.NewTaskView(task)
		if err != nil {
			logger.Logger.Printf("WARN: Error during initialization of task view model for task '%v'\nError: '%v'\n", task.Name, err.Error())
		}

		items = append(items, viewItem)
	}

	return project, items, nil
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) View() string {
	// TODO add additional icons and other components here
	return m.List.View()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	logger.Logger.Println("Projects model update")
	switch msg := msg.(type) {

	case tea.KeyMsg:

		if constants.KeyMatch(msg, constants.KeyUp1, constants.KeyUp2, constants.KeyDown1, constants.KeyDown2) {
			logger.Logger.Println("Up or down")
			l, cmd := m.List.Update(msg)
			m.List = l
			logger.Logger.Println("Returning list update cmd")
			return m, cmd
		}

		if constants.KeyMatch(msg, constants.KeySelect1, constants.KeySelect2) {
			return m, func() tea.Msg {
				item := m.List.Items()[m.List.GlobalIndex()]
				if project, ok := item.(viewmodels.ProjectView); ok {
					m.Selected = &project
					logger.Logger.Println(fmt.Sprintf("project selected: '%v'", m.Selected.ProjectID))
					return constants.ProjectSelectedMsg{ProjectID: m.Selected.ProjectID}
				}
				logger.Logger.Println("Failed to parse Item to ProjectView")
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
