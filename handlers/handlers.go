package handlers

import (
	"context"
	"fmt"
	"strconv"

	"github.com/frozenkro/mcpsequencer/services"
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
	projectIdStr, err := request.RequireString("ProjectID")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	projectId, err := strconv.Atoi(projectIdStr)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	name, err := request.RequireString("Name")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	err = svc.RenameProject(ctx, int64(projectId), name)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultText("Project renamed successfully!"), nil
}
