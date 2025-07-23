# MCP Sequencer

A project planning and task management tool designed for LLM-driven development workflows. MCP Sequencer provides both an MCP (Model Context Protocol) server for seamless LLM integration and a terminal user interface (TUI) for direct developer control.

## Features

- **Dual Interface**: MCP server for LLM integration + TUI for manual project management
- **Project Organization**: Create and manage development projects with structured task lists
- **Task Dependencies**: Define task dependencies and execution order
- **Task Status Tracking**: Track task progress (pending, in-progress, completed)
- **SQLite Database**: Persistent storage for projects and tasks
- **HTTP & Stdio Transport**: Flexible MCP server deployment options

## Architecture

### MCP Server
The MCP server exposes project and task management functionality to LLMs through standardized tools:

- **Project Management**: Create, update, delete projects
- **Task Management**: Add tasks, track progress, manage dependencies
- **Data Retrieval**: Get project lists, task details, and instructions

### TUI Application
A terminal-based interface for direct developer interaction:

- **Project Browser**: Navigate and manage projects
- **Task Viewer**: View and update task status
- **Real-time Updates**: Synchronized with the database

## Installation

### Prerequisites
- Go 1.24.3 or later
- SQLite3

### Build from Source
```bash
# Clone the repository
git clone https://github.com/frozenkro/mcpsequencer.git
cd mcpsequencer

# Build both binaries
make build

# Or build individually
make build.mcp  # MCP server
make build.tui  # TUI application
```

## Usage

### MCP Server

#### Stdio Mode (Default)
```bash
# Run directly
go run cmd/mcp/main.go

# Or use built binary
./build/mcpsequencer-mcp
```

#### HTTP Mode
```bash
# Default port (8080)
go run cmd/mcp/main.go --http

# Custom port
go run cmd/mcp/main.go --http 3000

# Using Makefile
make run.mcp
```

### TUI Application
```bash
# Run directly
go run cmd/tui/main.go

# Or use built binary
./build/mcpsequencer-tui

# Using Makefile
make run.tui
```

### MCP Tools Available

The MCP server exposes the following tools for LLM integration:

| Tool | Description |
|------|-------------|
| `createProject` | Create a new project with tasks |
| `updateProject` | Update project details |
| `deleteProject` | Delete a project and its tasks |
| `addTask` | Add a new task to a project |
| `beginTask` | Mark a task as in-progress |
| `completeTask` | Mark a task as completed |
| `getAllProjects` | Get list of all projects |
| `getTasksForProject` | Get all tasks for a specific project |
| `getTaskListInstructions` | Get task formatting instructions |

## Development

### Running Tests
```bash
make test
```

### Development Commands
```bash
# Run MCP server in development mode
make run.mcp

# Run TUI application
make run.tui

# Debug MCP server
make debug.mcp
```

### Project Structure
```
mcpsequencer/
├── cmd/
│   ├── mcp/           # MCP server entry point
│   └── tui/           # TUI application entry point
├── internal/
│   ├── db/            # Database layer
│   ├── mcp/           # MCP server implementation
│   │   ├── handlers/  # MCP tool handlers
│   │   └── tools/     # MCP tool definitions
│   ├── services/      # Business logic
│   ├── tui/           # TUI implementation
│   └── models/        # Data models
└── build/             # Compiled binaries
```

## Configuration

The application uses environment variables and command-line flags for configuration:

- `--http`: Enable HTTP transport mode for MCP server
- `--dev`: Development mode (when debugging)

Database files are created automatically in the working directory.

## Integration with LLMs

The MCP server is designed to work with LLM clients that support the Model Context Protocol. Configure your LLM client to connect to the MCP server:

**Stdio Mode**: Point your client to the `mcpsequencer-mcp` binary
**HTTP Mode**: Connect to `http://localhost:8080` (or your configured port)

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Run tests: `make test`
5. Submit a pull request

## License

This project is open source. Please check the repository for license details.
