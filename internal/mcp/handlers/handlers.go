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

	tasks, err := request.RequireStringSlice(string(globals.Tasks))
	if err != nil {
		return requiredParamError(globals.Tasks, err), nil
	}

	err = h.svc.CreateProject(ctx, name, tasks)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Project %s created successfully!", name)), nil
}

func (h *Handlers) RenameProjectHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projectId, err := request.RequireInt(string(globals.ProjectId))
	if err != nil {
		return requiredParamError(globals.ProjectId, err), nil
	}

	name, err := request.RequireString(string(globals.ProjectName))
	if err != nil {
		return requiredParamError(globals.ProjectName, err), nil
	}

	err = h.svc.RenameProject(ctx, int64(projectId), name)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultText("Project renamed successfully!"), nil
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
		return requiredParamError(globals.ProjectId, err), nil
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

	return mcp.NewToolResultText("Task completed successfully"), nil
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
	description: string, // Include a detailed description of the task
	sortId: int, // Used for visually ordering tasks and as a FK for dependencies
	dependencies: []int // the sortIds of any tasks that must be completed before this task.
}

You should return a json array of 'Task' items.
	`
	return mcp.NewToolResultText(instructions), nil
}

func requiredParamError(param globals.McpArg, err error) *mcp.CallToolResult {
	return mcp.NewToolResultError(fmt.Sprintf("Required parameter '%v' returned error: %v", param, err.Error()))
}
