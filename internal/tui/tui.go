package tui

// import (
// 	"context"
// 	"fmt"

// 	"github.com/charmbracelet/bubbles/list"
// 	tea "github.com/charmbracelet/bubbletea"
// 	"github.com/charmbracelet/lipgloss"

// 	"github.com/frozenkro/mcpsequencer/internal/services"
// 	"github.com/frozenkro/mcpsequencer/internal/tui/viewmodels"
// )

// var svc services.Services
// var ctx context.Context

// type model struct {
// 	projects      list.Model
// 	activeProject *viewmodels.ProjectView
// 	tasks         list.Model
// 	cursor        int
// 	activePane    ActivePane
// 	width         int
// 	height        int
// }

// func (m *model) resizeList(width, height int) {
// 	m.width = width
// 	m.height = height

// 	// Calculate dimensions for side-by-side layout
// 	halfWidth := width / 2

// 	// Update list dimensions
// 	m.projects.SetSize(halfWidth-4, height-8) // Allow for borders and padding
// 	m.tasks.SetSize(halfWidth-4, height-8)    // Allow for borders and padding
// }

// func InitialModel() (model, error) {
// 	ctx = context.Background()
// 	svc = services.Services{}
// 	projectsData, err := svc.GetProjects(ctx)
// 	if err != nil {
// 		return model{}, err
// 	}

// 	var items []list.Item
// 	for _, project := range projectsData {
// 		items = append(items, viewmodels.NewProjectView(project))
// 	}

// 	// Create custom delegates with better styling
// 	projectDelegate := list.NewDefaultDelegate()
// 	projectDelegate.ShowDescription = true
// 	projectDelegate.Styles.SelectedTitle = projectDelegate.Styles.SelectedTitle.Foreground(lipgloss.Color("5")).Bold(true)
// 	projectDelegate.Styles.SelectedDesc = projectDelegate.Styles.SelectedDesc.Foreground(lipgloss.Color("12"))
// 	projectDelegate.Styles.NormalTitle = projectDelegate.Styles.NormalTitle.Foreground(lipgloss.Color("7"))
// 	projectDelegate.Styles.NormalDesc = projectDelegate.Styles.NormalDesc.Foreground(lipgloss.Color("8"))

// 	taskDelegate := list.NewDefaultDelegate()
// 	taskDelegate.ShowDescription = true
// 	taskDelegate.Styles.SelectedTitle = taskDelegate.Styles.SelectedTitle.Foreground(lipgloss.Color("5")).Bold(true)
// 	taskDelegate.Styles.SelectedDesc = taskDelegate.Styles.SelectedDesc.Foreground(lipgloss.Color("12"))
// 	taskDelegate.Styles.NormalTitle = taskDelegate.Styles.NormalTitle.Foreground(lipgloss.Color("7"))
// 	taskDelegate.Styles.NormalDesc = taskDelegate.Styles.NormalDesc.Foreground(lipgloss.Color("8"))

// 	projects := list.New(items, projectDelegate, 0, 0)
// 	projects.Title = "Projects"

// 	tasks := list.New([]list.Item{}, taskDelegate, 0, 0)
// 	tasks.Title = "Tasks"

// 	// Set default dimensions - these will be adjusted by WindowSizeMsg
// 	defaultWidth := 100
// 	defaultHeight := 30

// 	// Create model
// 	m := model{
// 		projects:   projects,
// 		activePane: ProjectPane,
// 		tasks:      tasks,
// 		width:      defaultWidth,
// 		height:     defaultHeight,
// 	}

// 	// Initialize list sizes
// 	m.resizeList(defaultWidth, defaultHeight)

// 	return m, nil
// }

// func (m model) Init() tea.Cmd {
// 	// Just return `nil`, which means "no I/O right now, please."
// 	return nil
// }

// func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
// 	switch msg := msg.(type) {
// 	case tea.WindowSizeMsg:
// 		m.resizeList(msg.Width, msg.Height)
// 		var projectsCmd, tasksCmd tea.Cmd

// 		// Make sure each list gets the resize message
// 		m.projects, projectsCmd = m.projects.Update(msg)
// 		m.tasks, tasksCmd = m.tasks.Update(msg)

