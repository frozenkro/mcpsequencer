package transformers

import (
	"encoding/json"

	"github.com/frozenkro/mcpsequencer/internal/models"
	"github.com/frozenkro/mcpsequencer/internal/projectsdb"
)

type TaskArrayTransformer struct{}

func (t TaskArrayTransformer) ParseFromJson(tasks []string, projectId int) ([]projectsdb.Task, error) {
	result := []projectsdb.Task{}
	for _, jsonT := range tasks {
		args := models.CreateTaskArgs{}

		if err := json.Unmarshal([]byte(jsonT), args); err != nil {
			return nil, TaskUnmarshalError{TaskJson: jsonT, Err: err}
		}

		jsonDeps, err := json.Marshal(args.Dependencies)
		if err != nil {
			return nil, DepsMarshalError{Deps: args.Dependencies, Err: err}
		}

		task := projectsdb.Task{
			ProjectID:        int64(projectId),
			Name:             args.Name,
			Description:      args.Description,
			Sort:             int64(args.SortId),
			DependenciesJson: string(jsonDeps),
		}

		result = append(result, task)
	}
	return result, nil
}
