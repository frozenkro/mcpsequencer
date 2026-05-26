package tui

import (
	"context"
	"log"

	"github.com/frozenkro/sqncr/internal/services"
	"github.com/frozenkro/sqncr/internal/tui/components/projects"
	"github.com/frozenkro/sqncr/internal/tui/components/taskdetail"
	"github.com/frozenkro/sqncr/internal/tui/components/tasks"
	tuilog "github.com/frozenkro/sqncr/internal/tui/logger"
)

func InitialModel(logger *log.Logger) (Model, error) {
	tuilog.SetLogger(logger)
	ctx := context.Background()
	svc := services.NewServices()

	width := 100 // Default width
	height := 30 // Default height

	// Initialize components
	projectsModel, err := projects.NewModel(svc, ctx, width, height)
	if err != nil {
		return Model{}, err
	}

	tasksModel := tasks.NewModel(svc, width, height)
	taskDetailModel := taskdetail.NewModel(svc, width, height)

	// Create main model
	model := Model{
		Projects:   projectsModel,
		Tasks:      tasksModel,
		TaskDetail: taskDetailModel,
		ActivePane: ProjectPane,
		Services:   svc,
		Context:    ctx,
		Width:      width,
		Height:     height,
	}

	return model, nil
}