// 		return m, tea.Batch(projectsCmd, tasksCmd)

// 	case tea.KeyMsg:

// 		switch msg.String() {

// 		case "ctrl+c", "q":
// 			return m, tea.Quit

// 		case "up", "k":
// 			if m.activePane == ProjectPane {
// 				var cmd tea.Cmd
// 				m.projects, cmd = m.projects.Update(msg)
// 				return m, cmd
// 			} else {
// 				var cmd tea.Cmd
// 				m.tasks, cmd = m.tasks.Update(msg)
// 				return m, cmd
// 			}

// 		case "left", "h":
// 			return m.handleLeft()

// 		case "right", "l":
// 			return m.handleRight()

// 		case "down", "j":
// 			if m.activePane == ProjectPane {
// 				var cmd tea.Cmd
// 				m.projects, cmd = m.projects.Update(msg)
// 				return m, cmd
// 			} else {
// 				var cmd tea.Cmd
// 				m.tasks, cmd = m.tasks.Update(msg)
// 				return m, cmd
// 			}

// 		case "enter", " ":
// 			return m.handleSelect()
// 		}
// 	}

// 	// Return the updated model to the Bubble Tea runtime for processing.
// 	// Note that we're not returning a command.
// 	return m, nil
// }

// func (m model) View() string {
// 	// Apply styles based on active pane
// 	var projectsStyle, tasksStyle lipgloss.Style

// 	if m.activePane == ProjectPane {
// 		projectsStyle = FocusedStyle.Copy().BorderForeground(lipgloss.Color("5"))
// 		tasksStyle = UnfocusedStyle
// 	} else {
// 		projectsStyle = UnfocusedStyle
// 		tasksStyle = FocusedStyle.Copy().BorderForeground(lipgloss.Color("5"))
// 	}

// 	// Generate views for both lists
// 	projectsView := projectsStyle.Render(m.projects.View())
// 	tasksView := tasksStyle.Render(m.tasks.View())

// 	// Join them horizontally
// 	row := lipgloss.JoinHorizontal(lipgloss.Top, projectsView, tasksView)

// 	// Add some help text at the bottom
// 	helpText := "\nNavigate: ←/h ↑/k ↓/j →/l • Select: Enter/Space • Quit: q or Ctrl+c"

// 	// Return the final view
// 	return AppStyle.Render(row + helpText)
// }

// func (m model) activeListLen() int {
// 	if m.activePane == ProjectPane {
// 		return len(m.projects.Items())
// 	}
// 	if m.activePane == TasksPane {
// 		return len(m.tasks.Items())
// 	}
// 	return 0
// }

// func (m model) handleSelect() (tea.Model, tea.Cmd) {
// 	if m.activePane == ProjectPane {
// 		selectedItem := m.projects.SelectedItem()
// 		if project, ok := selectedItem.(viewmodels.ProjectView); ok {
// 			m.activeProject = &project
// 		} else if projectPtr, ok := selectedItem.(*viewmodels.ProjectView); ok {
// 			m.activeProject = projectPtr
// 		}

// 		// Load tasks for the selected project
// 		tasksData, err := svc.GetTasksByProject(ctx, int64(m.activeProject.ProjectID))
// 		if err != nil {
// 			fmt.Printf("Error retrieving tasks for project: %v\n", err.Error())
// 			return m, nil
// 		}

// 		items := []list.Item{}
// 		for _, task := range tasksData {
// 			items = append(items, viewmodels.NewTaskView(task))
// 		}
// 		m.tasks.SetItems(items)

// 		// Switch to tasks pane
// 		m.activePane = TasksPane
// 		return m, nil
// 	}

// 	// For task pane, you could implement actions when tasks are selected
// 	// For now, just return the model
// 	return m, nil
// }

// func (m model) handleLeft() (tea.Model, tea.Cmd) {
// 	if m.activePane == TasksPane {
// 		m.activePane = ProjectPane
// 	}
// 	return m, nil
// }

// func (m model) handleRight() (tea.Model, tea.Cmd) {
// 	if m.activePane == ProjectPane {
// 		m.activePane = TasksPane
// 	}
// 	return m, nil
// }
