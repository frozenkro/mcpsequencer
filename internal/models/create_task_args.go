package models

// CreateTaskArgs contains the parameters needed to create a new task.
// It includes basic task information and dependency relationships.
//
// When creating a project with tasks, Dependencies is a list of SortIDs
// When creating an individual task, Dependencies is a list of TaskIDs
type CreateTaskArgs struct {
	Name            string  `json:"name"`
	Description     string  `json:"description"`
	SortId          int     `json:"sortId"`
	Dependencies    []int   `json:"dependencies"`
	Owner           *string `json:"owner,omitempty"`
	ScheduledDate   *string `json:"scheduledDate,omitempty"`
	PhaseId         *int    `json:"phaseId,omitempty"`
	EstimatedHours  *int    `json:"estimatedHours,omitempty"`
}

// UpdateTaskArgs contains the parameters needed to update an existing task.
type UpdateTaskArgs struct {
	TaskId int
	Fields CreateTaskArgs
}
