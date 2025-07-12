package tui

import (
	"context"

	"github.com/frozenkro/mcpsequencer/internal/services"
	"github.com/frozenkro/mcpsequencer/internal/tui/components/projects"
	"github.com/frozenkro/mcpsequencer/internal/tui/components/tasks"
)

func InitialModel() (Model, error) {
	ctx := context.Background()
	svc := services.Services{}

	width := 100 // Default width
	height := 30 // Default height

	// Initialize components
	projectsModel, err := projects.NewModel(svc, ctx, width, height)
	if err != nil {
		return Model{}, err
	}

	tasksModel := tasks.NewModel(width, height)

	// Create main model
	model := Model{
		Projects:   projectsModel,
		Tasks:      tasksModel,
		ActivePane: ProjectPane,
		Services:   svc,
		Context:    ctx,
		Width:      width,
		Height:     height,
	}

	return model, nil
}
