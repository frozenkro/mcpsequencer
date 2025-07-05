package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/frozenkro/mcpsequencer/db"
	"github.com/frozenkro/mcpsequencer/handlers"
	"github.com/frozenkro/mcpsequencer/services"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

var DefaultPort int = 8080
var Svc services.Services

func main() {
	// Create a new MCP server
	s := server.NewMCPServer(
		"LLM Project Sequencer",
		"0.0.1",
		server.WithToolCapabilities(false),
	)

	// Add tool
	createProjectTool := mcp.NewTool("createProject",
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

	renameProjectTool := mcp.NewTool("renameProject",
		mcp.WithDescription("Rename an LLM-driven project"),
		mcp.WithString("ProjectName",
			mcp.Required(),
			mcp.Description("New project name"),
		),
		mcp.WithNumber("ProjectID",
			mcp.Required(),
			mcp.Description("ProjectID of project to rename"),
		),
	)

	deleteProjectTool := mcp.NewTool("deleteProject",
		mcp.WithDescription("Delete an LLM-driven project"),
		mcp.WithNumber("ProjectID",
			mcp.Required(),
			mcp.Description("ProjectID of project to delete"),
		),
	)

	addTaskTool := mcp.NewTool("addTask",
		mcp.WithDescription("Add new task to an LLM-driven project"),
		mcp.WithNumber("ProjectID",
			mcp.Required(),
			mcp.Description("ID of Project to add task to"),
		),
		mcp.WithString("Description",
			mcp.Required(),
			mcp.Description("Name / Description of new task"),
		),
		mcp.WithNumber("Sort",
			mcp.Description("Sort order of new task. 0-indexed, -1 for last"),
		),
	)

	beginTaskTool := mcp.NewTool("beginTask",
		mcp.WithDescription("Indicate that a task is in-progress"),
		mcp.WithNumber("TaskID",
			mcp.Required(),
			mcp.Description("ID of task"),
		),
	)

	completeTaskTool := mcp.NewTool("completeTask",
		mcp.WithDescription("Indicate that a task is completed"),
		mcp.WithNumber("TaskID",
			mcp.Required(),
			mcp.Description("ID of task"),
		),
	)

	getProjectsTool := mcp.NewTool("getAllProjects",
		mcp.WithDescription("Get a list of all project names and IDs"),
	)

	getTasksTool := mcp.NewTool("getTasksForProject",
		mcp.WithDescription("Get all tasks for a project"),
		mcp.WithNumber("ProjectID",
			mcp.Required(),
			mcp.Description("ID of project"),
		),
	)

	// Add tool handler
	s.AddTool(createProjectTool, handlers.CreateProjectHandler)
	s.AddTool(renameProjectTool, handlers.RenameProjectHandler)
	s.AddTool(deleteProjectTool, handlers.DeleteProjectHandler)
	s.AddTool(addTaskTool, handlers.AddTaskHandler)
	s.AddTool(beginTaskTool, handlers.BeginTaskHandler)
	s.AddTool(completeTaskTool, handlers.CompleteTaskHandler)
	s.AddTool(getProjectsTool, handlers.GetProjectsHandler)
	s.AddTool(getTasksTool, handlers.GetTasksHandler)

	// Set up DB
	db.Init()
	Svc = services.Services{}

	if http, port := isHTTP(); http {
		fmt.Printf("Starting HTTP Server...")

		httpServer := server.NewStreamableHTTPServer(s)

		portStr := fmt.Sprintf(":%v", port)
		if err := httpServer.Start(portStr); err != nil {
			fmt.Printf("HTTP Server error: %v\n", err)
		}

	} else {
		fmt.Printf("Starting Stdio Server...")

		if err := server.ServeStdio(s); err != nil {
			fmt.Printf("Stdio Server error: %v\n", err)
		}
	}
}

func isHTTP() (bool, int) {

	for i, v := range os.Args {
		if v == "--http" {
			fmt.Printf("Running in http transport mode\n")

			if i < len(os.Args)-1 {
				port, err := strconv.Atoi(os.Args[i+1])
				if err != nil {
					fmt.Printf("Error parsing port, setting port to %v", DefaultPort)
					return true, DefaultPort
				}

				return true, port
			}
			return true, DefaultPort
		}
	}

	return false, 0
}
