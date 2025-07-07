package handlers

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/frozenkro/mcpsequencer/internal/services"
	"github.com/mark3labs/mcp-go/mcp"
)

var svc services.Services

func CreateProjectHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	name, err := request.RequireString("ProjectName")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	tasks, err := request.RequireStringSlice("Tasks")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	err = svc.CreateProject(ctx, name, tasks)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Project %s created successfully!", name)), nil
}

func RenameProjectHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projectId, err := request.RequireInt("ProjectID")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	name, err := request.RequireString("ProjectName")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	err = svc.RenameProject(ctx, int64(projectId), name)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultText("Project renamed successfully!"), nil
}

func DeleteProjectHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projectId, err := request.RequireInt("ProjectID")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	err = svc.DeleteProject(ctx, int64(projectId))
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Project %v deleted successfully", projectId)), nil
}

func AddTaskHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projectId, err := request.RequireInt("ProjectID")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	description, err := request.RequireString("Description")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	sort := request.GetInt("Sort", int(services.TaskSortLast))

	err = svc.AddTask(ctx, int64(projectId), description, int64(sort))
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultText("Task added successfully"), nil
}

func BeginTaskHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	taskId, err := request.RequireInt("TaskID")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	err = svc.UpdateTaskState(ctx, int64(taskId), services.StateInProgress)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultText("Task completed successfully"), nil
}

func CompleteTaskHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	taskId, err := request.RequireInt("TaskID")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	err = svc.UpdateTaskState(ctx, int64(taskId), services.StateComplete)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultText("Task completed successfully"), nil
}

func GetProjectsHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projects, err := svc.GetProjects(ctx)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	projectsJson, err := json.Marshal(projects)
	return mcp.NewToolResultText(string(projectsJson)), nil
}

func GetTasksHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projectId, err := request.RequireInt("ProjectID")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	tasks, err := svc.GetTasksByProject(ctx, int64(projectId))
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	tasksJson, err := json.Marshal(tasks)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultText(string(tasksJson)), nil
}
