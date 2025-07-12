package projects

import (
	"context"

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
func (m *Model) SelectProject(svc services.Services, ctx context.Context) (*viewmodels.ProjectView, []list.Item, error) {
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
	tasksData, err := svc.GetTasksByProject(ctx, int64(project.ProjectID))
	if err != nil {
		return project, nil, err
	}

	items := []list.Item{}
	for _, task := range tasksData {
		items = append(items, viewmodels.NewTaskView(task))
	}

	return project, items, nil
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) View() string {
	return ""
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
				item := m.List.Items()[m.List.GlobalIndex()]
				if project, ok := item.(viewmodels.ProjectView); ok {
					m.Selected = &project
					return constants.ProjectSelectedMsg{ProjectID: m.Selected.ProjectID}
				}
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
