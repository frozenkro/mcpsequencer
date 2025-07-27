package viewmodels

import (
	"encoding/json"
	"fmt"

	"github.com/frozenkro/mcpsequencer/internal/projectsdb"
)

type ProjectView struct {
	ProjectID int
	Name      string
}

func NewProjectView(p projectsdb.Project) ProjectView {
	return ProjectView{
		ProjectID: int(p.ProjectID),
		Name:      p.Name,
	}
}

func (p ProjectView) FilterValue() string {
	return p.Name
}

// Title returns the title for display in lists
func (p ProjectView) Title() string {
	return p.Name
}

// Description returns an empty string - could be customized in the future
func (p ProjectView) Description() string {
	return "Project ID: " + fmt.Sprintf("%d", p.ProjectID)
}

type TaskView struct {
	TaskID       int
	Name         string
	DescProp     string
	ProjectID    int
	Sort         int
	IsCompleted  bool
	IsInProgress bool
	Deps         []int
	Notes        string
}

func (t TaskView) FilterValue() string {
	return t.Title()
}

// Title returns the title for display in lists
func (t TaskView) Title() string {
	return t.Name
}

// Description returns information about the task status
func (t TaskView) Description() string {
	status := ""
	if t.IsCompleted {
		status = "Completed"
	} else if t.IsInProgress {
		status = "In Progress"
	} else {
		status = "Not Started"
	}
	return status
}

func NewTaskView(t projectsdb.Task) (TaskView, error) {
	isCompleted := false
	if t.IsCompleted == 1 {
		isCompleted = true
	}

	isInProgress := false
	if t.IsInProgress == 1 {
		isInProgress = true
	}

	deps := []int{}
	err := json.Unmarshal([]byte(t.DependenciesJson), &deps)

	return TaskView{
		TaskID:       int(t.TaskID),
		Name:         t.Name,
		DescProp:     t.Description,
		ProjectID:    int(t.ProjectID),
		Sort:         int(t.Sort),
		IsCompleted:  isCompleted,
		IsInProgress: isInProgress,
	}, err
}
