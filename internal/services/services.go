package services

import (
	"context"
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
type DependencyValidator interface {
	Validate([]models.Dependency) error
}
type TaskArrayTransformer interface {
	ParseFromJson(string, int) ([]projectsdb.Task, []models.Dependency, error)
	TaskIdMapFromTasks([]projectsdb.Task) models.SortIdTaskIdMap
}
type DependencyTransformer interface {
	FromInts([]int, int, models.DependencyDiscriminator) []models.Dependency
	FromDbRows([]projectsdb.Dependency) []models.Dependency
}

type Services struct {
	Queries               *projectsdb.Queries
	TaskArrayValidator    TaskArrayValidator
	DependencyValidator   DependencyValidator
	TaskArrayTransformer  TaskArrayTransformer
	DependencyTransformer DependencyTransformer
}

func NewServices() Services {
	s := Services{}
	s.Queries = projectsdb.New(db.DB)
	s.TaskArrayValidator = validators.TaskArrayValidator{}
	s.DependencyValidator = validators.DependencyValidator{}
	s.TaskArrayTransformer = transformers.TaskArrayTransformer{}
	s.DependencyTransformer = transformers.DependencyTransformer{}
	return s
}

func (s *Services) CreateProject(ctx context.Context, args models.CreateProjectArgs) error {
	params := projectsdb.CreateProjectParams{
		Name:         args.Name,
		Description:  args.Description,
		AbsolutePath: args.Directory,
	}
	p, err := s.Queries.CreateProject(ctx, params)
	if err != nil {
		return err
	}

	tasks, deps, err := s.TaskArrayTransformer.ParseFromJson(args.TasksJson, int(p.ProjectID))

	if err = s.TaskArrayValidator.Validate(tasks); err != nil {
		return err
	}
	if err = s.DependencyValidator.Validate(deps); err != nil {
		return err
	}

	sortToTaskIdMap := map[int64]int64{}
	for _, t := range tasks {
		task := projectsdb.CreateTaskParams{
			Name:        t.Name,
			Description: t.Description,
			Sort:        t.Sort,
			ProjectID:   p.ProjectID,
		}
		newTask, err := s.Queries.CreateTask(ctx, task)
		if err != nil {
			return fmt.Errorf("Unable to insert new task '%v'\nError: %w", task.Name, err)
		}

		sortToTaskIdMap[t.Sort] = newTask.TaskID
	}

	for _, d := range deps {
		depParams := projectsdb.AddDependencyForTaskParams{
			TaskID:    int64(sortToTaskIdMap[int64(d.Id)]),
			DependsOn: int64(sortToTaskIdMap[int64(d.DependsOn)]),
		}
		s.Queries.AddDependencyForTask(ctx, depParams)
	}

	return nil
}

func (s *Services) GetProjects(ctx context.Context) ([]projectsdb.Project, error) {
	return s.Queries.GetAllProjects(ctx)
}

func (s *Services) UpdateProject(ctx context.Context, args models.UpdateProjectArgs) error {
	params := projectsdb.UpdateProjectParams{
		ProjectID:    int64(args.ProjectId),
		Name:         args.Fields.Name,
		Description:  args.Fields.Description,
		AbsolutePath: args.Fields.Directory,
	}
	if _, err := s.Queries.UpdateProject(ctx, params); err != nil {
		return err
	}
	return nil
}

func (s *Services) DeleteProject(ctx context.Context, projectId int64) error {
	if err := s.Queries.DeleteProject(ctx, projectId); err != nil {
		return err
	}
	return nil
}

func (s *Services) GetTasksByProject(ctx context.Context, projectId int64) ([]models.Task, error) {
	taskRows, err := s.Queries.GetTasksByProject(ctx, projectId)
	if err != nil {
		return nil, fmt.Errorf("Error retrieving Tasks from DB: %w", err)
	}

	result := []models.Task{}
	for _, t := range taskRows {

		depRows, err := s.Queries.GetDependenciesForTask(ctx, t.TaskID)
		if err != nil {
			return nil, fmt.Errorf("Error retrieving Dependencies from DB: %w", err)
		}

		deps := s.DependencyTransformer.FromDbRows(depRows)

		result = append(result, models.NewTask(t, deps))
	}
	return result, nil
}

func (s *Services) AddTask(ctx context.Context, projectId int64, args models.CreateTaskArgs) error {

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
		Name:        args.Name,
		Sort:        int64(args.SortId),
		Description: args.Description,
		ProjectID:   projectId,
	}
	validationTasks := append(tasks, newTask)
	if err = s.TaskArrayValidator.Validate(validationTasks); err != nil {
		return fmt.Errorf("Error validating Task array: %w", err)
	}

	params := projectsdb.CreateTaskParams{
		Name:        newTask.Name,
		Description: newTask.Description,
		ProjectID:   projectId,
		Sort:        int64(args.SortId),
	}

	savedTask, err := s.Queries.CreateTask(ctx, params)
	if err != nil {
		return fmt.Errorf("Error saving new Task to DB: %w", err)
	}
	deps := s.DependencyTransformer.FromInts(args.Dependencies, int(savedTask.TaskID), models.TaskId)

	err = s.DependencyValidator.Validate(deps)
	if err != nil {
		return fmt.Errorf("Error validating dependency structure: %w", err)
	}

	for _, d := range deps {

		args := projectsdb.AddDependencyForTaskParams{
			TaskID:    savedTask.TaskID,
			DependsOn: int64(d.DependsOn),
		}

		err = s.Queries.AddDependencyForTask(ctx, args)
		if err != nil {
			return fmt.Errorf("Error saving dependency '{ %v, %v }' to database: %w", args.TaskID, args.DependsOn, err)
		}
	}

	return nil
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
