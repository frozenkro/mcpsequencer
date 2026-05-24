package viewmodels

import (
	"fmt"

	"github.com/frozenkro/mcpsequencer/internal/models"
	"github.com/frozenkro/mcpsequencer/internal/projectsdb"
)

type ProjectItem struct {
	ProjectID int
	Name      string
}

func NewProjectItem(p projectsdb.Project) ProjectItem {
	return ProjectItem{
		ProjectID: int(p.ProjectID),
		Name:      p.Name,
	}
}

func (p ProjectItem) FilterValue() string {
	return p.Name
}

// Title returns the title for display in lists
func (p ProjectItem) Title() string {
	return p.Name
}

// Description returns an empty string - could be customized in the future
func (p ProjectItem) Description() string {
	return "Project ID: " + fmt.Sprintf("%d", p.ProjectID)
}

type PhaseItem struct {
	PhaseID     int
	Name        string
	StartDate   string
	EndDate     string
	ExitDone    int
	ExitTotal   int
}

func NewPhaseItem(phase models.Phase, exitDone, exitTotal int) PhaseItem {
	return PhaseItem{
		PhaseID:   phase.PhaseId,
		Name:      phase.Name,
		StartDate: phase.StartDate,
		EndDate:   phase.EndDate,
		ExitDone:  exitDone,
		ExitTotal: exitTotal,
	}
}

func (p PhaseItem) FilterValue() string { return p.Name }
func (p PhaseItem) Title() string     { return p.Name }
func (p PhaseItem) Description() string {
	return fmt.Sprintf("%s → %s | exit criteria: %d/%d", p.StartDate, p.EndDate, p.ExitDone, p.ExitTotal)
}

type TaskItem struct {
	TaskID         int
	Name           string
	DescProp       string
	ProjectID      int
	Sort           int
	Status         models.Status
	Deps           []int
	Notes          string
	Owner          models.Owner
	ScheduledDate  *string
	PhaseId        *int
	BlockerText    *string
	BlockedAt      *string
	EstimatedHours *int
}

func (t TaskItem) FilterValue() string {
	return t.Title()
}

// Title returns the title for display in lists
func (t TaskItem) Title() string {
	return t.Name
}

func (t TaskItem) Description() string {
	return string(t.Status)
}

func NewTaskItem(t models.Task) (TaskItem, error) {

	deps := []int{}
	for _, d := range t.Dependencies {
		if d.Discriminator != models.TaskId {
			return TaskItem{}, fmt.Errorf("Received Dependency with SortIds, expecting TaskIds")
		}
		deps = append(deps, d.DependsOn)
	}

	return TaskItem{
		TaskID:         int(t.TaskId),
		Name:           t.Name,
		DescProp:       t.Description,
		ProjectID:      int(t.ProjectId),
		Sort:           int(t.Sort),
		Status:         t.Status,
		Deps:           deps,
		Owner:          t.Owner,
		ScheduledDate:  t.ScheduledDate,
		PhaseId:        t.PhaseId,
		BlockerText:    t.BlockerText,
		BlockedAt:      t.BlockedAt,
		EstimatedHours: t.EstimatedHours,
	}, nil
}
