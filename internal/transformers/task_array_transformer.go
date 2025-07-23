package transformers

import (
	"encoding/json"

	"github.com/frozenkro/mcpsequencer/internal/models"
	"github.com/frozenkro/mcpsequencer/internal/projectsdb"
)

type TaskArrayTransformer struct{}

func (t TaskArrayTransformer) ParseFromJson(tasks string, projectId int) ([]projectsdb.Task, error) {
	argsSl := []models.CreateTaskArgs{}
	if err := json.Unmarshal([]byte(tasks), &argsSl); err != nil {
		return nil, TasksUnmarshalError{}
	}

	result := []projectsdb.Task{}
	for _, taskArgs := range argsSl {

		jsonDeps, err := json.Marshal(taskArgs.Dependencies)
		if err != nil {
			return nil, DepsMarshalError{Deps: taskArgs.Dependencies, Err: err}
		}

		task := projectsdb.Task{
			ProjectID:        int64(projectId),
			Name:             taskArgs.Name,
			Description:      taskArgs.Description,
			Sort:             int64(taskArgs.SortId),
			DependenciesJson: string(jsonDeps),
		}

		result = append(result, task)
	}
	return result, nil
}
