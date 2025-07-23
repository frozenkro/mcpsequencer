package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/frozenkro/mcpsequencer/internal/db"
	"github.com/frozenkro/mcpsequencer/internal/globals"
	"github.com/frozenkro/mcpsequencer/internal/mcp/handlers"
	"github.com/frozenkro/mcpsequencer/internal/mcp/tools"
	"github.com/mark3labs/mcp-go/server"
)

var DefaultPort int = 8080

func main() {
	if err := globals.Init(); err != nil {
		log.Fatalf("Application Initialization failed. \nError: %v\n", err.Error())
	}
	db.Init()
	h := handlers.NewHandlers()

	// Create a new MCP server
	s := server.NewMCPServer(
		"LLM Project Sequencer",
		"0.0.1",
		server.WithToolCapabilities(false),
	)

	s.AddTool(tools.CreateProjectTool, h.CreateProjectHandler)
	s.AddTool(tools.UpdateProjectTool, h.UpdateProjectHandler)
	s.AddTool(tools.DeleteProjectTool, h.DeleteProjectHandler)
	s.AddTool(tools.AddTaskTool, h.AddTaskHandler)
	s.AddTool(tools.BeginTaskTool, h.BeginTaskHandler)
	s.AddTool(tools.CompleteTaskTool, h.CompleteTaskHandler)
	s.AddTool(tools.GetProjectsTool, h.GetProjectsHandler)
	s.AddTool(tools.GetTasksTool, h.GetTasksHandler)
	s.AddTool(tools.GetTaskListInstructionsTool, h.GetTaskListInstructionsHandler)

	if http, port := isHTTP(); http {
		log.Printf("Starting HTTP Server...")

		httpServer := server.NewStreamableHTTPServer(s)

		portStr := fmt.Sprintf(":%v", port)
		if err := httpServer.Start(portStr); err != nil {
			log.Printf("HTTP Server error: %v\n", err)
		}

	} else {
		log.Printf("Starting Stdio Server...")

		if err := server.ServeStdio(s); err != nil {
			log.Printf("Stdio Server error: %v\n", err)
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
