package models

type CreateProjectArgs struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Directory   string `json:"directory"`
	TasksJson   string `json:"tasksJson"`
}

type UpdateProjectArgs struct {
	ProjectId int
	Fields    CreateProjectArgs
}
