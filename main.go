package main

import (
	"context"
	"fmt"

	"github.com/frozenkro/mcpsequencer/db"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func createProjectHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	name, err := request.RequireString("ProjectName")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	tasks, err := request.RequireStringSlice("Tasks")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	project := db.Project{
		Name:  name,
		Tasks: tasks,
	}
	err = db.InsertProject(project)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Project %s created successfully!", name)), nil
}

func main() {
	// Create a new MCP server
	s := server.NewMCPServer(
		"LLM Project Sequencer",
		"0.0.1",
		server.WithToolCapabilities(false),
	)

	// Add tool
	tool := mcp.NewTool("create_project",
		mcp.WithDescription("Create a new LLM-driven project (collection of small tasks)"),
		mcp.WithString("ProjectName",
			mcp.Required(),
			mcp.Description("Name of the new project"),
		),
		mcp.WithArray("Tasks",
			mcp.Required(),
			mcp.Description("Ordered list of tasks to complete project"),
		),
	)

	// Add tool handler
	s.AddTool(tool, createProjectHandler)

	// Start the stdio server
	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
