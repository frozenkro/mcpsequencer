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
	s.AddTool(tool, handlers.CreateProjectHandler)

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
