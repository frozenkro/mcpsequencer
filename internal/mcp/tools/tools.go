package tools

import "github.com/mark3labs/mcp-go/mcp"

var CreateProjectTool = mcp.NewTool("createProject",
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

var RenameProjectTool = mcp.NewTool("renameProject",
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

var DeleteProjectTool = mcp.NewTool("deleteProject",
	mcp.WithDescription("Delete an LLM-driven project"),
	mcp.WithNumber("ProjectID",
		mcp.Required(),
		mcp.Description("ProjectID of project to delete"),
	),
)

var AddTaskTool = mcp.NewTool("addTask",
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

var BeginTaskTool = mcp.NewTool("beginTask",
	mcp.WithDescription("Indicate that a task is in-progress"),
	mcp.WithNumber("TaskID",
		mcp.Required(),
		mcp.Description("ID of task"),
	),
)

var CompleteTaskTool = mcp.NewTool("completeTask",
	mcp.WithDescription("Indicate that a task is completed"),
	mcp.WithNumber("TaskID",
		mcp.Required(),
		mcp.Description("ID of task"),
	),
)

var GetProjectsTool = mcp.NewTool("getAllProjects",
	mcp.WithDescription("Get a list of all project names and IDs"),
)

var GetTasksTool = mcp.NewTool("getTasksForProject",
	mcp.WithDescription("Get all tasks for a project"),
	mcp.WithNumber("ProjectID",
		mcp.Required(),
		mcp.Description("ID of project"),
	),
)

var GetTaskListInstructionsTool = mcp.NewTool("getTaskListInstructions",
	mcp.WithDescription("Get instructions on how tasks are defined and organized"),
)
