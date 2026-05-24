package models

import "github.com/frozenkro/mcpsequencer/internal/projectsdb"

// Domain model for Task
// Contains structured list of dependencies
type Task struct {
	TaskId       int
	Name         string
	Description  string
	ProjectId    int
	Sort         int
	Status       Status
	Notes        string
	Dependencies []Dependency
}

type Status string

const (
	NotStarted Status = "Not Started"
	InProgress Status = "In Progress"
	Completed  Status = "Completed"
	Failed     Status = "Failed"
)

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

	return Task{
		TaskId:       int(dbTask.TaskID),
		Name:         dbTask.Name,
		Description:  dbTask.Description,
		ProjectId:    int(dbTask.ProjectID),
		Sort:         int(dbTask.Sort),
		Status:       status,
		Notes:        notes,
		Dependencies: deps,
	}
}
