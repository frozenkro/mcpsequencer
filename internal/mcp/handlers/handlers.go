package handlers

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/frozenkro/mcpsequencer/internal/globals"
	"github.com/frozenkro/mcpsequencer/internal/models"
	"github.com/frozenkro/mcpsequencer/internal/services"
	"github.com/mark3labs/mcp-go/mcp"
)

type Handlers struct {
	svc services.Services
}

func NewHandlers() *Handlers {
	return &Handlers{svc: services.NewServices()}
}

func (h *Handlers) CreateProjectHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	name, err := request.RequireString(string(globals.ProjectName))
	if err != nil {
		return requiredParamError(globals.ProjectName, err), nil
	}

	description, err := request.RequireString(string(globals.ProjectDesc))
	if err != nil {
		return requiredParamError(globals.ProjectName, err), nil
	}

	dir := request.GetString(string(globals.ProjectDir), "")

	tasks, err := request.RequireString(string(globals.Tasks))
	if err != nil {
		return requiredParamError(globals.Tasks, err), nil
	}

	args := models.CreateProjectArgs{
		Name:        name,
		Description: description,
		Directory:   dir,
		TasksJson:   tasks,
	}

	err = h.svc.CreateProject(ctx, args)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Project %s created successfully!", name)), nil
}

func (h *Handlers) UpdateProjectHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projectId, err := request.RequireInt(string(globals.ProjectId))
	if err != nil {
		return requiredParamError(globals.ProjectId, err), nil
	}

	name, err := request.RequireString(string(globals.ProjectName))
	if err != nil {
		return requiredParamError(globals.ProjectName, err), nil
	}

	description, err := request.RequireString(string(globals.ProjectDesc))
	if err != nil {
		return requiredParamError(globals.ProjectName, err), nil
	}

	dir := request.GetString(string(globals.ProjectDir), "")

	args := models.UpdateProjectArgs{
		ProjectId: projectId,
		Fields: models.CreateProjectArgs{
			Name:        name,
			Description: description,
			Directory:   dir,
		},
	}

	err = h.svc.UpdateProject(ctx, args)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Project %v updated successfully!", projectId)), nil
}

func (h *Handlers) DeleteProjectHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projectId, err := request.RequireInt(string(globals.ProjectId))
	if err != nil {
		return requiredParamError(globals.ProjectId, err), nil
	}

	err = h.svc.DeleteProject(ctx, int64(projectId))
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Project %v deleted successfully", projectId)), nil
}

func (h *Handlers) AddTaskHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projectId, err := request.RequireInt(string(globals.ProjectId))
	if err != nil {
		return requiredParamError(globals.ProjectId, err), nil
	}

	name, err := request.RequireString(string(globals.TaskName))
	if err != nil {
		return requiredParamError(globals.TaskName, err), nil
	}

	description, err := request.RequireString(string(globals.Description))
	if err != nil {
		return requiredParamError(globals.Description, err), nil
	}

	sort, err := request.RequireInt(string(globals.SortId))
	if err != nil {
		return requiredParamError(globals.SortId, err), nil
	}

	deps, err := request.RequireIntSlice(string(globals.Dependencies))
	if err != nil {
		return requiredParamError(globals.SortId, err), nil
	}

	args := models.CreateTaskArgs{
		Name:         name,
		Description:  description,
		SortId:       sort,
		Dependencies: deps,
	}
	err = h.svc.AddTask(ctx, int64(projectId), args)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultText("Task added successfully"), nil
}

func (h *Handlers) BeginTaskHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	taskId, err := request.RequireInt(string(globals.TaskId))
	if err != nil {
		return requiredParamError(globals.TaskId, err), nil
	}

	err = h.svc.UpdateTaskState(ctx, int64(taskId), services.StateInProgress)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultText("Task in progress"), nil
}

func (h *Handlers) CompleteTaskHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	taskId, err := request.RequireInt(string(globals.TaskId))
	if err != nil {
		return requiredParamError(globals.TaskId, err), nil
	}

	err = h.svc.UpdateTaskState(ctx, int64(taskId), services.StateComplete)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultText("Task completed successfully"), nil
}

func (h *Handlers) GetProjectsHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projects, err := h.svc.GetProjects(ctx)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	projectsJson, err := json.Marshal(projects)
	return mcp.NewToolResultText(string(projectsJson)), nil
}

func (h *Handlers) GetTasksHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projectId, err := request.RequireInt(string(globals.ProjectId))
	if err != nil {
		return requiredParamError(globals.ProjectId, err), nil
	}

	tasks, err := h.svc.GetTasksByProject(ctx, int64(projectId))
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	tasksJson, err := json.Marshal(tasks)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultText(string(tasksJson)), nil
}

func (h *Handlers) GetTaskListInstructionsHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	instructions := `
Tasks should be defined in json, with the goal of being as parallelizable as possible.

Tasks are defined as follows:
{
	name: string,
	description: string, // Include a detailed, markdown-formatted description of the task
	sortId: int, // Used for visually ordering tasks and as a FK for dependencies
	dependencies: []int // the sortIds of any tasks that must be completed before this task.
}

You should pass a json array of 'Task' items.
	`
	return mcp.NewToolResultText(instructions), nil
}

func (h *Handlers) CreatePhaseHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projectId, err := request.RequireInt(string(globals.ProjectId))
	if err != nil {
		return requiredParamError(globals.ProjectId, err), nil
	}

	name, err := request.RequireString(string(globals.PhaseName))
	if err != nil {
		return requiredParamError(globals.PhaseName, err), nil
	}

	description := request.GetString(string(globals.Description), "")
	startDate := request.GetString(string(globals.StartDate), "")
	endDate := request.GetString(string(globals.EndDate), "")

	sort, err := request.RequireInt(string(globals.SortId))
	if err != nil {
		return requiredParamError(globals.SortId, err), nil
	}

	phase, err := h.svc.CreatePhase(ctx, models.Phase{
		ProjectId:   projectId,
		Name:        name,
		Description: description,
		StartDate:   startDate,
		EndDate:     endDate,
		Sort:        sort,
	})
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	phaseJson, err := json.Marshal(phase)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultText(string(phaseJson)), nil
}

func (h *Handlers) GetPhasesHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projectId, err := request.RequireInt(string(globals.ProjectId))
	if err != nil {
		return requiredParamError(globals.ProjectId, err), nil
	}

	phases, err := h.svc.GetPhases(ctx, int64(projectId))
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	phasesJson, err := json.Marshal(phases)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultText(string(phasesJson)), nil
}

func (h *Handlers) UpdatePhaseHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	phaseId, err := request.RequireInt(string(globals.PhaseId))
	if err != nil {
		return requiredParamError(globals.PhaseId, err), nil
	}

	name, err := request.RequireString(string(globals.PhaseName))
	if err != nil {
		return requiredParamError(globals.PhaseName, err), nil
	}

	description := request.GetString(string(globals.Description), "")
	startDate := request.GetString(string(globals.StartDate), "")
	endDate := request.GetString(string(globals.EndDate), "")

	sort, err := request.RequireInt(string(globals.SortId))
	if err != nil {
		return requiredParamError(globals.SortId, err), nil
	}

	err = h.svc.UpdatePhase(ctx, models.Phase{
		PhaseId:     phaseId,
		Name:        name,
		Description: description,
		StartDate:   startDate,
		EndDate:     endDate,
		Sort:        sort,
	})
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultText("Phase updated successfully"), nil
}

func (h *Handlers) DeletePhaseHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	phaseId, err := request.RequireInt(string(globals.PhaseId))
	if err != nil {
		return requiredParamError(globals.PhaseId, err), nil
	}

	err = h.svc.DeletePhase(ctx, int64(phaseId))
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultText("Phase deleted successfully"), nil
}

func (h *Handlers) SetTaskOwnerHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	taskId, err := request.RequireInt(string(globals.TaskId))
	if err != nil {
		return requiredParamError(globals.TaskId, err), nil
	}

	owner, err := request.RequireString(string(globals.Owner))
	if err != nil {
		return requiredParamError(globals.Owner, err), nil
	}

	err = h.svc.SetTaskOwner(ctx, int64(taskId), models.Owner(owner))
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultText("Task owner updated successfully"), nil
}

func (h *Handlers) SetTaskScheduleHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	taskId, err := request.RequireInt(string(globals.TaskId))
	if err != nil {
		return requiredParamError(globals.TaskId, err), nil
	}

	scheduledDate := request.GetString(string(globals.ScheduledDate), "")

	var phaseIdPtr *int
	phaseIdVal, perr := request.RequireInt(string(globals.PhaseId))
	if perr == nil {
		phaseIdPtr = &phaseIdVal
	}

	err = h.svc.SetTaskSchedule(ctx, int64(taskId), scheduledDate, phaseIdPtr)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultText("Task schedule updated successfully"), nil
}

func (h *Handlers) SetTaskBlockerHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	taskId, err := request.RequireInt(string(globals.TaskId))
	if err != nil {
		return requiredParamError(globals.TaskId, err), nil
	}

	blockerText, err := request.RequireString(string(globals.BlockerText))
	if err != nil {
		return requiredParamError(globals.BlockerText, err), nil
	}

	err = h.svc.SetTaskBlocker(ctx, int64(taskId), blockerText)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultText("Task blocker set successfully"), nil
}

func (h *Handlers) ResolveTaskBlockerHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	taskId, err := request.RequireInt(string(globals.TaskId))
	if err != nil {
		return requiredParamError(globals.TaskId, err), nil
	}

	err = h.svc.ResolveTaskBlocker(ctx, int64(taskId))
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultText("Task blocker resolved successfully"), nil
}

func (h *Handlers) AddTaskNoteHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	taskId, err := request.RequireInt(string(globals.TaskId))
	if err != nil {
		return requiredParamError(globals.TaskId, err), nil
	}

	note, err := request.RequireString(string(globals.Note))
	if err != nil {
		return requiredParamError(globals.Note, err), nil
	}

	err = h.svc.AddTaskNote(ctx, int64(taskId), note)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultText("Task note added successfully"), nil
}

func (h *Handlers) GetTaskNotesHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	taskId, err := request.RequireInt(string(globals.TaskId))
	if err != nil {
		return requiredParamError(globals.TaskId, err), nil
	}

	notes, err := h.svc.GetTaskNotes(ctx, int64(taskId))
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	notesJson, err := json.Marshal(notes)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultText(string(notesJson)), nil
}

func (h *Handlers) GetTodaysTasksHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projectId, err := request.RequireInt(string(globals.ProjectId))
	if err != nil {
		return requiredParamError(globals.ProjectId, err), nil
	}

	tasks, err := h.svc.GetTodaysTasks(ctx, int64(projectId))
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	tasksJson, err := json.Marshal(tasks)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultText(string(tasksJson)), nil
}

func (h *Handlers) GetTasksByOwnerHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projectId, err := request.RequireInt(string(globals.ProjectId))
	if err != nil {
		return requiredParamError(globals.ProjectId, err), nil
	}

	owner, err := request.RequireString(string(globals.Owner))
	if err != nil {
		return requiredParamError(globals.Owner, err), nil
	}

	tasks, err := h.svc.GetTasksByOwner(ctx, int64(projectId), models.Owner(owner))
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	tasksJson, err := json.Marshal(tasks)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultText(string(tasksJson)), nil
}

func (h *Handlers) AddExitCriterionHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	phaseId, err := request.RequireInt(string(globals.PhaseId))
	if err != nil {
		return requiredParamError(globals.PhaseId, err), nil
	}

	description, err := request.RequireString(string(globals.Description))
	if err != nil {
		return requiredParamError(globals.Description, err), nil
	}

	sort, err := request.RequireInt(string(globals.SortId))
	if err != nil {
		return requiredParamError(globals.SortId, err), nil
	}

	criterion, err := h.svc.AddExitCriterion(ctx, int64(phaseId), description, sort)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	criterionJson, err := json.Marshal(criterion)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultText(string(criterionJson)), nil
}

func (h *Handlers) CompleteExitCriterionHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	criterionId, err := request.RequireInt(string(globals.CriterionId))
	if err != nil {
		return requiredParamError(globals.CriterionId, err), nil
	}

	err = h.svc.CompleteExitCriterion(ctx, int64(criterionId))
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultText("Exit criterion completed successfully"), nil
}

func (h *Handlers) GetExitCriteriaHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	phaseId, err := request.RequireInt(string(globals.PhaseId))
	if err != nil {
		return requiredParamError(globals.PhaseId, err), nil
	}

	criteria, err := h.svc.GetExitCriteria(ctx, int64(phaseId))
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	criteriaJson, err := json.Marshal(criteria)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultText(string(criteriaJson)), nil
}

func (h *Handlers) GetDownstreamTasksHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	taskId, err := request.RequireInt(string(globals.TaskId))
	if err != nil {
		return requiredParamError(globals.TaskId, err), nil
	}

	tasks, err := h.svc.GetDownstreamTasks(ctx, int64(taskId))
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	tasksJson, err := json.Marshal(tasks)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultText(string(tasksJson)), nil
}

func (h *Handlers) GetUpstreamTasksHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	taskId, err := request.RequireInt(string(globals.TaskId))
	if err != nil {
		return requiredParamError(globals.TaskId, err), nil
	}

	tasks, err := h.svc.GetUpstreamTasks(ctx, int64(taskId))
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	tasksJson, err := json.Marshal(tasks)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultText(string(tasksJson)), nil
}

func requiredParamError(param globals.McpArg, err error) *mcp.CallToolResult {
	return mcp.NewToolResultError(fmt.Sprintf("Required parameter '%v' returned error: %v", param, err.Error()))
}
