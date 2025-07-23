package tools

import (
	"fmt"

	"github.com/frozenkro/mcpsequencer/internal/globals"
	"github.com/mark3labs/mcp-go/mcp"
)

var CreateProjectTool = mcp.NewTool("createProject",
	mcp.WithDescription("Create a new LLM-driven project (collection of small tasks)"),
	mcp.WithString(string(globals.ProjectName),
		mcp.Required(),
		mcp.Description("Name of the new project"),
	),
	mcp.WithString("Tasks",
		mcp.Required(),
		mcp.Description("Json array that follows specific schema and rules. Make an mcp call to the getTaskListInstructions tool for instructions."),
	),
)

var RenameProjectTool = mcp.NewTool("renameProject",
	mcp.WithDescription("Rename an LLM-driven project"),
	mcp.WithString(string(globals.ProjectName),
		mcp.Required(),
		mcp.Description("New project name"),
	),
	mcp.WithNumber(string(globals.ProjectId),
		mcp.Required(),
		mcp.Description(fmt.Sprintf("%v of project to rename", globals.ProjectId)),
	),
)

var DeleteProjectTool = mcp.NewTool("deleteProject",
	mcp.WithDescription("Delete an LLM-driven project"),
	mcp.WithNumber(string(globals.ProjectId),
		mcp.Required(),
		mcp.Description(fmt.Sprintf("%v of project to delete", globals.ProjectId)),
	),
)

var AddTaskTool = mcp.NewTool("addTask",
	mcp.WithDescription("Add new task to an LLM-driven project"),
	mcp.WithNumber(string(globals.ProjectId),
		mcp.Required(),
		mcp.Description("ID of Project to add task to"),
	),
	mcp.WithString(string(globals.TaskName),
		mcp.Required(),
		mcp.Description("Name of the task"),
	),
	mcp.WithString(string(globals.Description),
		mcp.Required(),
		mcp.Description("Detailed description of the task"),
	),
	mcp.WithNumber(string(globals.SortId),
		mcp.Required(),
		mcp.Description("Sort order of new task. 0 for first, -1 for last"),
	),
	mcp.WithArray(string(globals.Dependencies),
		mcp.Description("The sortIds of any tasks that must be completed before this task. Tasks are otherwise assumed to be parallelizable."),
	),
)

var BeginTaskTool = mcp.NewTool("beginTask",
	mcp.WithDescription("Indicate that a task is in-progress"),
	mcp.WithNumber(string(globals.TaskId),
		mcp.Required(),
		mcp.Description(fmt.Sprintf("%v of task (NOT the %v)", globals.TaskId, globals.SortId)),
	),
)

var CompleteTaskTool = mcp.NewTool("completeTask",
	mcp.WithDescription("Indicate that a task is completed"),
	mcp.WithNumber(string(globals.TaskId),
		mcp.Required(),
		mcp.Description(fmt.Sprintf("%v of task (NOT the %v)", globals.TaskId, globals.SortId)),
	),
)

var GetProjectsTool = mcp.NewTool("getAllProjects",
	mcp.WithDescription("Get a list of all project names and IDs"),
)

var GetTasksTool = mcp.NewTool("getTasksForProject",
	mcp.WithDescription("Get all tasks for a project"),
	mcp.WithNumber(string(globals.ProjectId),
		mcp.Required(),
		mcp.Description("ID of project"),
	),
)

var GetTaskListInstructionsTool = mcp.NewTool("getTaskListInstructions",
	mcp.WithDescription("Get instructions on how tasks are defined and organized"),
)
