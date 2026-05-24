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
	Validate([]models.Dependency, []int) error
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
	if err != nil {
		return err
	}

	if err = s.TaskArrayValidator.Validate(tasks); err != nil {
		return err
	}

	taskIds := []int{}
	for _, t := range tasks {
		taskIds = append(taskIds, int(t.Sort))
	}
	if err = s.DependencyValidator.Validate(deps, taskIds); err != nil {
		return err
	}

	sortToTaskIdMap := map[int64]int64{}
	for _, t := range tasks {
		var sched interface{}
		if sd := t.ScheduledDate; sd != nil {
			if str, ok := sd.(string); ok {
				sched = str
			}
		}
		var phaseId interface{}
		if pid, ok := t.PhaseID.(int64); ok {
			phaseId = pid
		}
		var estHours interface{}
		if eh, ok := t.EstimatedHours.(int64); ok {
			estHours = eh
		}

		task := projectsdb.CreateTaskParams{
			Name:           t.Name,
			Description:    t.Description,
			Sort:           t.Sort,
			ProjectID:      p.ProjectID,
			Owner:          t.Owner,
			ScheduledDate:  sched,
			PhaseID:        phaseId,
			EstimatedHours: estHours,
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

	owner := string(models.User)
	if args.Owner != nil {
		owner = *args.Owner
	}

	var sched interface{}
	if args.ScheduledDate != nil {
		sched = *args.ScheduledDate
	}
	var phaseId interface{}
	if args.PhaseId != nil {
		phaseId = int64(*args.PhaseId)
	}
	var estHours interface{}
	if args.EstimatedHours != nil {
		estHours = int64(*args.EstimatedHours)
	}

	newTask := projectsdb.Task{
		Name:           args.Name,
		Sort:           int64(args.SortId),
		Description:    args.Description,
		ProjectID:      projectId,
		Owner:          owner,
		ScheduledDate:  sched,
		PhaseID:        phaseId,
		EstimatedHours: estHours,
	}
	validationTasks := append(tasks, newTask)
	if err = s.TaskArrayValidator.Validate(validationTasks); err != nil {
		return fmt.Errorf("Error validating Task array: %w", err)
	}

	params := projectsdb.CreateTaskParams{
		Name:           newTask.Name,
		Description:    newTask.Description,
		ProjectID:      projectId,
		Sort:           int64(args.SortId),
		Owner:          owner,
		ScheduledDate:  sched,
		PhaseID:        phaseId,
		EstimatedHours: estHours,
	}

	savedTask, err := s.Queries.CreateTask(ctx, params)
	if err != nil {
		return fmt.Errorf("Error saving new Task to DB: %w", err)
	}
	deps := s.DependencyTransformer.FromInts(args.Dependencies, int(savedTask.TaskID), models.TaskId)

	taskIds := []int{int(savedTask.TaskID)}
	for _, t := range tasks {
		taskIds = append(taskIds, int(t.TaskID))
	}

	err = s.DependencyValidator.Validate(deps, taskIds)
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

// --- Phase Management ---

func (s *Services) CreatePhase(ctx context.Context, p models.Phase) (models.Phase, error) {
	params := projectsdb.CreatePhaseParams{
		ProjectID:   int64(p.ProjectId),
		Name:          p.Name,
		Description:   p.Description,
		StartDate:     p.StartDate,
		EndDate:       p.EndDate,
		Sort:          int64(p.Sort),
	}
	dbPhase, err := s.Queries.CreatePhase(ctx, params)
	if err != nil {
		return models.Phase{}, err
	}
	return models.Phase{
		PhaseId:     int(dbPhase.PhaseID),
		ProjectId:   int(dbPhase.ProjectID),
		Name:        dbPhase.Name,
		Description: fmt.Sprintf("%v", dbPhase.Description),
		StartDate:   dbPhase.StartDate,
		EndDate:     dbPhase.EndDate,
		Sort:        int(dbPhase.Sort),
	}, nil
}

func (s *Services) GetPhases(ctx context.Context, projectId int64) ([]models.Phase, error) {
	rows, err := s.Queries.GetPhasesForProject(ctx, projectId)
	if err != nil {
		return nil, err
	}
	result := make([]models.Phase, len(rows))
	for i, r := range rows {
		desc := ""
		if r.Description != nil {
			desc = fmt.Sprintf("%v", r.Description)
		}
		result[i] = models.Phase{
			PhaseId:     int(r.PhaseID),
			ProjectId:   int(r.ProjectID),
			Name:        r.Name,
			Description: desc,
			StartDate:   r.StartDate,
			EndDate:     r.EndDate,
			Sort:        int(r.Sort),
		}
	}
	return result, nil
}

func (s *Services) UpdatePhase(ctx context.Context, p models.Phase) error {
	params := projectsdb.UpdatePhaseParams{
		PhaseID:     int64(p.PhaseId),
		Name:          p.Name,
		Description:   p.Description,
		StartDate:     p.StartDate,
		EndDate:       p.EndDate,
		Sort:          int64(p.Sort),
	}
	_, err := s.Queries.UpdatePhase(ctx, params)
	return err
}

func (s *Services) DeletePhase(ctx context.Context, phaseId int64) error {
	if err := s.Queries.UnlinkTasksFromPhase(ctx, phaseId); err != nil {
		return fmt.Errorf("failed to unlink tasks from phase: %w", err)
	}
	if err := s.Queries.DeletePhase(ctx, phaseId); err != nil {
		return fmt.Errorf("failed to delete phase: %w", err)
	}
	return nil
}

// --- Task Ownership, Scheduling, Blockers ---

func (s *Services) SetTaskOwner(ctx context.Context, taskId int64, owner models.Owner) error {
	return s.Queries.SetTaskOwner(ctx, projectsdb.SetTaskOwnerParams{
		TaskID: taskId,
		Owner:  string(owner),
	})
}

func (s *Services) SetTaskSchedule(ctx context.Context, taskId int64, date string, phaseId *int) error {
	var pid interface{}
	if phaseId != nil {
		pid = int64(*phaseId)
	}
	return s.Queries.SetTaskSchedule(ctx, projectsdb.SetTaskScheduleParams{
		TaskID:        taskId,
		ScheduledDate: interface{}(date),
		PhaseID:       pid,
	})
}

func (s *Services) SetTaskBlocker(ctx context.Context, taskId int64, text string) error {
	return s.Queries.SetTaskBlocker(ctx, projectsdb.SetTaskBlockerParams{
		TaskID:      taskId,
		BlockerText: interface{}(text),
	})
}

func (s *Services) ResolveTaskBlocker(ctx context.Context, taskId int64) error {
	return s.Queries.ResolveTaskBlocker(ctx, taskId)
}

func (s *Services) GetTodaysTasks(ctx context.Context, projectId int64) ([]models.Task, error) {
	rows, err := s.Queries.GetTodaysTasks(ctx, projectId)
	if err != nil {
		return nil, err
	}
	return s.rowsToTasks(ctx, rows)
}

func (s *Services) GetTasksByOwner(ctx context.Context, projectId int64, owner models.Owner) ([]models.Task, error) {
	rows, err := s.Queries.GetTasksForOwner(ctx, projectsdb.GetTasksForOwnerParams{
		ProjectID: projectId,
		Owner:     string(owner),
	})
	if err != nil {
		return nil, err
	}
	return s.rowsToTasks(ctx, rows)
}

func (s *Services) rowsToTasks(ctx context.Context, rows []projectsdb.Task) ([]models.Task, error) {
	result := make([]models.Task, 0, len(rows))
	for _, t := range rows {
		depRows, err := s.Queries.GetDependenciesForTask(ctx, t.TaskID)
		if err != nil {
			return nil, fmt.Errorf("Error retrieving Dependencies from DB: %w", err)
		}
		deps := s.DependencyTransformer.FromDbRows(depRows)
		result = append(result, models.NewTask(t, deps))
	}
	return result, nil
}

// --- Task Notes ---

func (s *Services) AddTaskNote(ctx context.Context, taskId int64, note string) error {
	return s.Queries.AddTaskNote(ctx, projectsdb.AddTaskNoteParams{
		TaskID: taskId,
		Note:   note,
	})
}

func (s *Services) GetTaskNotes(ctx context.Context, taskId int64) ([]models.TaskNote, error) {
	rows, err := s.Queries.GetTaskNotes(ctx, taskId)
	if err != nil {
		return nil, err
	}
	result := make([]models.TaskNote, len(rows))
	for i, r := range rows {
		result[i] = models.TaskNote{
			NoteID:    int(r.NoteID),
			TaskID:    int(r.TaskID),
			Note:      r.Note,
			CreatedAt: r.CreatedAt,
		}
	}
	return result, nil
}

// --- Exit Criteria ---

func (s *Services) AddExitCriterion(ctx context.Context, phaseId int64, desc string, sort int) (models.ExitCriterion, error) {
	row, err := s.Queries.CreateExitCriterion(ctx, projectsdb.CreateExitCriterionParams{
		PhaseID:     phaseId,
		Description: desc,
		IsCompleted: 0,
		Sort:        int64(sort),
	})
	if err != nil {
		return models.ExitCriterion{}, err
	}
	return models.ExitCriterion{
		CriterionID: int(row.CriterionID),
		PhaseID:     int(row.PhaseID),
		Description: row.Description,
		IsCompleted: row.IsCompleted != 0,
		Sort:        int(row.Sort),
	}, nil
}

func (s *Services) GetExitCriteria(ctx context.Context, phaseId int64) ([]models.ExitCriterion, error) {
	rows, err := s.Queries.GetExitCriteriaForPhase(ctx, phaseId)
	if err != nil {
		return nil, err
	}
	result := make([]models.ExitCriterion, len(rows))
	for i, r := range rows {
		result[i] = models.ExitCriterion{
			CriterionID: int(r.CriterionID),
			PhaseID:     int(r.PhaseID),
			Description: r.Description,
			IsCompleted: r.IsCompleted != 0,
			Sort:        int(r.Sort),
		}
	}
	return result, nil
}

func (s *Services) CompleteExitCriterion(ctx context.Context, criterionId int64) error {
	return s.Queries.CompleteExitCriterion(ctx, criterionId)
}

// --- Cross-task dependency lookups ---

func (s *Services) GetDownstreamTasks(ctx context.Context, taskId int64) ([]int64, error) {
	return s.Queries.GetDownstreamTasks(ctx, taskId)
}

func (s *Services) GetUpstreamTasks(ctx context.Context, taskId int64) ([]int64, error) {
	return s.Queries.GetUpstreamTasks(ctx, taskId)
}
