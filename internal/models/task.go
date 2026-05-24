package models

import "github.com/frozenkro/mcpsequencer/internal/projectsdb"

// Domain model for Task
// Contains structured list of dependencies
type Task struct {
	TaskId         int
	Name           string
	Description    string
	ProjectId      int
	Sort           int
	Status         Status
	Notes          string
	Owner          Owner
	ScheduledDate  *string
	PhaseId        *int
	BlockerText    *string
	BlockedAt      *string
	EstimatedHours *int
	Dependencies   []Dependency
}

type Status string

const (
	NotStarted Status = "Not Started"
	InProgress Status = "In Progress"
	Completed  Status = "Completed"
	Failed     Status = "Failed"
)

type Owner string

const (
	User      Owner = "USER"
	AiAgent   Owner = "AI_AGENT"
	Collabor  Owner = "COLLAB"
)

// Phase represents a milestone date range within a project
type Phase struct {
	PhaseId     int
	ProjectId   int
	Name        string
	Description string
	StartDate   string
	EndDate     string
	Sort        int
}

// TaskNote represents a chronological daily update attached to a task
type TaskNote struct {
	NoteID    int
	TaskID    int
	Note      string
	CreatedAt string
}

// ExitCriterion is a checklist item belonging to a phase
type ExitCriterion struct {
	CriterionID int
	PhaseID     int
	Description string
	IsCompleted bool
	Sort        int
}

func NewTask(dbTask projectsdb.Task, deps []Dependency) Task {
	status := NotStarted
	if dbTask.IsCompleted != int64(0) {
		status = Completed
	} else if dbTask.IsInProgress != int64(0) {
		status = InProgress
	}

	notes, ok := dbTask.Notes.(string)
	if !ok {
		notes = ""
	}

	owner := Owner(dbTask.Owner)
	if owner == "" {
		owner = User
	}

	var scheduledDate *string
	if sd, ok := dbTask.ScheduledDate.(string); ok && sd != "" {
		scheduledDate = &sd
	}

	var phaseID *int
	if pid, ok := dbTask.PhaseID.(int64); ok {
		v := int(pid)
		phaseID = &v
	}

	var blockerText *string
	if bt, ok := dbTask.BlockerText.(string); ok && bt != "" {
		blockerText = &bt
	}

	var blockedAt *string
	if ba, ok := dbTask.BlockedAt.(string); ok && ba != "" {
		blockedAt = &ba
	}

	var estimatedHours *int
	if eh, ok := dbTask.EstimatedHours.(int64); ok {
		v := int(eh)
		estimatedHours = &v
	}

	return Task{
		TaskId:         int(dbTask.TaskID),
		Name:           dbTask.Name,
		Description:    dbTask.Description,
		ProjectId:      int(dbTask.ProjectID),
		Sort:           int(dbTask.Sort),
		Status:         status,
		Notes:          notes,
		Owner:          owner,
		ScheduledDate:  scheduledDate,
		PhaseId:        phaseID,
		BlockerText:    blockerText,
		BlockedAt:      blockedAt,
		EstimatedHours: estimatedHours,
		Dependencies:   deps,
	}
}
