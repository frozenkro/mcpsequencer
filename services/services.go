package services

import (
	"context"
	"database/sql"

	"github.com/frozenkro/mcpsequencer/db"
	"github.com/frozenkro/mcpsequencer/projectsdb"
)

const TaskSortLast int64 = -1
const TaskSortFirst int64 = 0

type Services struct {
	Queries *projectsdb.Queries
}

func (s *Services) tryInit() {
	if s.Queries == nil {
		s.Queries = projectsdb.New(db.DB)
	}
}

func (s *Services) CreateProject(ctx context.Context, name string, tasks []string) error {
	s.tryInit()

	p, err := s.Queries.CreateProject(ctx, name)
	if err != nil {
		return err
	}

	for i, t := range tasks {
		task := projectsdb.CreateTaskParams{
			Description: t,
			Sort:        int64(i),
			ProjectID:   p.ProjectID,
		}
		if _, err := s.Queries.CreateTask(ctx, task); err != nil {
			// TODO either wait for all errs then return list,
			// or implement a rollback
			return err
		}
	}

	return nil
}

func (s *Services) RenameProject(ctx context.Context, project_id int64, name string) error {
	s.tryInit()

	params := projectsdb.UpdateProjectParams{
		ProjectID: project_id,
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

func (s *Services) AddTask(ctx context.Context, project_id int64, task string, sort int64) error {
	s.tryInit()

	if sort != TaskSortLast {
		tasks, err := s.Queries.GetTasksByProject(ctx, project_id)
		if err != nil {
			return err
		}

		for _, t := range tasks {
			if t.Sort >= sort {
				t.Sort = t.Sort + 1

				params := projectsdb.UpdateTaskSortParams{
					TaskID: t.TaskID,
					Sort:   t.Sort,
				}
				err := s.Queries.UpdateTaskSort(ctx, params)
				if err != nil {
					return err
				}
			}
		}
	}

	params := projectsdb.CreateTaskParams{
		Description: task,
		ProjectID:   project_id,
		Sort:        sort,
	}

	_, err := s.Queries.CreateTask(ctx, params)
	return err
}
