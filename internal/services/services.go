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

type Services struct {
	Queries *projectsdb.Queries
}

func (s *Services) tryInit() {
	if s.Queries == nil {
		s.Queries = projectsdb.New(db.DB)
	}
}

func (s *Services) CreateProject(ctx context.Context, name string, tasksJson []string) error {
	s.tryInit()

	p, err := s.Queries.CreateProject(ctx, name)
	if err != nil {
		return err
	}

	tasks, err := transformers.ParseTasksArray(tasksJson, int(p.ProjectID))

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
	s.tryInit()

	return s.Queries.GetAllProjects(ctx)
}

func (s *Services) RenameProject(ctx context.Context, projectId int64, name string) error {
	s.tryInit()

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
	s.tryInit()

	if err := s.Queries.DeleteProject(ctx, project_id); err != nil {
		return err
	}
	return nil
}

func (s *Services) GetTasksByProject(ctx context.Context, projectId int64) ([]projectsdb.Task, error) {
	s.tryInit()

	return s.Queries.GetTasksByProject(ctx, projectId)
}

func (s *Services) AddTask(ctx context.Context, projectId int64, args models.CreateTaskArgs) error {
	s.tryInit()

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
	validators.ValidateTasksArray(tasks)

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
	s.tryInit()

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
