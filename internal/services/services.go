package services

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/frozenkro/mcpsequencer/internal/db"
	"github.com/frozenkro/mcpsequencer/internal/models"
	"github.com/frozenkro/mcpsequencer/internal/projectsdb"
	"github.com/frozenkro/mcpsequencer/internal/transformers"
	"github.com/frozenkro/mcpsequencer/internal/validators"
)

const TaskSortLast int64 = -1
const TaskSortFirst int64 = 0

type TaskState int

const (
	StatePending TaskState = iota
	StateInProgress
	StateComplete
)

type TaskArrayValidator interface {
	Validate([]projectsdb.Task) error
}
type TaskArrayTransformer interface {
	ParseFromJson(string, int) ([]projectsdb.Task, error)
}

type Services struct {
	Queries              *projectsdb.Queries
	TaskArrayValidator   TaskArrayValidator
	TaskArrayTransformer TaskArrayTransformer
}

func NewServices() Services {
	s := Services{}
	s.Queries = projectsdb.New(db.DB)
	s.TaskArrayValidator = validators.TaskArrayValidator{}
	s.TaskArrayTransformer = transformers.TaskArrayTransformer{}
	return s
}

func (s *Services) CreateProject(ctx context.Context, name string, tasksJson string) error {
	p, err := s.Queries.CreateProject(ctx, name)
	if err != nil {
		return err
	}

	tasks, err := s.TaskArrayTransformer.ParseFromJson(tasksJson, int(p.ProjectID))

	for _, t := range tasks {
		task := projectsdb.CreateTaskParams{
			Name:        t.Name,
			Description: t.Description,
			Sort:        t.Sort,
			ProjectID:   p.ProjectID,
		}
		if _, err := s.Queries.CreateTask(ctx, task); err != nil {
			return fmt.Errorf("Unable to insert new task '%v'\nError: %w", task.Name, err)
		}
	}

	return nil
}

func (s *Services) GetProjects(ctx context.Context) ([]projectsdb.Project, error) {
	return s.Queries.GetAllProjects(ctx)
}

func (s *Services) RenameProject(ctx context.Context, projectId int64, name string) error {
	params := projectsdb.UpdateProjectParams{
		ProjectID: projectId,
		Name:      name,
	}
	if _, err := s.Queries.UpdateProject(ctx, params); err != nil {
		return err
	}
	return nil
}

func (s *Services) DeleteProject(ctx context.Context, project_id int64) error {
	if err := s.Queries.DeleteProject(ctx, project_id); err != nil {
		return err
	}
	return nil
}

func (s *Services) GetTasksByProject(ctx context.Context, projectId int64) ([]projectsdb.Task, error) {
	return s.Queries.GetTasksByProject(ctx, projectId)
}

func (s *Services) AddTask(ctx context.Context, projectId int64, args models.CreateTaskArgs) error {
	newTaskDepsJson, err := json.Marshal(args.Dependencies)
	if err != nil {
		return fmt.Errorf("Unable to parse Dependencies array\nError: %w", err)
	}

	tasks, err := s.Queries.GetTasksByProject(ctx, projectId)
	if err != nil {
		return err
	}

	if int64(args.SortId) != TaskSortLast {
		for _, t := range tasks {
			if t.Sort >= int64(args.SortId) {
				t.Sort = t.Sort + 1

				params := projectsdb.UpdateTaskSortParams{
					TaskID: t.TaskID,
					Sort:   t.Sort,
				}
				err := s.Queries.UpdateTaskSort(ctx, params)
				if err != nil {
					return fmt.Errorf("Unable to update Sort ID to '%v' for TaskID '%v'\nError: %w", t.Sort, t.TaskID, err)
				}
			}
		}
	}

	newTask := projectsdb.Task{
		Name:             args.Name,
		Sort:             int64(args.SortId),
		Description:      args.Description,
		DependenciesJson: string(newTaskDepsJson),
		ProjectID:        projectId,
	}
	tasks = append(tasks, newTask)
	if err = s.TaskArrayValidator.Validate(tasks); err != nil {
		return err
	}

	params := projectsdb.CreateTaskParams{
		Name:             newTask.Name,
		Description:      newTask.Description,
		ProjectID:        projectId,
		Sort:             int64(args.SortId),
		DependenciesJson: newTask.DependenciesJson,
	}

	_, err = s.Queries.CreateTask(ctx, params)
	return err
}

func (s *Services) UpdateTaskState(ctx context.Context, taskId int64, state TaskState) error {
	inProgress := 0
	complete := 0
	if state == StateInProgress {
		inProgress = 1
	}
	if state == StateComplete {
		complete = 1
	}

	params := projectsdb.UpdateTaskStatusParams{
		TaskID:       taskId,
		IsInProgress: int64(inProgress),
		IsCompleted:  int64(complete),
	}
	_, err := s.Queries.UpdateTaskStatus(ctx, params)
	return err
}
